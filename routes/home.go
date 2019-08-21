package routes

import (
	"anonymoe/pkg/context"
)

const (
	HOME = "home"
)

func Home(c *context.Context) {
	c.Success(HOME)
}
