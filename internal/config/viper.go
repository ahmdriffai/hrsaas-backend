package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	// set config file name
	config.SetConfigName("config")
	config.SetConfigType("yaml")

	// add config patth
	config.AddConfigPath("./config/")    // Docker container path
	config.AddConfigPath("./../config/") // Relative path
	config.AddConfigPath("./")           // Current directory
	config.AddConfigPath("../")          // Parent directory

	// enable environment variable support
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set environment variable prefix
	config.SetEnvPrefix("")

	// Read config file
	err := config.ReadInConfig()
	if err != nil {
		// Check if config path is provided via environment variable
		if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
			config.SetConfigFile(configPath)
			err = config.ReadInConfig()
		}

		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w ", err))
		}
	}

	// Set default values
	setDefaults(config)

	return config
}

// setDefaults sets default values for configuration
func setDefaults(config *viper.Viper) {
	// App defaults
	config.SetDefault("app.name", "Auth Service")

	// Web defaults
	config.SetDefault("web.port", 8080)
	config.SetDefault("web.prefork", false)

	// Log defaults
	config.SetDefault("log.level", 4) // Info level

	// Database defaults
	config.SetDefault("database.host", "localhost")
	config.SetDefault("database.port", 5432)
	config.SetDefault("database.username", "postgres")
	config.SetDefault("database.password", "postgres")
	config.SetDefault("database.name", "auth_service")
	config.SetDefault("database.sslmode", "disable")
	config.SetDefault("database.timezone", "Asia/Jakarta")
	config.SetDefault("database.pool.idle", 10)
	config.SetDefault("database.pool.max", 100)
	config.SetDefault("database.pool.lifetime", 300)

	// JWT defaults
	config.SetDefault("jwt.secret", "supersecret")
	config.SetDefault("jwt.issuer", "bw-auth-service")
	config.SetDefault("jwt.audience", "bw-users")
	config.SetDefault("jwt.expires_in", 30000)

	// Mail defaults
	config.SetDefault("mail.smtp.host", "localhost")
	config.SetDefault("mail.smtp.port", 587)
	config.SetDefault("mail.smtp.username", "")
	config.SetDefault("mail.smtp.password", "")
	config.SetDefault("mail.smtp.from", "noreply@bank-wonosobo.com")

	// Redis defaults
	config.SetDefault("redis.host", "localhost")
	config.SetDefault("redis.port", 6379)
	config.SetDefault("redis.password", "")
	config.SetDefault("redis.db", 0)

	// Monitoring defaults
	config.SetDefault("monitoring.enabled", true)
	config.SetDefault("monitoring.port", 9090)
	config.SetDefault("monitoring.path", "/metrics")
}
