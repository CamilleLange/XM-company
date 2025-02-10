package datasources

import "github.com/Aloe-Corporation/mongodb"

type Config struct {
	Mongo mongodb.Conf `mapstructure:"mongo"`
}
