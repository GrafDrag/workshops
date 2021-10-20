package config

type Config struct {
	Rest Rest `toml:"rest"`
	GRPC Rest `toml:"gRPC"`

	Jwt     JWTConfig     `toml:"jwt"`
	DB      DBConfig      `toml:"database"`
	Session SessionConfig `toml:"session"`
}

type GRPC struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
}

type Rest struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
}

type JWTConfig struct {
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

type SessionConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewConfig() *Config {
	return &Config{}
}
