package models

type User struct {
	Name string `xorm:"UNIQUE NOT NULL"`
}
