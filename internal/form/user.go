package form

type Register struct {
	Name      string `binding:"Required;MaxSize(35)"`
	LoginName string `binding:"Required;AlphaDashDot;MaxSize(35)"`
	Email     string `binding:"Required;Email;MaxSize(254)"`
	Password  string `binding:"Required;MaxSize(255)"`
	Retype    string
}

type SignIn struct {
	Email    string `binding:"Required;MaxSize(254)"`
	Password string `binding:"Required;MaxSize(255)"`
}
