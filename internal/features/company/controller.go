package company

import (
	"errors"
	"fmt"

	"github.com/CamilleLange/XM-company/internal/interfaces/events"
	"github.com/google/uuid"
)

var (
	_ CompanyFeatures = (*companyController)(nil)
)

// companyController is the controller of the feature.
type companyController struct {
	companyRepository iCompanyRepository
	companyEvent      *events.CompanyEventHandler
}

// newCompanyController is a factory method to create the feature's controller.
func newCompanyController(companyRepository iCompanyRepository, companyEvent *events.CompanyEventHandler) *companyController {
	return &companyController{
		companyRepository: companyRepository,
		companyEvent:      companyEvent,
	}
}

// Create use a DTO to create a company.
// DTO fields validation must be done before.
// This check if the company name already exist.
func (c *companyController) Create(company CompanyCreateDTO) (uuid.UUID, error) {
	companyToCreate := company.ReverseCreateDTO()
	companyToCreate.UUID = uuid.New()

	if !c.isNameValid(company.Name) {
		return uuid.Nil, fmt.Errorf("company name is invalid")
	}

	companyUUID, err := c.companyRepository.Create(*companyToCreate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create company: %w", err)
	}

	e := events.NewEvent(companyToCreate, events.CREATE)
	if err := c.companyEvent.WriteMessage(e); err != nil {
		return uuid.Nil, fmt.Errorf("can't send event : %w", err)
	}

	return companyUUID, nil
}

// ReadByID read the company with it's uuid.
func (c *companyController) ReadByID(uuid uuid.UUID) (*CompanyPublicDTO, error) {
	company, err := c.companyRepository.ReadByID(uuid)
	if err != nil {
		return nil, fmt.Errorf("can't read company by id: %w", err)
	}

	return FactoryCompanyPublicDTO(company), nil
}

// ReadAll read all compagnies.
func (c *companyController) ReadAll() ([]CompanyPublicDTO, error) {
	company, err := c.companyRepository.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can't read all company: %w", err)
	}

	var publicCompany []CompanyPublicDTO
	for _, company := range company {
		publicCompany = append(publicCompany, *FactoryCompanyPublicDTO(&company))
	}

	return publicCompany, nil
}

// Update use a DTO to update a company.
// DTO fields validation must be done before.
// This check if the company name already exist.
func (c *companyController) Update(uuid uuid.UUID, company CompanyUpdateDTO) error {
	if !c.isNameValid(company.Name) {
		return fmt.Errorf("company name is invalid")
	}

	if err := c.companyRepository.Update(uuid, company); err != nil {
		return fmt.Errorf("can't update company: %w", err)
	}

	updatedCompany := company.ReverseUpdateDTO()
	updatedCompany.UUID = uuid

	e := events.NewEvent(updatedCompany, events.UPDATE)
	if err := c.companyEvent.WriteMessage(e); err != nil {
		return fmt.Errorf("can't send event : %w", err)
	}

	return nil
}

// Delete the company using it's uuid.
func (c *companyController) Delete(u uuid.UUID) error {
	if err := c.companyRepository.Delete(u); err != nil {
		return fmt.Errorf("can't delete company: %w", err)
	}

	e := events.NewEvent(struct {
		DeletedUUID uuid.UUID `json:"deleted_uuid"`
	}{
		DeletedUUID: u,
	}, events.UPDATE)

	if err := c.companyEvent.WriteMessage(e); err != nil {
		return fmt.Errorf("can't send event : %w", err)
	}

	return nil
}

// isNameValid ensure if the provided name can be use in a company.
func (c *companyController) isNameValid(name string) bool {
	previousCompany, err := c.companyRepository.ReadByName(name)
	if err != nil && !errors.Is(err, ErrCompanyNotFound) {
		return false
	}

	return previousCompany == nil
}
