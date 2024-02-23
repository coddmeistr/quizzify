package testshandlers

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/helpers/user"
	testsservice "github.com/maxik12233/quizzify-online-tests/backend/tests/internal/service/tests"
	ahttp "github.com/maxik12233/quizzify-online-tests/backend/tests/internal/transport/http"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/httputil"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	CreateTest(ctx context.Context, test domain.Test) error
	UpdateTest(ctx context.Context, testID string, test domain.Test) error
	DeleteTest(ctx context.Context, testID string) error
	GetTestByID(ctx context.Context, testID string) (*domain.Test, error)
	GetTests(ctx context.Context) ([]*domain.Test, error)
}

const (
	createTestUrl        = "/tests"
	updateTestPreviewUrl = "/tests/{test_id}/preview"
	deleteTestUrl        = "/tests/{test_id}"
	getTestsUrl          = "/tests"
	getTestUrl           = "/tests/{test_id}"
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
		user.AuthMiddleware(user.Admin),
	)
	auth.Methods(http.MethodPost).Path(createTestUrl).HandlerFunc(h.CreateTest)
	auth.Methods(http.MethodPut).Path(updateTestPreviewUrl).HandlerFunc(h.UpdateTestPreview)
	auth.Methods(http.MethodDelete).Path(deleteTestUrl).HandlerFunc(h.DeleteTest)
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

	test, err := h.srv.GetTestByID(r.Context(), testID)
	if err != nil {
		if errors.Is(err, testsservice.ErrNotFound) {
			log.Error("test not found", zap.Error(err))
			ahttp.WriteError(w, ahttp.ErrNotFound)
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

	userInfo, ok := user.SubjectUserFromContext(r.Context())
	if !ok {
		log.Error("failed to get user info", zap.Any("user_info", userInfo))
		ahttp.WriteError(w, ahttp.ErrForbidden)
		return
	}

	dt := *req.Test.ToDomain()
	dt.UserID = &userInfo.ID

	if err := h.srv.CreateTest(r.Context(), dt); err != nil {
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

	ahttp.WriteResponse(w, http.StatusCreated, "test was created")
}

// UpdateTestPreview TODO: Refactor subject id logic, handler and service dont validating that test with given testID belongs to subject and auth user
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

	//userInfo, ok := user.SubjectUserFromContext(r.Context())
	//if !ok {
	//	log.Error("failed to get user info", zap.Any("user_info", userInfo))
	//	ahttp.WriteError(w, ahttp.ErrForbidden)
	//	return
	//}

	dt := *req.ToDomain()

	if err := h.srv.UpdateTest(r.Context(), testID, dt); err != nil {
		if errors.Is(err, testsservice.ErrNotFound) {
			ahttp.WriteError(w, ahttp.ErrNotFound)
			return
		}
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusOK, "test preview was updated")
}
