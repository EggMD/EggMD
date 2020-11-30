package conf

type SecurityOpts struct {
	SecretKey string
}

// Security settings
var Security SecurityOpts

type SessionOpts struct {
	CookieName     string
	GCInterval     int64
	MaxLifeTime    int64
	CSRFCookieName string
}

// Session settings
var Session SessionOpts

type ServerOpts struct {
	LandingURL string
	Subpath    string

	HTTPAddr string
	HTTPPort string
}

// Server settings
var Server ServerOpts

type DatabaseOpts struct {
	Host         string
	Name         string
	User         string
	Password     string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}

// Database settings
var Database DatabaseOpts
