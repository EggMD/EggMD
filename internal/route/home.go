package route

import (
	"net/http"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/route/user"
	"gopkg.in/macaron.v1"
)

func Home(c *context.Context) {
	if c.IsLogged {
		user.Dashboard(c)
		return
	}

	c.Data["PageIsHome"] = true
	c.Success("home")
}

func NotFound(c *macaron.Context) {
	c.Data["Title"] = "页面不存在"
	c.HTML(http.StatusNotFound, "404")
}
