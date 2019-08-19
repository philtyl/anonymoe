package models

import (
	"time"

	"github.com/go-xorm/xorm"
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

func getUserByName(name string) (user *User, err error) {
	user = &User{
		Name: name,
	}
	has, err := x.Get(user)
	if has && err != nil {
		return user, err
	} else {
		return nil, err
	}
}

func createUser(e *xorm.Session, name string) (_ *User, err error) {
	userItem := &User{
		Name: name,
	}
	_, err = e.Insert(userItem)
	return userItem, err
}

func GetOrCreateUserByName(name string) (user *User, err error) {
	user, err = getUserByName(name)
	if err != nil {
		return nil, err
	}

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	if user == nil {
		user, err = createUser(sess, name)
	}

	return user, sess.Commit()
}
