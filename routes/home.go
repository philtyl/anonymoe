package routes

import (
	"anonymoe/pkg/context"
	"anonymoe/pkg/setting"
)

const (
	HOME = "home"
)

func Home(c *context.Context) {
	if c.IsLogged() {
		Inbox(c)
		return
	}

	c.Success(HOME)
}

func GlobalInit() {
	setting.NewContext()
}
