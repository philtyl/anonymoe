package models

import (
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

type RawMailItem struct {
	From      string
	Recipient []string
	Data      string
	Complete  bool
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

func createMail(e *xorm.Session, raw *RawMailItem) (_ *Mail, _ []MailRecipient, err error) {
	log.Trace("Received new mail item: %v", raw)
	r := strings.NewReader(raw.Data)
	m, err := parsemail.Parse(r)
	if err != nil {
		log.Warn("Unable to parse raw email data, ignoring: %v", err)
		return
	}

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
	}

	mailItem := &Mail{
		From:    raw.From,
		Sent:    m.Date,
		Subject: m.Subject,
		Body:    policy.Sanitize(body),
	}
	if _, err = e.Insert(mailItem); err != nil {
		log.Warn("Unable to insert Mail item into database: %v", err)
		return
	}
	log.Info("Created Mail: [ID:%d, From:<%s>, Subject:\"%s\"]", mailItem.Id, mailItem.From, mailItem.Subject)

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
		return nil, nil, err
	}

	return mail, recipients, sess.Commit()
}

func getUserMailCount(u *User) (int64, error) {
	return x.Count(&MailRecipient{RecipientId: u.Id})
}
