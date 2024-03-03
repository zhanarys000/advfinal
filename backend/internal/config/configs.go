package config

type Config struct {
	AppPort       string `env:"APP_PORT"`
	DbHost        string `env:"DB_HOST"`
	DbUser        string `env:"DB_USER"`
	DbPassword    string `env:"DB_PASSWORD"`
	DbName        string `env:"DB_NAME"`
	DbPort        string `env:"DB_PORT"`
	DbSSLMode     string `env:"DB_SSLMODE"`
	RedisAddress  string `env:"REDIS_ADDRESS"`
	Redispassword string `env:"REDIS_PASSWORD"`
	RedisDb       int    `env:"REDIS_DB"`
	EmailFrom     string `env:"EMAIL_FROM"`
	EmailPassword string `env:"EMAIL_PASSWORD"`
	SMTPHost      string `env:"SMTP_HOST"`
	SMTPPort      string `env:"SMTP_PORT"`
}
