package models

import (
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/philtyl/anonymoe/pkg/setting"
	"github.com/philtyl/parsemail"
	log "gopkg.in/clog.v1"
)

type MailRecipient struct {
	Id          int64
	MailId      int64
	RecipientId int64 `xorm:"INDEX"`
}

func (m *MailRecipient) String() string {
	return fmt.Sprintf("[ID:%d, MailID:%d, UserID:%d]", m.Id, m.MailId, m.RecipientId)
}

type Mail struct {
	Id      int64
	From    string
	Subject string
	Body    string `xorm:"TEXT"`

	Received     time.Time `xorm:"created"`
	ReceivedUnix int64     `xorm:"created"`
	Sent         time.Time
	SentUnix     int64
}

func (m *Mail) String() string {
	return fmt.Sprintf("[ID:%d, From:%s, Subject:%s]", m.Id, m.From, m.Subject)
}

func (m *Mail) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "sent":
		m.SentUnix = m.Sent.Unix()
	}
}

func (m *Mail) AfterLoad() {
	m.Received = time.Unix(m.ReceivedUnix, 0).Local()
	m.Sent = time.Unix(m.SentUnix, 0).Local()
}

type RawMailItem struct {
	From      string
	Recipient []string
	Data      string
	Complete  bool
}

func createMail(e *xorm.Session, raw *RawMailItem) (_ *Mail, _ []MailRecipient, err error) {
	log.Trace("Received new mail item: %v", raw)
	r := strings.NewReader(raw.Data)
	m, err := parsemail.Parse(r)
	if err != nil {
		log.Warn("Unable to parse raw email data, ignoring: %v", err)
		return
	}

	mailItem := &Mail{
		From:    raw.From,
		Sent:    m.Date,
		Subject: m.Subject,
	}
	if _, err = e.Insert(mailItem); err != nil {
		log.Warn("Unable to insert Mail item into database: %v", err)
		return
	}
	log.Info("Created Mail: [ID:%d, From:<%s>, Subject:\"%s\"]", mailItem.Id, mailItem.From, mailItem.Subject)

	body := m.HTMLBody
	if len(body) == 0 {
		log.Trace("HTMLBody is empty, falling back to TextBody")
		var bytesBody []byte
		bytesBody, err = ioutil.ReadAll(quotedprintable.NewReader(strings.NewReader(m.TextBody)))
		if err != nil {
			log.Warn("Both HTMLBody and TextBody are not able to be parsed, ignoring: %v", err)
			return
		}
		body = string(bytesBody)
	} else {
		body = strings.ReplaceAll(body, "<img src=\"cid:", fmt.Sprintf("<img src=\"%s/inbox/embed/%d/", setting.Config.AppURL, mailItem.Id))
	}

	log.Trace("Updating body for Mail item [ID:%d]:\n%s", mailItem.Id, body)
	mailItem.Body = policy.Sanitize(body)
	if _, err = e.Update(mailItem); err != nil {
		log.Warn("Unable to update body on Mail [ID:%d]", mailItem.Id)
		return
	}

	var recipients []MailRecipient
	for _, recipient := range raw.Recipient {
		if strings.HasSuffix(recipient, "@"+setting.Config.AppDomain) {
			var user *User
			user, err = GetOrCreateUserByName(e, recipient)
			if err != nil || user == nil {
				log.Warn("Unable to create or retrieve user <%s>: %v", recipient, err)
				return
			}

			mailRecipient := &MailRecipient{
				MailId:      mailItem.Id,
				RecipientId: user.Id,
			}
			if _, err = e.Insert(mailRecipient); err != nil {
				log.Warn("Unable to link [Mail:%d] to [User:%d]: %v", mailItem.Id, user.Id, err)
				return
			}
			log.Info("Created Mail Recipient: [ID:%d, MailID:%d, UserID:%d, User:%s]",
				mailRecipient.Id, mailRecipient.MailId, mailRecipient.RecipientId, user.Name)
			recipients = append(recipients, *mailRecipient)
		} else {
			log.Warn("Receiving mail for outside address: <%s>, skipping linkage...", recipient)
		}
	}
	if len(recipients) == 0 {
		return nil, nil, fmt.Errorf("No valid recepients, nothing to do.  [Recipients:%v]", raw.Recipient)
	}

	for _, rawAttachment := range m.Attachments {
		data, _ := ioutil.ReadAll(rawAttachment.Data)
		attachment := &Attachment{
			MailId:      mailItem.Id,
			FileName:    rawAttachment.Filename,
			ContentType: rawAttachment.ContentType,
			Data:        string(data),
		}
		log.Trace("Creating Attachment [MailID:%d, FileName:%s]", attachment.MailId, attachment.FileName)
		attachment, warning := CreateAttachment(e, attachment)
		if warning != nil {
			log.Warn("Unable to insert Attachment (%s) into database: %v", rawAttachment.Filename, warning)
		}
	}

	for _, rawEmbeddedFile := range m.EmbeddedFiles {
		data, _ := ioutil.ReadAll(rawEmbeddedFile.Data)
		embeddedFile := &EmbeddedFile{
			MailId:      mailItem.Id,
			ContentId:   rawEmbeddedFile.CID,
			ContentType: rawEmbeddedFile.ContentType,
			Data:        string(data),
		}
		log.Trace("Creating EmbeddedFile [MailID:%d, ContentID:%s]", embeddedFile.MailId, embeddedFile.ContentId)
		embeddedFile, warning := CreateEmbeddedFile(e, embeddedFile)
		if warning != nil {
			log.Warn("Unable to insert EmbeddedFile [MailID:%d, ContentID:%s] into database: %v", mailItem.Id, rawEmbeddedFile.CID, warning)
		}
	}

	return mailItem, recipients, err
}

func CreateMail(raw *RawMailItem) (mail *Mail, recipients []MailRecipient, err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, nil, err
	}

	mail, recipients, err = createMail(sess, raw)
	if err != nil {
		defer sess.Rollback()
		return nil, nil, err
	}

	return mail, recipients, sess.Commit()
}
