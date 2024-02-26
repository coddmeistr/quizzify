package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

const (
	testsCollection   = "tests"
	resultsCollection = "results"
)

type Storage struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) SaveUserResult(ctx context.Context, result domain.Result) error {
	const op = "mongo.storage.SaveUserResult"

	_, err := s.db.Collection(resultsCollection).InsertOne(ctx, result)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetTests(ctx context.Context) ([]*domain.Test, error) {
	const op = "mongo.storage.GetTests"

	opt := options.Find().SetProjection(getProjectionForAnwers())

	tests := make([]*domain.Test, 0)
	cursor, err := s.db.Collection(testsCollection).Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	if err = cursor.All(ctx, &tests); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tests, nil
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

	res, err := s.db.Collection(testsCollection).DeleteOne(ctx, bson.D{{"_id", testID}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return nil
}

func (s *Storage) GetTestByID(ctx context.Context, testID string, includeAnswers bool) (*domain.Test, error) {
	const op = "mongo.storage.GetTestByID"

	opt := options.FindOne()
	if !includeAnswers {
		opt.SetProjection(getProjectionForAnwers())
	}

	var test domain.Test
	err := s.db.Collection(testsCollection).FindOne(ctx, bson.D{{"_id", testID}}, opt).Decode(&test)
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
	bsonToUpdate := make([]bson.E, 0)
	val := reflect.ValueOf(toUpdate)
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
