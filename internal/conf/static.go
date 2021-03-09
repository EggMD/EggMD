package conf

// Security 包含安全相关的配置。
var Security SecurityOpts

type SecurityOpts struct {
	SecretKey string
}

// Session 包含网站 Session 会话凭证相关的配置。
var Session SessionOpts

type SessionOpts struct {
	CookieName     string
	GCInterval     int64
	MaxLifeTime    int64
	CSRFCookieName string
}

// Server 包含 Web 服务器相关配置。
var Server ServerOpts

type ServerOpts struct {
	SubPath string

	HTTPAddr string
	HTTPPort string
}

// Database 包含数据库相关配置。
var Database DatabaseOpts

type DatabaseOpts struct {
	Host         string
	Name         string
	User         string
	Password     string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}
