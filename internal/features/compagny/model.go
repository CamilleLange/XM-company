package company

import "github.com/google/uuid"

const (
	Corporation        CompagnyType = "Corporation"
	NonProfit          CompagnyType = "Non Profit"
	Cooperative        CompagnyType = "Cooperative"
	SoleProprietorship CompagnyType = "Sole Proprietorship"
)

// CompagnyType is an enum for the company type.
type CompagnyType string

// Compagny is the data model of the feature.
type Compagny struct {
	UUID            uuid.UUID    `json:"uuid" bson:"uuid"`
	Name            string       `json:"name" bson:"name"`
	Description     string       `json:"description" bson:"description"`
	EmployeesNumber int          `json:"employees_number" bson:"employees_number"`
	Registered      bool         `json:"registered" bson:"registered"`
	Type            CompagnyType `json:"type" bson:"type"`
}

// CompagnyPublicDTO is the public data model representation.
type CompagnyPublicDTO struct {
	Compagny
}

func (dto *CompagnyPublicDTO) ReversePublicDTO() *Compagny {
	return &dto.Compagny
}

func FactoryCompagnyPublicDTO(company *Compagny) *CompagnyPublicDTO {
	return &CompagnyPublicDTO{
		Compagny: *company,
	}
}

// CompagnyCreateDTO is the data model representation to create company.
type CompagnyCreateDTO struct {
	Name            string `json:"name" validate:"required,max=15"`
	Description     string `json:"description" validate:"max=3000"`
	EmployeesNumber int    `json:"employees_number" validate:"required"`
	Registered      *bool  `json:"registered" validate:"required"`
	Type            string `json:"type" validate:"required"`
}

func (dto *CompagnyCreateDTO) ReverseCreateDTO() *Compagny {
	return &Compagny{
		Name:            dto.Name,
		Description:     dto.Description,
		EmployeesNumber: dto.EmployeesNumber,
		Registered:      *dto.Registered,
		Type:            StringToCompagnyType(dto.Type),
	}
}

func FactoryCompagnyCreateDTO(company *Compagny) *CompagnyCreateDTO {
	return &CompagnyCreateDTO{
		Name:            company.Name,
		Description:     company.Description,
		EmployeesNumber: company.EmployeesNumber,
		Registered:      &company.Registered,
		Type:            string(company.Type),
	}
}

// CompagnyUpdateDTO is the data model representation to update company.
type CompagnyUpdateDTO struct {
	Name            string `json:"name" bson:"name" validate:"required,max=15"`
	Description     string `json:"description" bson:"description" validate:"max=3000"`
	EmployeesNumber int    `json:"employees_number" bson:"employees_number" validate:"required"`
	Registered      *bool  `json:"registered" bson:"registered" validate:"required"`
	Type            string `json:"type" bson:"type" validate:"required"`
}

func (dto *CompagnyUpdateDTO) ReverseUpdateDTO() *Compagny {
	return &Compagny{
		Name:            dto.Name,
		Description:     dto.Description,
		EmployeesNumber: dto.EmployeesNumber,
		Registered:      *dto.Registered,
		Type:            StringToCompagnyType(dto.Type),
	}
}

func FactoryCompagnyUpdateDTO(company *Compagny) *CompagnyUpdateDTO {
	return &CompagnyUpdateDTO{
		Name:            company.Name,
		Description:     company.Description,
		EmployeesNumber: company.EmployeesNumber,
		Registered:      &company.Registered,
		Type:            string(company.Type),
	}
}

// StringToCompagnyType convert a string to a CompagnyType
func StringToCompagnyType(s string) CompagnyType {
	var companyType CompagnyType
	switch s {
	case "Corporation":
		companyType = Corporation
	case "Non Profit":
		companyType = NonProfit
	case "Cooperative":
		companyType = Cooperative
	case "Sole Proprietorship":
		companyType = SoleProprietorship
	}
	return companyType
}
