package conf

var Server struct {
	AppVersion string `toml:"-"`

	HTTPAddr string
	HTTPPort string
}

var Database struct {
	DSN          string `toml:"-"` // DSN is set when connect to the database.
	Host         string
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}
