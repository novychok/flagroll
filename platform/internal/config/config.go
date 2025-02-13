package config

import (
	"os"
	"time"

	"github.com/novychok/flagroll/platform/internal/handler/platformapiv1"
	"github.com/novychok/flagroll/platform/internal/pkg/postgres"
	"github.com/spf13/viper"
)

type JWT struct {
	ExpiresIn        time.Duration `mapstructure:"JWT_EXPIRES_IN"`
	RefreshExpiresIn time.Duration `mapstructure:"JWT_REFRESH_EXPIRES_IN"`
}

type Config struct {
	Postgres      postgres.Config      `mapstructure:",squash"`
	PlatfromAPIV1 platformapiv1.Config `mapstructure:",squash"`
	JWT           JWT                  `mapstructure:",squash"`
}

func New() (*Config, error) {
	viper.AutomaticEnv()

	configPath := "platform.env"

	_, err := os.Stat(configPath)
	if err == nil {
		viper.AddConfigPath("env")
		viper.SetConfigType("env")
		viper.SetConfigName("platform")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetPostgres(cfg *Config) *postgres.Config {
	return &cfg.Postgres
}

func GetPlatfromAPIV1(cfg *Config) *platformapiv1.Config {
	return &cfg.PlatfromAPIV1
}

func GetJWT(cfg *Config) *JWT {
	return &cfg.JWT
}
