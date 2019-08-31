package models

import (
	"time"

	"github.com/go-xorm/xorm"
)

type Attachment struct {
	Id          int64
	MailId      int64
	FileName    string
	ContentType string
	Data        string `xorm:"TEXT"`

	Received     time.Time `xorm:"created"`
	ReceivedUnix int64     `xorm:"created"`
}

func CreateAttachment(e *xorm.Session, attachment *Attachment) (*Attachment, error) {
	_, err := e.Insert(attachment)
	return attachment, err
}
