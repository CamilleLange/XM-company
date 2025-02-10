package company

import (
	"fmt"

	"github.com/Aloe-Corporation/mongodb"
	"github.com/google/uuid"
)

// CompagnyFeatures is the interface for the feature company.
type CompagnyFeatures interface {
	Create(company CompagnyCreateDTO) (uuid.UUID, error)
	ReadByID(uuid uuid.UUID) (*CompagnyPublicDTO, error)
	ReadAll() ([]CompagnyPublicDTO, error)
	Update(uuid uuid.UUID, company CompagnyUpdateDTO) error
	Delete(uuid uuid.UUID) error
}

// NewCompagnyFeatures is a factory method to create a new instance of the company feature.
func NewCompagnyFeatures(connectorType string, connector any) (CompagnyFeatures, error) {
	var companyRepo iCompagnyRepository

	switch connectorType {
	case "mongo":
		db, castable := connector.(*mongodb.Connector)
		if !castable {
			return nil, fmt.Errorf("can't cast connector to mongo database")
		}
		companyRepo = newCompagnyRepository(db.Collection("company"))
	}

	return newCompagnyController(companyRepo), nil
}
