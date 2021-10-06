package store

type Config struct {
	FilePath     string `toml:"file_path"`
	HashFileName string `toml:"hash_file"`
}

func NewConfig() *Config {
	return &Config{
		FilePath:     "./files",
		HashFileName: "hash_list",
	}
}
