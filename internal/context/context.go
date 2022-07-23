package context

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/dbutil"
)

// Context represents context of a request.
type Context struct {
	flamego.Context

	User     *db.User
	IsLogged bool
	Team     *db.Team
}

func (c *Context) Success(data ...interface{}) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	c.ResponseWriter().WriteHeader(http.StatusOK)

	var d interface{}
	if len(data) == 1 {
		d = data[0]
	} else {
		d = ""
	}

	err := json.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"error": 0,
			"data":  d,
		},
	)
	if err != nil {
		log.Error("Failed to encode: %v", err)
	}
	return nil
}

func (c *Context) ServerError() error {
	return c.Error(http.StatusInternalServerError*100, "Internal server error")
}

func (c *Context) Error(errorCode uint, message string, v ...interface{}) error {
	statusCode := int(errorCode / 100)

	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	c.ResponseWriter().WriteHeader(statusCode)

	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	err := json.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"error": errorCode,
			"msg":   message,
		},
	)
	if err != nil {
		log.Error("Failed to encode: %v", err)
	}
	return nil
}

func (c *Context) Status(statusCode int) {
	c.ResponseWriter().WriteHeader(statusCode)
}

// Contexter initializes a classic context for a request.
func Contexter(gormDB *gorm.DB) flamego.Handler {
	return func(ctx flamego.Context, session session.Session) {
		c := Context{
			Context: ctx,
		}

		c.User = authenticateUser(ctx.Request().Context(), session)
		if c.User != nil {
			c.IsLogged = true
		}

		c.MapTo(gormDB, (*dbutil.Transactor)(nil))
		ctx.Map(c)
	}
}
