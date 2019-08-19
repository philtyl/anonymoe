package models

import (
	"regexp"
	"time"

	"github.com/go-xorm/xorm"
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

	Sent         time.Time `xorm:"-" json:"-"`
	SentUnix     int64
	Received     time.Time `xorm:"-" json:"-"`
	ReceivedUnix int64
}

type RawMailItem struct {
	From      string
	Recipient []string
	Data      string
}

const DataPattern = "Date: (.*)\nFrom: .*\nSubject: (.*)\nTo: (.*)\n\n(.*)."

func (m *Mail) BeforeInsert() {
	m.ReceivedUnix = time.Now().Unix()
}

func (m *Mail) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "sent_unix":
		m.Sent = time.Unix(m.SentUnix, 0).Local()
	case "received_unix":
		m.Received = time.Unix(m.ReceivedUnix, 0).Local()
	}
}

func createMail(e *xorm.Session, raw *RawMailItem) (_ *Mail, _ []MailRecipient, err error) {
	re := regexp.MustCompile(DataPattern)
	match := re.FindStringSubmatch(raw.Data)
	received, _ := time.Parse("Thu, 21 May 2008 05:33:29 -0700", match[0])
	mailItem := &Mail{
		From:     raw.From,
		Received: received,
		Subject:  match[1],
		Body:     match[2],
	}
	if _, err = e.Insert(mailItem); err != nil {
		return nil, nil, err
	}

	var recipients []MailRecipient
	for _, recipient := range raw.Recipient {
		user, err := GetOrCreateUserByName(recipient)
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
