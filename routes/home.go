package routes

import (
	"anonymoe/pkg/context"
	"anonymoe/pkg/setting"
)

const (
	HOME = "home"
)

func Home(c *context.Context) {
	c.Success(HOME)
}

func GlobalInit() {
	setting.NewContext()
}
