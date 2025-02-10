package compagny

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	_ CompagnyFeatures = (*compagnyController)(nil)
)

// compagnyController is the controller of the feature.
type compagnyController struct {
	compagnyRepository iCompagnyRepository
}

// newCompagnyController is a factory method to create the feature's controller.
func newCompagnyController(compagnyRepository iCompagnyRepository) *compagnyController {
	return &compagnyController{
		compagnyRepository: compagnyRepository,
	}
}

// Create use a DTO to create a compagny.
// DTO fields validation must be done before.
// This check if the compagny name already exist.
func (c *compagnyController) Create(compagny CompagnyCreateDTO) (uuid.UUID, error) {
	compagnyToCreate := compagny.ReverseCreateDTO()
	compagnyToCreate.UUID = uuid.New()

	if !c.isNameValid(compagny.Name) {
		return uuid.Nil, fmt.Errorf("compagny name is invalid")
	}

	compagnyUUID, err := c.compagnyRepository.Create(*compagnyToCreate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create compagny: %w", err)
	}

	return compagnyUUID, nil
}

// ReadByID read the compagny with it's uuid.
func (c *compagnyController) ReadByID(uuid uuid.UUID) (*CompagnyPublicDTO, error) {
	compagny, err := c.compagnyRepository.ReadByID(uuid)
	if err != nil {
		return nil, fmt.Errorf("can't read compagny by id: %w", err)
	}

	return FactoryCompagnyPublicDTO(compagny), nil
}

// ReadAll read all compagnies.
func (c *compagnyController) ReadAll() ([]CompagnyPublicDTO, error) {
	compagny, err := c.compagnyRepository.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can't read all compagny: %w", err)
	}

	var publicCompagny []CompagnyPublicDTO
	for _, compagny := range compagny {
		publicCompagny = append(publicCompagny, *FactoryCompagnyPublicDTO(&compagny))
	}

	return publicCompagny, nil
}

// Update use a DTO to update a compagny.
// DTO fields validation must be done before.
// This check if the compagny name already exist.
func (c *compagnyController) Update(uuid uuid.UUID, compagny CompagnyUpdateDTO) error {
	if !c.isNameValid(compagny.Name) {
		return fmt.Errorf("compagny name is invalid")
	}

	if err := c.compagnyRepository.Update(uuid, compagny); err != nil {
		return fmt.Errorf("can't update compagny: %w", err)
	}

	return nil
}

// Delete the compagny using it's uuid.
func (c *compagnyController) Delete(uuid uuid.UUID) error {
	if err := c.compagnyRepository.Delete(uuid); err != nil {
		return fmt.Errorf("can't delete compagny: %w", err)
	}

	return nil
}

// isNameValid ensure if the provided name can be use in a compagny.
func (c *compagnyController) isNameValid(name string) bool {
	previousCompagny, err := c.compagnyRepository.ReadByName(name)
	if err != nil && !errors.Is(err, ErrCompagnyNotFound) {
		return false
	}

	return previousCompagny == nil
}
