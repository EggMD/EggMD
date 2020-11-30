package route

import (
	"github.com/EggMD/EggMD/internal/context"
)

func Home(c *context.Context) {
	if c.IsLogged {
		// TODO
		return
	}

	c.Data["PageIsHome"] = true
	c.Success("home")
}
