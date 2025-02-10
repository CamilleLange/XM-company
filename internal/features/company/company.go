package company

import (
	"fmt"

	"github.com/Aloe-Corporation/mongodb"
	"github.com/google/uuid"
)

// CompanyFeatures is the interface for the feature company.
type CompanyFeatures interface {
	Create(company CompanyCreateDTO) (uuid.UUID, error)
	ReadByID(uuid uuid.UUID) (*CompanyPublicDTO, error)
	ReadAll() ([]CompanyPublicDTO, error)
	Update(uuid uuid.UUID, company CompanyUpdateDTO) error
	Delete(uuid uuid.UUID) error
}

// NewCompanyFeatures is a factory method to create a new instance of the company feature.
func NewCompanyFeatures(connectorType string, connector any) (CompanyFeatures, error) {
	var companyRepo iCompanyRepository

	switch connectorType {
	case "mongo":
		db, castable := connector.(*mongodb.Connector)
		if !castable {
			return nil, fmt.Errorf("can't cast connector to mongo database")
		}
		companyRepo = newCompanyRepository(db.Collection("company"))
	}

	return newCompanyController(companyRepo), nil
}
