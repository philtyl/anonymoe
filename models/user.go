package models

import (
	"sort"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/lunny/log"
	"github.com/philtyl/anonymoe/pkg/setting"
)

type User struct {
	Id   int64
	Name string `xorm:"UNIQUE NOT NULL"`

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
}

func (u *User) BeforeInsert() {
	u.CreatedUnix = time.Now().Unix()
}

func (u *User) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		u.Created = time.Unix(u.CreatedUnix, 0).Local()
	}
}

func (u *User) getUserMail() (mailItems []*Mail, err error) {
	log.Infof("loading mail items for '%s'", u.Name)
	count, err := getUserMailCount(u)
	log.Infof("%d mail items found for '%s'", count, u.Name)
	if err != nil {
		return
	}
	sess := x.Where("recipient_id=?", u.Id)
	recipients := make([]*MailRecipient, 0, count)
	if err = sess.Find(&recipients); err != nil {
		return
	}

	mailItems = make([]*Mail, 0, count)
	for _, recipient := range recipients {
		mailItem := &Mail{Id: recipient.MailId}
		if has, err := x.Get(mailItem); has {
			mailItems = append(mailItems, mailItem)
		} else if err != nil {
			return mailItems, err
		}
	}
	sort.Slice(mailItems, func(i, j int) bool {
		return mailItems[i].ReceivedUnix < mailItems[j].ReceivedUnix
	})
	return
}

func getUserByName(name string) (user *User, has bool, err error) {
	user = &User{
		Name: name,
	}
	has, err = x.Get(user)
	log.Infof("user: %+v, has: %b, err: %+v", user, has, err)
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
		log.Infof("Created User: %+v", user)
	}

	return user, nil
}

func GetMail(username string) ([]*Mail, error) {
	user, has, err := getUserByName(username + "@" + setting.AppDomain)
	if err == nil && user != nil && has {
		return user.getUserMail()
	}
	return nil, err
}
