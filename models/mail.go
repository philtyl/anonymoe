package models

import (
	"io"
	"strings"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/go-xorm/xorm"
	"github.com/philtyl/anonymoe/pkg/setting"
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

	Received     time.Time `xorm:"-" json:"-"`
	ReceivedUnix int64
	Sent         time.Time `xorm:"-" json:"-"`
	SentUnix     int64
}

type RawMailItem struct {
	From      string
	Recipient []string
	Data      io.Reader
	Complete  bool
}

func (m *Mail) BeforeInsert() {
	m.ReceivedUnix = time.Now().Unix()
}

func (m *Mail) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "received_unix":
		m.Received = time.Unix(m.ReceivedUnix, 0).Local()
	case "sent_unix":
		m.Sent = time.Unix(m.SentUnix, 0).Local()
	}
}

func createMail(e *xorm.Session, raw *RawMailItem) (_ *Mail, _ []MailRecipient, err error) {
	log.Trace("Raw Mail Item: %+v:\n", raw)
	m, err := parsemail.Parse(raw.Data)
	if err != nil {
		log.Warn("Unable to parse raw email data: %v", err)
		return
	}

	mailItem := &Mail{
		From:     raw.From,
		Sent:     m.Date,
		Received: time.Now(),
		Subject:  m.Subject,
		Body:     m.HTMLBody,
	}
	if _, err = e.Insert(mailItem); err != nil {
		return nil, nil, err
	}

	var recipients []MailRecipient
	for _, recipient := range raw.Recipient {
		if strings.HasSuffix(recipient, "@"+setting.AppDomain) {
			user, err := GetOrCreateUserByName(e, recipient)
			if err != nil {
				return nil, nil, err
			}
			if user != nil {
				mailRecipient := &MailRecipient{
					MailId:      mailItem.Id,
					RecipientId: user.Id,
				}
				if _, err = e.Insert(mailRecipient); err != nil {
					return nil, nil, err
				}
				recipients = append(recipients, *mailRecipient)
			}
		} else {
			log.Info("Receiving mail for outside address: %s, skipping linkage...", recipient)
		}
	}

	return mailItem, recipients, nil
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
