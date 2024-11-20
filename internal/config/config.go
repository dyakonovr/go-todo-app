package config

import (
	"os"
	"time"
	"todo-app/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort           = 8000
	defaultHTTPMaxHeaderBytes = 1
	defaultHTTPReadTimeout    = time.Second * 10
	defaultHTTPWriteTimeout   = time.Second * 10
)

type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
}

type HTTPConfig struct {
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

type PostgresConfig struct {
	URI string
}

func Init(configDir string, configFileName string) (*Config, error) {
	setDefaultConfigValues()

	if err := parseConfigFile(configDir, configFileName); err != nil {
		return nil, err
	}

	var config Config

	if err := unmarshal(&config); err != nil {
		return nil, err
	}

	setFromEnv(&config)

	return &config, nil
}

func parseConfigFile(dir string, name string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName(name)
	viper.SetConfigType("yml")

	if error := viper.ReadInConfig(); error != nil {
		return error
	}

	return nil
}

func unmarshal(config *Config) error {
	if err := viper.UnmarshalKey("http", &config.HTTP); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	cfg.Postgres.URI = os.Getenv("POSTGRES_URI")
}

func setDefaultConfigValues() {
	viper.SetDefault("http.host", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderBytes", defaultHTTPMaxHeaderBytes)
	viper.SetDefault("http.readTimeout", defaultHTTPReadTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPWriteTimeout)
}

