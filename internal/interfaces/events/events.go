package events

// Config struct of the event package.
type Config struct {
	BrockerAddr string `mapstructure:"brocker_addr"`
	Topic       string `mapstructure:"topic"`
}
