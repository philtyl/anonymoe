package routes

import (
	"anonymoe/pkg/context"
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
