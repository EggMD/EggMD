package context

import (
	"context"

	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/db"
)

const (
	UserIDSessionKey = "_UserID"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
}

func Toggle(options *ToggleOptions) flamego.Handler {
	return func(c Context) {
		if options.SignOutRequired && c.IsLogged {
			_ = c.Error(40300, "你已经登录")
			return
		}

		if options.SignInRequired && !c.IsLogged {
			_ = c.Error(40100, "请先登录")
			return
		}
	}
}

// authenticateUser checks if user has signed in.
func authenticateUser(ctx context.Context, sess session.Session) *db.User {
	userIDItf := sess.Get(UserIDSessionKey)
	if userIDItf == nil {
		return nil
	}
	if userID, ok := userIDItf.(uint); ok {
		u, err := db.Users.GetByID(ctx, userID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Error("Failed to get user by ID: %v", err)
			}
			return nil
		}
		return u
	}
	return nil
}
