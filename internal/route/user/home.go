package user

import (
	"github.com/EggMD/EggMD/internal/context"
	log "unknwon.dev/clog/v2"
)

const (
	DASHBOARD = "user/dashboard/dashboard"
)

func Dashboard(c *context.Context) {
	docs, err := c.User.GetDocuments(1, 10)
	if err != nil {
		log.Error("%v", err.Error())
	}
	c.Data["Docs"] = docs

	c.Success(DASHBOARD)
}
