package config

import (
	"fmt"
	"strings"

	"github.com/CamilleLange/XM-company/internal/interfaces/datasources"
	"github.com/CamilleLange/XM-company/internal/interfaces/events"
	"github.com/CamilleLange/XM-company/internal/interfaces/http"
	"github.com/spf13/viper"
)

// Config struct of the API.
type Config struct {
	Datasources datasources.Config `mapstructure:"datasources"`
	Event       events.Config      `mapstructure:"events"`
	Router      http.Config        `mapstructure:"router"`
}

// Load API configuration using the provided path.
// This function auto-detect config.yaml file in the provided repository.
// This function override config file with environnement variables.
func Load(path, prefix string) (*Config, error) {
	config := new(Config)

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("can't load API configuration : %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal API configuration : %w", err)
	}

	return config, nil
}
