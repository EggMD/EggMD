package form

type SignUp struct {
	LoginName string `valid:"required" label:"用户名"`
	Email     string `valid:"required;email;maxlen:100" label:"电子邮箱"`
	Password  string `valid:"required" label:"密码"`
}

type SignIn struct {
	Email    string `valid:"required;email;maxlen:100" label:"电子邮箱"`
	Password string `valid:"required" label:"密码"`
}
