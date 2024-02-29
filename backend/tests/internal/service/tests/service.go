package testsservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/helpers/user"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/storage"
	p "github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/pointer"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/slice"
	"go.uber.org/zap"
)

//go:generate mockery --name Storage
type Storage interface {
	CreateTest(ctx context.Context, test domain.Test) error
	UpdateTest(ctx context.Context, testID string, test domain.Test) error
	DeleteTest(ctx context.Context, testID string) error
	GetTestByID(ctx context.Context, testID string, includeAnswers bool) (*domain.Test, error)
	GetTests(ctx context.Context) ([]*domain.Test, error)
	SaveUserResult(ctx context.Context, result domain.Result) error
}

//go:generate mockery --name Validator
type Validator interface {
	ValidateTest(test domain.Test) error
	ValidateUserAnswers(q domain.Question, a domain.UserAnswerModel) error
}

var (
	ErrNoRights             = errors.New("no rights to perform")
	ErrInvalidTestType      = errors.New("invalid test type")
	ErrFailedTestValidation = errors.New("failed test validation")
	ErrNotFound             = errors.New("not found")
	ErrNoUserAnswer         = errors.New("no user answer")
)

type Service struct {
	cfg        *config.Config
	log        *zap.Logger
	storage    Storage
	validation Validator
}

func New(log *zap.Logger, cfg *config.Config, storage Storage, validator Validator) *Service {
	return &Service{
		cfg:        cfg,
		log:        log,
		storage:    storage,
		validation: validator,
	}
}

func (s *Service) ApplyTest(ctx context.Context, testID string, UserID int, answers map[int]domain.UserAnswerModel) error {
	const op = "service.testsservice.ApplyTest"
	log := s.log.With(zap.String("op", op))
	log.Info("applying test")

	test, err := s.storage.GetTestByID(ctx, testID, true)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("test not found")
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		log.Error("failed to get test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	handleQuestions := func(handler func(domain.Question, domain.UserAnswerModel)) ([]domain.UserAnswerModel, error) {
		ua := make([]domain.UserAnswerModel, 0, len(answers))
		for _, q := range *test.Questions {
			answer, has := answers[q.ID]
			if !has {
				if q.Required {
					log.Warn("no user answer on required question", zap.Int("question_id", q.ID))
					return nil, fmt.Errorf("%s: %w", op, ErrNoUserAnswer)
				}
				log.Warn("no user answer on question", zap.Int("question_id", q.ID))
				handler(*q, domain.UserAnswerModel{QuestionID: 0}) // When no answers, pass empty struct to the handler
				continue
			}
			if err := s.validation.ValidateUserAnswers(*q, answer); err != nil {
				log.Error("failed to validate user answers", zap.Error(err))
				return nil, fmt.Errorf("%s: %w", op, ErrFailedTestValidation)
			}
			handler(*q, answer)
			ua = append(ua, answer)
		}
		return ua, nil
	}

	saveResults := func(r domain.Result) error {
		err := s.storage.SaveUserResult(ctx, r)
		if err != nil {
			log.Error("failed to save user test result", zap.Error(err))
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	switch *test.Type {
	case domain.TestTypeForm:
		ua, err := handleQuestions(func(q domain.Question, a domain.UserAnswerModel) {})
		if err != nil {
			return err
		}

		err = saveResults(domain.Result{
			TestID:      testID,
			UserID:      UserID,
			UserAnswers: ua,
		})
		if err != nil {
			return err
		}
	case domain.TestTypeQuiz:
		ua, err := handleQuestions(func(q domain.Question, a domain.UserAnswerModel) {})
		if err != nil {
			return err
		}

		err = saveResults(domain.Result{
			TestID:      testID,
			UserID:      UserID,
			UserAnswers: ua,
		})
		if err != nil {
			return err
		}
	case domain.TestTypeStrictTest:
		maxPoints := 0
		points := 0
		ua, err := handleQuestions(func(q domain.Question, a domain.UserAnswerModel) {
			if a.QuestionID == 0 {
				maxPoints += *q.Points
				return
			}
			maxPoints += *q.Points
			got := q.ComparePreciseResults(a)
			points += int((float64(got) / 100.0) * float64(*q.Points))
		})
		if err != nil {
			return err
		}

		err = saveResults(domain.Result{
			TestID:      testID,
			UserID:      UserID,
			UserAnswers: ua,
			Percentage:  p.Int(int(float64(points) / float64(maxPoints) * 100.0)),
		})
		if err != nil {
			return err
		}
	case domain.TestTypeTest:
		log.Error("NOT IMPLEMENTED")
		return fmt.Errorf("%s: %w", op, errors.New("SAVING RESULTS FOR TEST OF TYPE TEST NOT IMPLEMENTED"))
	default:
		log.Error("invalid test type", zap.String("type", *test.Type))
		return fmt.Errorf("%s: %w", op, ErrInvalidTestType)
	}

	log.Info("test was applied successfully")
	return nil
}

func (s *Service) GetTests(ctx context.Context) ([]*domain.Test, error) {
	const op = "service.testsservice.GetTests"
	log := s.log.With(zap.String("op", op))
	log.Info("getting tests")

	tests, err := s.storage.GetTests(ctx)
	if err != nil {
		log.Error("failed to get tests", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("tests were gotten successfully")
	return tests, nil
}

func (s *Service) GetTestByID(ctx context.Context, testID string, provideAnswers bool) (*domain.Test, error) {
	const op = "service.testsservice.GetTestByID"
	log := s.log.With(zap.String("op", op))
	log.Info("getting test")

	test, err := s.storage.GetTestByID(ctx, testID, provideAnswers)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("test not found")
			return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		log.Error("failed to get test", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if provideAnswers {
		authUser, ok := user.AuthUserFromContext(ctx)
		if !ok || (authUser.ID != *test.UserID && slice.MaxInt(authUser.Permissions) < user.Admin) {
			log.Error("forbidden action")
			return nil, fmt.Errorf("%s: %w", op, ErrNoRights)
		}
	}

	log.Info("test was gotten successfully")
	return test, nil
}

func (s *Service) DeleteTest(ctx context.Context, testID string) error {
	const op = "service.testsservice.DeleteTest"
	log := s.log.With(zap.String("op", op))
	log.Info("starting test deletion")

	test, err := s.storage.GetTestByID(ctx, testID, false)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("test not found")
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		log.Error("failed to get test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	authUser, ok := user.AuthUserFromContext(ctx)
	if !ok || (authUser.ID != *test.UserID && slice.MaxInt(authUser.Permissions) < user.Admin) {
		log.Error("forbidden action")
		return fmt.Errorf("%s: %w", op, ErrNoRights)
	}

	if err := s.storage.DeleteTest(ctx, testID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("test not found")
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		log.Error("failed to delete test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("test was deleted successfully")
	return nil
}

func (s *Service) CreateTest(ctx context.Context, test domain.Test) error {
	const op = "service.testsservice.CreateTest"
	log := s.log.With(zap.String("op", op))
	log.Info("creating new test")

	if err := s.validation.ValidateTest(test); err != nil {
		log.Error("failed to validate test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, ErrFailedTestValidation)
	}

	if err := s.storage.CreateTest(ctx, test); err != nil {
		s.log.Error("failed to create test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("test was created successfully")
	return nil
}

func (s *Service) UpdateTest(ctx context.Context, testID string, update domain.Test) error {
	const op = "service.testsservice.UpdateTest"
	log := s.log.With(zap.String("op", op))
	log.Info("updating test")

	test, err := s.storage.GetTestByID(ctx, testID, true)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("test not found")
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		log.Error("failed to get test", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	authUser, ok := user.AuthUserFromContext(ctx)
	if !ok || (authUser.ID != *test.UserID && slice.MaxInt(authUser.Permissions) < user.Admin) {
		log.Error("forbidden action")
		return fmt.Errorf("%s: %w", op, ErrNoRights)
	}

	if err := s.storage.UpdateTest(ctx, testID, update); err != nil {
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
