package routes

import (
	"strings"

	"anonymoe/models"
	"anonymoe/pkg/context"
	"anonymoe/pkg/setting"
	"github.com/Pallinder/go-randomdata"
	"github.com/lunny/log"
)

const (
	INBOX = "inbox"
)

func NewInbox(c *context.Context) {
	c.Redirect(setting.AppURL + "/inbox/" + strings.ToLower(randomdata.SillyName()))
}

func Inbox(c *context.Context) {
	username := strings.ToLower(c.Params(":user"))
	c.Data["User"] = username

	if setting.IsPrivateAccount(username) {
		c.Data["Private"] = true
	} else {
		mail, err := models.GetMail(username)
		log.Infof("Mail for %s: %+v", username, mail)
		if err == nil {
			c.Data["Mail"] = mail
		} else {
			log.Infof("Error loading mail for '%s': %v", username, err)
		}
	}

	c.Success(INBOX)
}
