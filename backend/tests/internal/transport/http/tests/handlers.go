package testshandlers

import (
	"context"
	"errors"
	"github.com/coddmeistr/quizzify/backend/tests/internal/config"
	"github.com/coddmeistr/quizzify/backend/tests/internal/domain"
	"github.com/coddmeistr/quizzify/backend/tests/internal/helpers/user"
	testsservice "github.com/coddmeistr/quizzify/backend/tests/internal/service/tests"
	ahttp "github.com/coddmeistr/quizzify/backend/tests/internal/transport/http"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/httputil"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/slice"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Service interface {
	CreateTest(ctx context.Context, test domain.Test) (string, error)
	UpdateTest(ctx context.Context, testID string, test domain.Test) error
	DeleteTest(ctx context.Context, testID string) error
	GetTestByID(ctx context.Context, testID string, provideAnswers bool) (*domain.Test, error)
	GetTests(ctx context.Context) ([]*domain.Test, error)
	ApplyTest(ctx context.Context, testID string, UserID int, answers map[int]domain.UserAnswerModel) error
}

const (
	createTestUrl        = "/tests"
	updateTestPreviewUrl = "/tests/{test_id}/preview"
	deleteTestUrl        = "/tests/{test_id}"
	getTestsUrl          = "/tests"
	getTestUrl           = "/tests/{test_id}"
	applyTestUrl         = "/tests/{test_id}/apply"
)

type Handlers struct {
	val *ahttp.Validator
	log *zap.Logger
	srv Service
}

func New(log *zap.Logger, cfg *config.Config, srv Service) *Handlers {
	testVal := NewTestValidator(cfg, validator.New())
	testVal.Register()
	val := ahttp.NewValidator(log, cfg, testVal)

	return &Handlers{
		val: val,
		log: log,
		srv: srv,
	}
}

func (h *Handlers) Register(router *mux.Router) {
	router.Methods(http.MethodGet).Path(getTestsUrl).HandlerFunc(h.GetTests)
	router.Methods(http.MethodGet).Path(getTestUrl).HandlerFunc(h.GetTest)

	auth := router.PathPrefix("").Subrouter()
	auth.Use(
		user.AuthMiddleware(0),
	)
	auth.Methods(http.MethodPost).Path(applyTestUrl).HandlerFunc(h.ApplyTest)
	auth.Methods(http.MethodPost).Path(createTestUrl).HandlerFunc(h.CreateTest)
	auth.Methods(http.MethodPut).Path(updateTestPreviewUrl).HandlerFunc(h.UpdateTestPreview)
	auth.Methods(http.MethodDelete).Path(deleteTestUrl).HandlerFunc(h.DeleteTest)
}

func (h *Handlers) ApplyTest(w http.ResponseWriter, r *http.Request) {
	const op = "tests.handlers.ApplyTest"
	log := h.log.With(zap.String("op", op))

	testID, ok := mux.Vars(r)["test_id"]
	if !ok || testID == "" {
		log.Error("failed to get test id from url path")
		ahttp.WriteErrorMessage(w, ahttp.ErrNoRequiredValue, "no test id in url path")
		return
	}

	var req ApplyTestRequest
	if err := httputil.UnmarshalJSONBody(r.Body, &req); err != nil {
		log.Error("failed to parse body", zap.Error(err))
		ahttp.WriteError(w, ahttp.ErrInvalidJSONBody)
		return
	}

	if ok := h.val.Validate(w, req); !ok {
		log.Error("interrupting request due to failed validation")
		return
	}

	authUser, ok := user.AuthUserFromContext(r.Context())
	if !ok || authUser.ID == 0 {
		log.Error("forbidden action")
		ahttp.WriteError(w, ahttp.ErrForbidden)
		return
	}

	answers := make(map[int]domain.UserAnswerModel)
	for _, a := range req.UserAnswers {
		da := a.ToDomain()
		if _, has := answers[da.QuestionID]; has {
			log.Error("duplicate answers for question", zap.Int("question_id", da.QuestionID))
			ahttp.WriteError(w, ahttp.ErrFailedValidation)
			return
		}
		answers[da.QuestionID] = *da
	}

	if err := h.srv.ApplyTest(r.Context(), testID, authUser.ID, answers); err != nil {
		log.Error("failed to apply test", zap.Error(err))
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusOK, "test was applied")
}

func (h *Handlers) GetTests(w http.ResponseWriter, r *http.Request) {
	const op = "tests.handlers.GetTests"
	log := h.log.With(zap.String("op", op))

	tests, err := h.srv.GetTests(r.Context())
	if err != nil {
		log.Error("failed to get tests", zap.Error(err))
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusOK, tests)
}

func (h *Handlers) GetTest(w http.ResponseWriter, r *http.Request) {
	const op = "tests.handlers.GetTest"
	log := h.log.With(zap.String("op", op))

	testID, ok := mux.Vars(r)["test_id"]
	if !ok || testID == "" {
		log.Error("failed to get test id from url path")
		ahttp.WriteErrorMessage(w, ahttp.ErrNoRequiredValue, "no test id in url path")
		return
	}

	var withAnswers bool
	withAnswers, err := strconv.ParseBool(r.URL.Query().Get("withAnswers"))
	if err != nil {
		withAnswers = false
	}

	test, err := h.srv.GetTestByID(r.Context(), testID, withAnswers)
	if err != nil {
		if errors.Is(err, testsservice.ErrNotFound) {
			log.Error("test not found", zap.Error(err))
			ahttp.WriteError(w, ahttp.ErrNotFound)
			return
		}
		if errors.Is(err, testsservice.ErrNoRights) {
			log.Error("forbidden action", zap.Error(err))
			ahttp.WriteErrorMessage(w, ahttp.ErrForbidden, "no rights to get test with answers")
			return
		}
		log.Error("failed to get test", zap.Error(err))
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusOK, test)
}

func (h *Handlers) DeleteTest(w http.ResponseWriter, r *http.Request) {
	const op = "tests.handlers.DeleteTest"
	log := h.log.With(zap.String("op", op))

	testID, ok := mux.Vars(r)["test_id"]
	if !ok || testID == "" {
		log.Error("failed to get test id from url path")
		ahttp.WriteErrorMessage(w, ahttp.ErrNoRequiredValue, "no test id in url path")
		return
	}

	if err := h.srv.DeleteTest(r.Context(), testID); err != nil {
		if errors.Is(err, testsservice.ErrNotFound) {
			log.Error("test not found", zap.Error(err))
			ahttp.WriteError(w, ahttp.ErrNotFound)
			return
		}
		if errors.Is(err, testsservice.ErrNoRights) {
			log.Error("forbidden action", zap.Error(err))
			ahttp.WriteErrorMessage(w, ahttp.ErrForbidden, "no rights to delete test")
			return
		}
		log.Error("failed to delete test", zap.Error(err))
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusOK, "test was deleted")
}

func (h *Handlers) CreateTest(w http.ResponseWriter, r *http.Request) {
	const op = "tests.handlers.CreateTest"
	log := h.log.With(zap.String("op", op))

	var req CreateTestRequest
	if err := httputil.UnmarshalJSONBody(r.Body, &req); err != nil {
		log.Error("failed to parse body", zap.Error(err))
		ahttp.WriteError(w, ahttp.ErrInvalidJSONBody)
		return
	}

	if ok := h.val.Validate(w, req); !ok {
		log.Error("interrupting request due to failed validation")
		return
	}

	authUser, ok := user.AuthUserFromContext(r.Context())
	if !ok || (authUser.ID != *req.Test.CreatorID && slice.MaxInt(authUser.Permissions) < user.Admin) {
		log.Error("forbidden action")
		ahttp.WriteError(w, ahttp.ErrForbidden)
		return
	}

	id, err := h.srv.CreateTest(r.Context(), *req.Test.ToDomain())
	if err != nil {
		if errors.Is(err, testsservice.ErrInvalidTestType) {
			ahttp.WriteError(w, ahttp.ErrInvalidTestType)
			return
		}
		if errors.Is(err, testsservice.ErrFailedTestValidation) {
			ahttp.WriteError(w, ahttp.ErrInvalidTestStructure)
			return
		}
		if errors.Is(err, testsservice.ErrNoRights) {
			ahttp.WriteErrorMessage(w, ahttp.ErrForbidden, "no rights to create test")
			return
		}
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusCreated, id)
}

func (h *Handlers) UpdateTestPreview(w http.ResponseWriter, r *http.Request) {
	const op = "tests.handlers.UpdateTestPreview"
	log := h.log.With(zap.String("op", op))

	testID, ok := mux.Vars(r)["test_id"]
	if !ok || testID == "" {
		log.Error("failed to get test id from url path")
		ahttp.WriteErrorMessage(w, ahttp.ErrNoRequiredValue, "no test id in url path")
		return
	}

	var req UpdateTestPreviewRequest
	if err := httputil.UnmarshalJSONBody(r.Body, &req); err != nil {
		log.Error("failed to parse body", zap.Error(err))
		ahttp.WriteError(w, ahttp.ErrInvalidJSONBody)
		return
	}

	if ok := h.val.Validate(w, req); !ok {
		log.Error("interrupting request due to failed validation")
		return
	}

	if err := h.srv.UpdateTest(r.Context(), testID, *req.ToDomain()); err != nil {
		if errors.Is(err, testsservice.ErrNotFound) {
			ahttp.WriteError(w, ahttp.ErrNotFound)
			return
		}
		if errors.Is(err, testsservice.ErrNoRights) {
			ahttp.WriteErrorMessage(w, ahttp.ErrForbidden, "no rights to update test")
			return
		}
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusOK, "test preview was updated")
}
