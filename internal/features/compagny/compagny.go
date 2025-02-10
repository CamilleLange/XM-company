package compagny

import (
	"fmt"

	"github.com/Aloe-Corporation/mongodb"
	"github.com/google/uuid"
)

// CompagnyFeatures is the interface for the feature compagny.
type CompagnyFeatures interface {
	Create(compagny CompagnyCreateDTO) (uuid.UUID, error)
	ReadByID(uuid uuid.UUID) (*CompagnyPublicDTO, error)
	ReadAll() ([]CompagnyPublicDTO, error)
	Update(uuid uuid.UUID, compagny CompagnyUpdateDTO) error
	Delete(uuid uuid.UUID) error
}

// NewCompagnyFeatures is a factory method to create a new instance of the compagny feature.
func NewCompagnyFeatures(connectorType string, connector any) (CompagnyFeatures, error) {
	var compagnyRepo iCompagnyRepository

	switch connectorType {
	case "mongo":
		db, castable := connector.(*mongodb.Connector)
		if !castable {
			return nil, fmt.Errorf("can't cast connector to mongo database")
		}
		compagnyRepo = newCompagnyRepository(db.Collection("compagny"))
	}

	return newCompagnyController(compagnyRepo), nil
}
