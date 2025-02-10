package company

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	_ iCompagnyRepository = (*companyRepository)(nil)

	// ErrCompagnyNotFound indicates that the company was not found in the database.
	ErrCompagnyNotFound = errors.New("company not found")
)

// CompagnyFeatures is the interface for the company's repository.
type iCompagnyRepository interface {
	Create(company Compagny) (uuid.UUID, error)
	ReadByID(uuid uuid.UUID) (*Compagny, error)
	ReadByName(name string) (*Compagny, error)
	ReadAll() ([]Compagny, error)
	Update(uuid uuid.UUID, company CompagnyUpdateDTO) error
	Delete(uuid uuid.UUID) error
}

// companyRepository hold the *mongo.Collection, and do the database requests.
type companyRepository struct {
	company *mongo.Collection
}

// newCompagnyRepository is a factory method to create a new iCompagnyRepository.
func newCompagnyRepository(company *mongo.Collection) iCompagnyRepository {
	return &companyRepository{
		company: company,
	}
}

// Create the company in the database.
// This function perform no checks of any kind.
func (r *companyRepository) Create(company Compagny) (uuid.UUID, error) {
	if _, err := r.company.InsertOne(context.Background(), company); err != nil {
		return uuid.Nil, fmt.Errorf("can't create company : %w", err)
	}

	return company.UUID, nil
}

// ReadByID the requested company from the database.
func (r *companyRepository) ReadByID(uuid uuid.UUID) (*Compagny, error) {
	filter := bson.M{
		"uuid": uuid,
	}

	company := new(Compagny)
	if err := r.company.FindOne(context.TODO(), filter).Decode(company); err != nil {
		return nil, fmt.Errorf("can't find document : %w", err)
	}

	return company, nil
}

// ReadByName the requested company from the database.
func (r *companyRepository) ReadByName(name string) (*Compagny, error) {
	filter := bson.M{
		"name": name,
	}

	company := new(Compagny)
	if err := r.company.FindOne(context.TODO(), filter).Decode(company); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrCompagnyNotFound
		}
		return nil, fmt.Errorf("can't find document : %w", err)
	}

	return company, nil
}

// ReadAll compagnies from the database.
func (r *companyRepository) ReadAll() ([]Compagny, error) {
	docs, err := r.company.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("can't find all document : %w", err)
	}

	compagnies := make([]Compagny, 0)

	for docs.Next(context.TODO()) {
		company := new(Compagny)
		if err := docs.Decode(&company); err != nil {
			return nil, fmt.Errorf("can't decode document : %w", err)
		}
		compagnies = append(compagnies, *company)
	}

	return compagnies, nil
}

// Update the requested company in the database.
func (r *companyRepository) Update(uuid uuid.UUID, company CompagnyUpdateDTO) error {
	filter := bson.M{
		"uuid": uuid,
	}

	update := bson.M{
		"$set": bson.M{
			"name":             company.Name,
			"description":      company.Description,
			"employees_number": company.EmployeesNumber,
			"registered":       company.Registered,
			"type":             company.Type,
		},
	}

	result, err := r.company.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("can't update document : %w", err)
	}

	if result.MatchedCount != 1 || result.ModifiedCount != 1 {
		return fmt.Errorf("something unexcepted happened : matched document %d | modified document : %d", result.MatchedCount, result.ModifiedCount)
	}

	return nil
}

// Delete the requested company from the database.
func (r *companyRepository) Delete(uuid uuid.UUID) error {
	filter := bson.M{
		"uuid": uuid,
	}

	if _, err := r.company.DeleteOne(context.TODO(), filter); err != nil {
		return fmt.Errorf("can't delete document : %w", err)
	}

	return nil
}
