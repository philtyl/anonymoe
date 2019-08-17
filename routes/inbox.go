package routes

import (
	"strings"

	"anonymoe/pkg/context"
	"anonymoe/pkg/setting"
	"github.com/Pallinder/go-randomdata"
)

const (
	INBOX = "inbox"
)

func NewInbox(c *context.Context) {
	c.Redirect(setting.AppURL + "/inbox/" + strings.ToLower(randomdata.SillyName()))
}

func Inbox(c *context.Context) {
	c.Data["User"] = c.Params(":user")
	c.Success(INBOX)
}
