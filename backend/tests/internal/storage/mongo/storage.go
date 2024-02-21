package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

const (
	testsCollection = "tests"
)

type Storage struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateTest(ctx context.Context, test domain.Test) error {
	const op = "mongo.storage.CreateTest"

	_, err := s.db.Collection(testsCollection).InsertOne(ctx, test)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteTest(ctx context.Context, testID string) error {
	const op = "mongo.storage.DeleteTest"

	_, err := s.db.Collection(testsCollection).DeleteOne(ctx, domain.Test{ID: testID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetTestByID(ctx context.Context, testID string) (*domain.Test, error) {
	const op = "mongo.storage.GetTestByID"

	var test domain.Test
	err := s.db.Collection(testsCollection).FindOne(ctx, domain.Test{ID: testID}).Decode(&test)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &test, nil
}

func (s *Storage) UpdateTest(ctx context.Context, testID string, toUpdate domain.Test) error {
	const op = "mongo.storage.UpdateTest"

	// Prepare values to update only for non-nil fields
	val := reflect.ValueOf(toUpdate)
	bsonToUpdate := make([]bson.E, 0)
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Kind() == reflect.Pointer && val.Field(i).IsNil() {
			continue
		}

		if val.Field(i).IsZero() {
			continue
		}

		bsonToUpdate = append(bsonToUpdate, bson.E{Key: val.Type().Field(i).Tag.Get("bson"), Value: val.Field(i).Interface()})
	}
	update := bson.D{{"$set", bson.D(bsonToUpdate)}}

	res, err := s.db.Collection(testsCollection).UpdateOne(ctx, bson.D{{"_id", testID}}, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return nil
}
