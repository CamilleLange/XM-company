package company

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	_ CompagnyFeatures = (*companyController)(nil)
)

// companyController is the controller of the feature.
type companyController struct {
	companyRepository iCompagnyRepository
}

// newCompagnyController is a factory method to create the feature's controller.
func newCompagnyController(companyRepository iCompagnyRepository) *companyController {
	return &companyController{
		companyRepository: companyRepository,
	}
}

// Create use a DTO to create a company.
// DTO fields validation must be done before.
// This check if the company name already exist.
func (c *companyController) Create(company CompagnyCreateDTO) (uuid.UUID, error) {
	companyToCreate := company.ReverseCreateDTO()
	companyToCreate.UUID = uuid.New()

	if !c.isNameValid(company.Name) {
		return uuid.Nil, fmt.Errorf("company name is invalid")
	}

	companyUUID, err := c.companyRepository.Create(*companyToCreate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create company: %w", err)
	}

	return companyUUID, nil
}

// ReadByID read the company with it's uuid.
func (c *companyController) ReadByID(uuid uuid.UUID) (*CompagnyPublicDTO, error) {
	company, err := c.companyRepository.ReadByID(uuid)
	if err != nil {
		return nil, fmt.Errorf("can't read company by id: %w", err)
	}

	return FactoryCompagnyPublicDTO(company), nil
}

// ReadAll read all compagnies.
func (c *companyController) ReadAll() ([]CompagnyPublicDTO, error) {
	company, err := c.companyRepository.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can't read all company: %w", err)
	}

	var publicCompagny []CompagnyPublicDTO
	for _, company := range company {
		publicCompagny = append(publicCompagny, *FactoryCompagnyPublicDTO(&company))
	}

	return publicCompagny, nil
}

// Update use a DTO to update a company.
// DTO fields validation must be done before.
// This check if the company name already exist.
func (c *companyController) Update(uuid uuid.UUID, company CompagnyUpdateDTO) error {
	if !c.isNameValid(company.Name) {
		return fmt.Errorf("company name is invalid")
	}

	if err := c.companyRepository.Update(uuid, company); err != nil {
		return fmt.Errorf("can't update company: %w", err)
	}

	return nil
}

// Delete the company using it's uuid.
func (c *companyController) Delete(uuid uuid.UUID) error {
	if err := c.companyRepository.Delete(uuid); err != nil {
		return fmt.Errorf("can't delete company: %w", err)
	}

	return nil
}

// isNameValid ensure if the provided name can be use in a company.
func (c *companyController) isNameValid(name string) bool {
	previousCompagny, err := c.companyRepository.ReadByName(name)
	if err != nil && !errors.Is(err, ErrCompagnyNotFound) {
		return false
	}

	return previousCompagny == nil
}
