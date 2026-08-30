package config

type Config struct {
	Port   string
	DBPath string
}

// временно захардкодил
func Load() *Config {
	return &Config{
		Port:   ":8080",
		DBPath: "./habits.db",
	}
}
