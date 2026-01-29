package config

import "os"

type MySQL struct {
	Host              string
	Port              int
	User              string
	Password          string
	DBName            string
	MaxOpenConnection int
	MaxIdleConnection int
}

type JWT struct {
	Secret string
}

type Config struct {
	Port   string
	JWT    JWT
	Domain string

	MySQL MySQL
}

func LoadConfig() Config {
	return Config{
		Port:   getEnv("PORT", ":8000"),
		Domain: getEnv("DOMAIN", "http://localhost:8000"),
		JWT: JWT{
			Secret: getEnv("JWT_SECRET", "secret"),
		},
		MySQL: MySQL{
			Host:              getEnv("MY_SQL_HOST", "localhost"),
			Port:              3306,
			User:              getEnv("MY_SQL_USER", "user"),
			Password:          getEnv("MY_SQL_PASSWORD", "password"),
			DBName:            getEnv("MY_SQL_DB", "products_catalog_db"),
			MaxOpenConnection: 10,
			MaxIdleConnection: 5,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
