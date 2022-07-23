// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/flamego/session"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/form"
)

var User userRouter

type userRouter struct{}

func (userRouter) SignUp(ctx context.Context, f form.SignUp) error {
	_, err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		NickName:  f.LoginName,
		LoginName: f.LoginName,
		Email:     f.Email,
		Password:  f.Password,
	})
	if err != nil {
		if errors.Is(err, db.ErrLoginNameAlreadyExists) {
			return ctx.Error(40000, "用户名已存在")
		} else if errors.Is(err, db.ErrEmailAlreadyExists) {
			return ctx.Error(40001, "邮箱已被使用")
		}
		log.Error("Failed to create new user: %v", err)
		return ctx.ServerError()
	}
	return ctx.Success("用户注册成功")
}

func (userRouter) SignIn(ctx context.Context, f form.SignIn, session session.Session) error {
	user, err := db.Users.Authenticate(ctx.Request().Context(), f.Email, f.Password)
	if err != nil {
		if errors.Is(err, db.ErrBadCredentials) {
			return ctx.Error(40300, "用户名或密码错误")
		}
		log.Error("Failed to authenticate user: %v", err)
		return ctx.ServerError()
	}

	session.Set(context.UserIDSessionKey, user.ID)
	return ctx.Success(session.ID())
}

func (userRouter) SignOut(ctx context.Context, session session.Session) error {
	session.Flush()
	return ctx.Success("你已成功登出账号")
}

func (userRouter) GetProfile(ctx context.Context) error {
	return ctx.Success(ctx.User)
}
