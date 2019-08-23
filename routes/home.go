package routes

import (
	"github.com/philtyl/anonymoe/pkg/context"
)

const (
	HOME = "home"
)

func Home(c *context.Context) {
	c.Success(HOME)
}
