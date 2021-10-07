package httpserver

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`

	JwtSecretKey string `toml:"jwt_secret_key"`
	JwtExpHours  int64  `toml:"jwt_expiration_hours"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":5000",
		LogLevel: "debug",
	}
}
