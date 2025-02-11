package company

import "github.com/google/uuid"

const (
	Corporation        CompanyType = "Corporation"
	NonProfit          CompanyType = "Non Profit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

// CompanyType is an enum for the company type.
type CompanyType string

// Company is the data model of the feature.
type Company struct {
	UUID            uuid.UUID   `json:"uuid" bson:"uuid"`
	Name            string      `json:"name" bson:"name"`
	Description     string      `json:"description" bson:"description"`
	EmployeesNumber int         `json:"employees_number" bson:"employees_number"`
	Registered      bool        `json:"registered" bson:"registered"`
	Type            CompanyType `json:"type" bson:"type"`
}

// CompanyPublicDTO is the public data model representation.
type CompanyPublicDTO struct {
	Company
}

func (dto *CompanyPublicDTO) ReversePublicDTO() *Company {
	return &dto.Company
}

func FactoryCompanyPublicDTO(company *Company) *CompanyPublicDTO {
	return &CompanyPublicDTO{
		Company: *company,
	}
}

// CompanyCreateDTO is the data model representation to create company.
type CompanyCreateDTO struct {
	Name            string `json:"name" validate:"required,max=15"`
	Description     string `json:"description" validate:"max=3000"`
	EmployeesNumber int    `json:"employees_number" validate:"required"`
	Registered      *bool  `json:"registered" validate:"required"`
	Type            string `json:"type" validate:"required"`
}

func (dto *CompanyCreateDTO) ReverseCreateDTO() *Company {
	return &Company{
		Name:            dto.Name,
		Description:     dto.Description,
		EmployeesNumber: dto.EmployeesNumber,
		Registered:      *dto.Registered,
		Type:            StringToCompanyType(dto.Type),
	}
}

func FactoryCompanyCreateDTO(company *Company) *CompanyCreateDTO {
	return &CompanyCreateDTO{
		Name:            company.Name,
		Description:     company.Description,
		EmployeesNumber: company.EmployeesNumber,
		Registered:      &company.Registered,
		Type:            string(company.Type),
	}
}

// CompanyUpdateDTO is the data model representation to update company.
type CompanyUpdateDTO struct {
	Name            string `json:"name" bson:"name" validate:"required,max=15"`
	Description     string `json:"description" bson:"description" validate:"max=3000"`
	EmployeesNumber int    `json:"employees_number" bson:"employees_number" validate:"required"`
	Registered      *bool  `json:"registered" bson:"registered" validate:"required"`
	Type            string `json:"type" bson:"type" validate:"required"`
}

func (dto *CompanyUpdateDTO) ReverseUpdateDTO() *Company {
	return &Company{
		Name:            dto.Name,
		Description:     dto.Description,
		EmployeesNumber: dto.EmployeesNumber,
		Registered:      *dto.Registered,
		Type:            StringToCompanyType(dto.Type),
	}
}

func FactoryCompanyUpdateDTO(company *Company) *CompanyUpdateDTO {
	return &CompanyUpdateDTO{
		Name:            company.Name,
		Description:     company.Description,
		EmployeesNumber: company.EmployeesNumber,
		Registered:      &company.Registered,
		Type:            string(company.Type),
	}
}

// StringToCompanyType convert a string to a CompanyType
func StringToCompanyType(s string) CompanyType {
	var companyType CompanyType
	switch s {
	case "Corporation":
		companyType = Corporation
	case "Non Profit":
		companyType = NonProfit
	case "Cooperative":
		companyType = Cooperative
	case "Sole Proprietorship":
		companyType = SoleProprietorship

	// This default is not mandatory, thanks to the default value of string.
	// It's here to clarify logic.
	default:
		companyType = ""
	}
	return companyType
}
