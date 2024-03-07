package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	DatabaseName               = "database.name"
	DatabaseUser               = "database.user"
	DatabasePassword           = "database.password"
	DatabaseHost               = "database.host"
	DatabasePort               = "database.port"
	DatabaseSSL                = "database.ssl"
	DatabaseConnectionLifetime = "database.connection_lifetime"

	ServerPort = "server.port"
)

const (
	DefaultDatabaseName               = "lancelee"
	DefaultDatabaseUser               = "lancelee"
	DefaultDatabasePassword           = "password"
	DefaultDatabaseHost               = "localhost"
	DefaultDatabasePort               = 5432
	DefaultDatabaseSSL                = false
	DefaultDatabaseConnectionLifetime = 5 * time.Minute

	DefaultServerPort = "8080"
)

func InitializeConfig() {
	viper.AutomaticEnv()

	viper.SetDefault(DatabaseName, DefaultDatabaseName)
	viper.SetDefault(DatabaseUser, DefaultDatabaseUser)
	viper.SetDefault(DatabasePassword, DefaultDatabasePassword)
	viper.SetDefault(DatabaseHost, DefaultDatabaseHost)
	viper.SetDefault(DatabasePort, DefaultDatabasePort)
	viper.SetDefault(DatabaseSSL, DefaultDatabaseSSL)
	viper.SetDefault(DatabaseConnectionLifetime, DefaultDatabaseConnectionLifetime)

	viper.SetDefault(ServerPort, DefaultServerPort)

	viper.SetEnvKeyReplacer(defaultReplacer())
}

func defaultReplacer() *strings.Replacer {
	return strings.NewReplacer(".", "_")
}

func FormatDBConfig() string {
	ssl := "disable"
	if viper.GetBool(DatabaseSSL) {
		ssl = "enable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString(DatabaseUser),
		viper.GetString(DatabasePassword),
		viper.GetString(DatabaseHost),
		viper.GetInt(DatabasePort),
		viper.GetString(DatabaseName),
		ssl,
	)
}
