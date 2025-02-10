package datasources

import (
	"fmt"

	"github.com/Aloe-Corporation/mongodb"
)

// NewMongoDB use the config to create a *mongodb.Connector ready to be used.
func NewMongoDB(config mongodb.Conf) (*mongodb.Connector, error) {
	connector, err := mongodb.FactoryConnector(config)
	if err != nil {
		return nil, fmt.Errorf("fail to init MongoDB connector: %w", err)
	}

	err = connector.TryConnection()
	if err != nil {
		return nil, fmt.Errorf("fail to ping MongoDB : %w", err)
	}

	return connector, nil
}
