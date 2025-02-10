package compagny

import "github.com/google/uuid"

const (
	Corporation        CompagnyType = "Corporation"
	NonProfit          CompagnyType = "Non Profit"
	Cooperative        CompagnyType = "Cooperative"
	SoleProprietorship CompagnyType = "Sole Proprietorship"
)

type CompagnyType string

type Compagny struct {
	UUID            uuid.UUID    `json:"uuid" bson:"uuid"`
	Name            string       `json:"name" bson:"name"`
	Description     string       `json:"description" bson:"description"`
	EmployeesNumber int          `json:"employees_number" bson:"employees_number"`
	Registered      bool         `json:"registered" bson:"registered"`
	Type            CompagnyType `json:"type" bson:"type"`
}

type CompagnyPublicDTO struct {
	Compagny
}

func (dto *CompagnyPublicDTO) ReversePublicDTO() *Compagny {
	return &dto.Compagny
}

func FactoryCompagnyPublicDTO(compagny *Compagny) *CompagnyPublicDTO {
	return &CompagnyPublicDTO{
		Compagny: *compagny,
	}
}

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

func FactoryCompagnyCreateDTO(compagny *Compagny) *CompagnyCreateDTO {
	return &CompagnyCreateDTO{
		Name:            compagny.Name,
		Description:     compagny.Description,
		EmployeesNumber: compagny.EmployeesNumber,
		Registered:      &compagny.Registered,
		Type:            string(compagny.Type),
	}
}

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

func FactoryCompagnyUpdateDTO(compagny *Compagny) *CompagnyUpdateDTO {
	return &CompagnyUpdateDTO{
		Name:            compagny.Name,
		Description:     compagny.Description,
		EmployeesNumber: compagny.EmployeesNumber,
		Registered:      &compagny.Registered,
		Type:            string(compagny.Type),
	}
}

func StringToCompagnyType(s string) CompagnyType {
	var compagnyType CompagnyType
	switch s {
	case "Corporation":
		compagnyType = Corporation
	case "Non Profit":
		compagnyType = NonProfit
	case "Cooperative":
		compagnyType = Cooperative
	case "Sole Proprietorship":
		compagnyType = SoleProprietorship
	}
	return compagnyType
}
