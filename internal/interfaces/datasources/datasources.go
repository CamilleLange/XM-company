package datasources

import "github.com/Aloe-Corporation/mongodb"

// Config of the database package.
type Config struct {
	Mongo mongodb.Conf `mapstructure:"mongo"`
}
