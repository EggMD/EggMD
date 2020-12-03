package user

import "github.com/EggMD/EggMD/internal/context"

const (
	DASHBOARD = "user/dashboard/dashboard"
)

func Dashboard(c *context.Context) {
	c.Success(DASHBOARD)
}
