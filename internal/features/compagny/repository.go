package compagny

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	_ iCompagnyRepository = (*compagnyRepository)(nil)
)

type iCompagnyRepository interface {
	Create(compagny CompagnyCreateDTO) (uuid.UUID, error)
	ReadByID(uuid uuid.UUID) (*Compagny, error)
	ReadAll() ([]Compagny, error)
	Update(uuid uuid.UUID, compagny CompagnyUpdateDTO) error
	Delete(uuid uuid.UUID) error
}

type compagnyRepository struct {
	compagny *mongo.Collection
}

func newCompagnyRepository(compagny *mongo.Collection) iCompagnyRepository {
	return &compagnyRepository{
		compagny: compagny,
	}
}

func (r *compagnyRepository) Create(compagny CompagnyCreateDTO) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (r *compagnyRepository) ReadByID(uuid uuid.UUID) (*Compagny, error) {
	return nil, nil
}

func (r *compagnyRepository) ReadAll() ([]Compagny, error) {
	return nil, nil
}

func (r *compagnyRepository) Update(uuid uuid.UUID, compagny CompagnyUpdateDTO) error {
	return nil
}

func (r *compagnyRepository) Delete(uuid uuid.UUID) error {
	return nil
}
