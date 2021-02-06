package user

import (
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/context"
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
