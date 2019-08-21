package models

import (
	"io/ioutil"
	"net/mail"
	"strings"
	"time"

	"anonymoe/pkg/setting"
	"github.com/go-xorm/xorm"
	"github.com/lunny/log"
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
	Complete  bool
}

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
	r := strings.NewReader(raw.Data)
	m, err := mail.ReadMessage(r)
	if err != nil {
		return
	}

	header := m.Header
	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		return
	}

	received, _ := time.Parse("Thu, 21 May 2008 05:33:29 -0700", header.Get("Date"))
	mailItem := &Mail{
		From:     raw.From,
		Received: received,
		Subject:  header.Get("Subject"),
		Body:     string(body),
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
			log.Infof("Receiving mail for outside address: %s, skipping linkage...", recipient)
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
