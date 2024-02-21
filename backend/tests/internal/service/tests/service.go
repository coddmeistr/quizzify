package testsservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/storage"
	"go.uber.org/zap"
)

type Storage interface {
	CreateTest(ctx context.Context, test domain.Test) error
	UpdateTest(ctx context.Context, testID string, test domain.Test) error
}

var (
	ErrNoRights             = errors.New("forbidden action")
	ErrInvalidTestType      = errors.New("invalid test type")
	ErrFailedTestValidation = errors.New("failed test validation")
	ErrNotFound             = errors.New("not found")
)

type Service struct {
	cfg        *config.Config
	log        *zap.Logger
	storage    Storage
	validation *Validation
}

func New(log *zap.Logger, cfg *config.Config, storage Storage) *Service {
	return &Service{
		cfg:        cfg,
		log:        log,
		storage:    storage,
		validation: NewValidation(cfg, log),
	}
}

func (s *Service) CreateTest(ctx context.Context, test domain.Test) error {
	const op = "service.testsservice.CreateTest"
	log := s.log.With(zap.String("op", op))
	log.Info("creating new test")

	var validated bool
	switch test.Type {
	// Form that is to gather information from one person (or group of people)
	// Should NOT contain correct answers in each question
	case TestTypeForm:
		validated = s.validation.validateForm(test)
	// Quiz that is to gather information from one person or group of people
	// Used to collect a lot of respondents and combine and analyze final result (social quiz's)
	// Should NOT contain correct answers in each question
	case TestTypeQuiz:
		validated = s.validation.validateQuiz(test)
	// Test used to collect some not strict answers and produce some final result
	// This result is not strict and based on the answers provided
	// This test MUST contain correct answers in each question in correct syntax
	case TestTypeTest:
		validated = s.validation.validateTest(test)
	// Strict test used to collect some strict answers and produce some final result
	// This result contains percentage of right answers in all test
	// This test MUST contain correct answers in each question in correct syntax
	case TestTypeStrictTest:
		validated = s.validation.validateStrictTest(test)
	default:
		s.log.Error("invalid test type", zap.String("type", test.Type))
		return fmt.Errorf("%s: %w", op, ErrInvalidTestType)
	}
	if !validated {
		s.log.Error("failed test validation", zap.String("type", test.Type))
		return fmt.Errorf("%s: %w", op, ErrFailedTestValidation)
	}

	if err := s.storage.CreateTest(ctx, test); err != nil {
		s.log.Error("failed to create test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("test was created successfully")
	return nil
}

func (s *Service) UpdateTest(ctx context.Context, testID string, test domain.Test) error {
	const op = "service.testsservice.UpdateTest"
	log := s.log.With(zap.String("op", op))
	log.Info("updating test")

	if err := s.storage.UpdateTest(ctx, testID, test); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			s.log.Error("test not found", zap.Error(err))
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		log.Error("failed to update test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("test was updated successfully")
	return nil
}
