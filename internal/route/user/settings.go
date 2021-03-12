package user

import (
	"github.com/pkg/errors"

	"github.com/EggMD/EggMD/internal/context"
	"github.com/EggMD/EggMD/internal/db"
	"github.com/EggMD/EggMD/internal/form"
)

const (
	PROFILE_SETTING = "user/setting/profile"
)

func ProfileSetting(c *context.Context) {
	c.Success(PROFILE_SETTING)
}

func ProfileSettingPost(c *context.Context, f form.Profile) {
	if c.HasError() {
		c.Success(PROFILE_SETTING)
		return
	}

	if f.AvatarEmail == "" {
		f.AvatarEmail = f.Email
	}

	err := db.Users.UpdateByID(db.UpdateUserOpts{
		ID:               c.User.ID,
		Name:             f.Name,
		Email:            f.Email,
		KeepEmailPrivate: f.KeepEmailPrivate,
		AvatarEmail:      f.AvatarEmail,
	})
	if err != nil {
		switch errors.Cause(err) {
		case db.ErrEmailAlreadyUsed:
			c.FormErr("Email")
			c.RenderWithErr("电子邮箱已被绑定", PROFILE_SETTING, &f)
		default:
			c.Error(err)
		}
	}
	c.Flash.Success("修改个人信息成功")
	c.Redirect("/user/settings/profile")
}
