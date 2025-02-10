package compagny

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	_ iCompagnyRepository = (*compagnyRepository)(nil)
)

type iCompagnyRepository interface {
	Create(compagny Compagny) (uuid.UUID, error)
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

func (r *compagnyRepository) Create(compagny Compagny) (uuid.UUID, error) {
	if _, err := r.compagny.InsertOne(context.Background(), compagny); err != nil {
		return uuid.Nil, fmt.Errorf("can't create compagny : %w", err)
	}

	return compagny.UUID, nil
}

func (r *compagnyRepository) ReadByID(uuid uuid.UUID) (*Compagny, error) {
	filter := bson.M{
		"uuid": uuid,
	}

	compagny := new(Compagny)
	if err := r.compagny.FindOne(context.TODO(), filter).Decode(compagny); err != nil {
		return nil, fmt.Errorf("can't find document : %w", err)
	}

	return compagny, nil
}

func (r *compagnyRepository) ReadAll() ([]Compagny, error) {
	docs, err := r.compagny.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("can't find all document : %w", err)
	}

	compagnies := make([]Compagny, 0)

	for docs.Next(context.TODO()) {
		compagny := new(Compagny)
		if err := docs.Decode(&compagny); err != nil {
			return nil, fmt.Errorf("can't decode document : %w", err)
		}
		compagnies = append(compagnies, *compagny)
	}

	return compagnies, nil
}

func (r *compagnyRepository) Update(uuid uuid.UUID, compagny CompagnyUpdateDTO) error {
	filter := bson.M{
		"uuid": uuid,
	}

	update := bson.M{
		"$set": bson.M{
			"name":             compagny.Name,
			"description":      compagny.Description,
			"employees_number": compagny.EmployeesNumber,
			"registered":       compagny.Registered,
			"type":             compagny.Type,
		},
	}

	result, err := r.compagny.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("can't update document : %w", err)
	}

	if result.MatchedCount != 1 || result.ModifiedCount != 1 {
		return fmt.Errorf("something unexcepted happened : matched document %d | modified document : %d", result.MatchedCount, result.ModifiedCount)
	}

	return nil
}

func (r *compagnyRepository) Delete(uuid uuid.UUID) error {
	filter := bson.M{
		"uuid": uuid,
	}

	if _, err := r.compagny.DeleteOne(context.TODO(), filter); err != nil {
		return fmt.Errorf("can't delete document : %w", err)
	}

	return nil
}
