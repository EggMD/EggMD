package route

import (
	"fmt"
	"net/http"

	"gopkg.in/macaron.v1"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/route/user"
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
	c.HTML(http.StatusNotFound, fmt.Sprintf("status/%d", http.StatusNotFound))
}
