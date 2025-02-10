package datasources

import (
	"fmt"

	"github.com/Aloe-Corporation/mongodb"
)

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
