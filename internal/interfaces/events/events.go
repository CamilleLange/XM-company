package events

import (
	"fmt"

	"github.com/Aloe-Corporation/logs"
	"github.com/nats-io/nats.go"
)

var log = logs.Get()

type Config struct {
	NatsURL      string `yaml:"nats_url" mapstructure:"nats_url"`
	NatsUsername string `yaml:"nats_username" mapstructure:"nats_username"`
	NatsPassword string `yaml:"nats_password" mapstructure:"nats_password"`
}

func NewNats(config Config) (*nats.Conn, error) {
	conn, err := nats.Connect(config.NatsURL, nats.UserInfo(config.NatsUsername, config.NatsPassword))
	if err != nil {
		return nil, fmt.Errorf("can't open connection to Nats: %w", err)
	}

	return conn, nil
}
