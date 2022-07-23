package user

//
//const (
//	PROFILE_SETTING  = "user/setting/profile"
//	ACCOUNT_SETTING  = "user/setting/account"
//	SECURITY_SETTING = "user/setting/security"
//	DOCUMENT_SETTING = "user/setting/document"
//)
//
//func ProfileSetting(c *context.Context) {
//	c.Title("个人信息")
//	c.PageIs("ProfileSetting")
//	c.Success(PROFILE_SETTING)
//}
//
//func ProfileSettingPost(c *context.Context, f form.ProfileSettings) {
//	c.Title("个人信息")
//	c.PageIs("ProfileSetting")
//
//	if c.HasError() {
//		c.Success(PROFILE_SETTING)
//		return
//	}
//
//	if f.AvatarEmail == "" {
//		f.AvatarEmail = f.Email
//	}
//
//	err := db.Users.UpdateByID(db.UpdateUserOpts{
//		ID:               c.User.ID,
//		Name:             f.Name,
//		Email:            f.Email,
//		KeepEmailPrivate: f.KeepEmailPrivate,
//		AvatarEmail:      f.AvatarEmail,
//	})
//	if err != nil {
//		switch errors.Cause(err) {
//		case db.ErrEmailAlreadyUsed:
//			c.FormErr("Email")
//			c.RenderWithErr("电子邮箱已被绑定", PROFILE_SETTING, &f)
//		default:
//			c.Error(err)
//		}
//	}
//	c.Flash.Success("修改个人信息成功")
//	c.Redirect("/user/settings/profile")
//}
//
//func AccountSetting(c *context.Context) {
//	c.Title("修改密码")
//	c.PageIs("AccountSetting")
//	c.Success(ACCOUNT_SETTING)
//}
//
//func AccountSettingPost(c *context.Context, f form.AccountSettings) {
//	c.Title("修改密码")
//	c.PageIs("AccountSetting")
//
//	if f.NewPassword != f.Retype {
//		c.RenderWithErr("两次密码输入不一致", ACCOUNT_SETTING, &f)
//		return
//	}
//
//	u := db.User{
//		Password: f.Password,
//		Salt:     c.User.Salt,
//	}
//	u.EncodePassword()
//	if c.User.Password != u.Password {
//		c.RenderWithErr("密码不正确", ACCOUNT_SETTING, &f)
//		return
//	}
//
//	c.User.Password = f.NewPassword
//	c.User.EncodePassword()
//
//	err := db.Users.UpdateByID(db.UpdateUserOpts{
//		ID:       c.User.ID,
//		Password: c.User.Password,
//	})
//	if err != nil {
//		c.Error(err)
//		return
//	}
//	c.Flash.Success("修改密码成功")
//	c.Redirect("/user/settings/account")
//}
//
//func SecuritySetting(c *context.Context) {
//	c.Title("安全设置")
//	c.PageIs("SecuritySetting")
//	c.Success(SECURITY_SETTING)
//}
//
//func SecuritySettingPost(c *context.Context, f form.SecuritySettings) {
//	c.Title("修改密码")
//	c.PageIs("AccountSetting")
//
//	c.Redirect("/user/settings/security")
//}
//
//func DocumentSetting(c *context.Context) {
//	c.Title("文档撰写设置")
//	c.PageIs("DocumentSetting")
//	c.Success(DOCUMENT_SETTING)
//}
//
//func DocumentSettingPost(c *context.Context, f form.DocumentSettings) {
//	c.Title("修改密码")
//	c.PageIs("DocumentSetting")
//
//	c.Redirect("/user/settings/document")
//}
