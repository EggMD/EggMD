package user

//
//// Profile 为用户的个人信息页。
//func Profile(c *context.Context) {
//	loginName := c.Params(":name")
//	if loginName == "" {
//		c.Redirect("/user/login")
//		return
//	}
//
//	profileUser, err := db.Users.GetByLoginName(loginName)
//	if err != nil {
//		c.NotFound()
//		return
//	}
//	c.Data["Owner"] = profileUser
//	c.Title(profileUser.Name)
//
//	loggedUID := uint(0)
//	if c.IsLogged {
//		loggedUID = c.User.ID
//	}
//
//	documents, err := db.Documents.GetUserDocuments(&db.UserDocOptions{
//		UserID:    profileUser.ID,
//		LoggedUID: loggedUID,
//		Page:      0,
//		PageSize:  0,
//	})
//	if err != nil {
//		c.Error(err)
//		return
//	}
//	c.Data["Docs"] = documents
//
//	c.Success(PROFILE)
//}
