package routes

import (
	"log"
	"strings"

	"anonymoe/models"
	"anonymoe/pkg/context"
	"anonymoe/pkg/setting"
	"github.com/Pallinder/go-randomdata"
)

const (
	INBOX      = "inbox"
	INBOX_NODE = "model/mailinbox"
)

func NewInbox(c *context.Context) {
	c.Redirect(setting.AppURL + "/inbox/" + strings.ToLower(randomdata.SillyName()))
}

func Inbox(c *context.Context) {
	InboxContents(c)
	c.Success(INBOX)
}

func InboxNode(c *context.Context) {
	InboxContents(c)
	c.Success(INBOX_NODE)
}

func InboxContents(c *context.Context) {
	username := strings.ToLower(c.Params(":user"))
	c.Data["User"] = username

	if setting.IsPrivateAccount(username) {
		c.Data["Private"] = true
	} else {
		mail, err := models.GetMail(username)
		log.Printf("Mail for %s: %+v", username, mail)
		if err == nil {
			c.Data["Mail"] = mail
		} else {
			log.Printf("Error loading mail for '%s': %v", username, err)
		}
	}
}
