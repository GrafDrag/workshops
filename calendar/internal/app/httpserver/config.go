package httpserver

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`

	Jwt     jwtConfig     `toml:"jwt"`
	DB      DBConfig      `toml:"database"`
	Session sessionConfig `toml:"session"`
}

type jwtConfig struct {
	JwtSecretKey string `toml:"jwt_secret_key"`
	JwtExpHours  int64  `toml:"jwt_expiration_hours"`
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

type sessionConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":5000",
		LogLevel: "debug",
	}
}
