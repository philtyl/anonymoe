package models

import (
	"time"

	"github.com/go-xorm/xorm"
)

type EmbeddedFile struct {
	Id          int64
	MailId      int64
	ContentId   string
	ContentType string
	Data        string `xorm:"TEXT"`

	Received     time.Time `xorm:"created"`
	ReceivedUnix int64     `xorm:"created"`
}

func CreateEmbeddedFile(e *xorm.Session, file *EmbeddedFile) (*EmbeddedFile, error) {
	_, err := e.Insert(file)
	return file, err
}

func GetEmbeddedFile(mailId int64, contentId string) (embeddedFile *EmbeddedFile, has bool, err error) {
	embeddedFile = &EmbeddedFile{
		MailId:    mailId,
		ContentId: contentId,
	}
	has, err = x.Get(embeddedFile)
	return embeddedFile, has, err
}
