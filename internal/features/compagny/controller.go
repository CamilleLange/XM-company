package compagny

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	_ CompagnyFeatures = (*compagnyController)(nil)
)

type compagnyController struct {
	compagnyRepository iCompagnyRepository
}

func newCompagnyController(compagnyRepository iCompagnyRepository) *compagnyController {
	return &compagnyController{
		compagnyRepository: compagnyRepository,
	}
}

func (c *compagnyController) Create(compagny CompagnyCreateDTO) (uuid.UUID, error) {
	compagnyToCreate := compagny.ReverseCreateDTO()
	compagnyToCreate.UUID = uuid.New()

	compagnyUUID, err := c.compagnyRepository.Create(*compagnyToCreate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create compagny: %w", err)
	}

	return compagnyUUID, nil
}

func (c *compagnyController) ReadByID(uuid uuid.UUID) (*CompagnyPublicDTO, error) {
	compagny, err := c.compagnyRepository.ReadByID(uuid)
	if err != nil {
		return nil, fmt.Errorf("can't read compagny by id: %w", err)
	}

	return FactoryCompagnyPublicDTO(compagny), nil
}

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

func (c *compagnyController) Update(uuid uuid.UUID, compagny CompagnyUpdateDTO) error {
	if err := c.compagnyRepository.Update(uuid, compagny); err != nil {
		return fmt.Errorf("can't update compagny: %w", err)
	}

	return nil
}

func (c *compagnyController) Delete(uuid uuid.UUID) error {
	if err := c.compagnyRepository.Delete(uuid); err != nil {
		return fmt.Errorf("can't delete compagny: %w", err)
	}

	return nil
}
