package routes

import (
	"anonymoe/pkg/context"
)

const (
	INBOX = "inbox"
)

func Inbox(c *context.Context) {
	c.Success(INBOX)
}
