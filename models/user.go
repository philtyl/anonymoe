package models

import (
	"sort"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/philtyl/anonymoe/pkg/setting"
	log "gopkg.in/clog.v1"
)

type User struct {
	Id   int64
	Name string `xorm:"UNIQUE NOT NULL"`

	Created     time.Time `xorm:"created"`
	CreatedUnix int64     `xorm:"created"`
}

func getUserMailCount(u *User) (int64, error) {
	return x.Count(&MailRecipient{RecipientId: u.Id})
}

func (u *User) getUserMail() (_ []*Mail, err error) {
	log.Trace("Loading mail items for <%s>", u.Name)
	count, err := getUserMailCount(u)
	if err != nil {
		log.Warn("Unable to count mail items for User <%s>: %v", u.Name, err)
		return
	}
	log.Trace("%d mail items found for <%s>", count, u.Name)

	sess := x.Where("recipient_id=?", u.Id)
	recipients := make([]*MailRecipient, 0, count)
	if err = sess.Find(&recipients); err != nil {
		return
	}

	mailItems := make([]*Mail, 0, count)
	for _, recipient := range recipients {
		mailItem := &Mail{Id: recipient.MailId}
		if has, err := x.Get(mailItem); has {
			mailItems = append(mailItems, mailItem)
		} else if err != nil {
			log.Warn("Error loading [Mail:%d] for User <%s>, ignoring: %v", recipient.MailId, u.Name, err)
		}
	}
	sort.Slice(mailItems, func(i, j int) bool {
		return mailItems[i].ReceivedUnix > mailItems[j].ReceivedUnix
	})
	return mailItems, err
}

func getUserByName(name string) (user *User, has bool, err error) {
	user = &User{
		Name: name,
	}
	has, err = x.Get(user)
	return user, has, err
}

func createUser(e *xorm.Session, userItem *User) (*User, error) {
	_, err := e.Insert(userItem)
	return userItem, err
}

func GetOrCreateUserByName(e *xorm.Session, name string) (user *User, err error) {
	var has bool
	user, has, err = getUserByName(name)
	if err != nil {
		return nil, err
	}

	if !has {
		user, err = createUser(e, user)
		if err != nil {
			return nil, err
		}
		log.Info("Created User: [ID:%d, Name:%s]", user.Id, user.Name)
	}

	return user, nil
}

func GetMail(username string) ([]*Mail, error) {
	user, has, err := getUserByName(username + "@" + setting.Config.AppDomain)
	if err == nil && user != nil && has {
		return user.getUserMail()
	}
	return nil, err
}
