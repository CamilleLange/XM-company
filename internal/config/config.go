package config

import (
	"fmt"
	"strings"

	"github.com/CamilleLange/XM-compagny/internal/interfaces/datasources"
	"github.com/CamilleLange/XM-compagny/internal/interfaces/http"
	"github.com/spf13/viper"
)

type Config struct {
	Datasources datasources.Config `mapstructure:"datasources"`
	Router      http.Config        `mapstructure:"router"`
}

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
