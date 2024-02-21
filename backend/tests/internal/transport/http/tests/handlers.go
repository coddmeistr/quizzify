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
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/slice"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	CreateTest(ctx context.Context, test domain.Test) error
	UpdateTest(ctx context.Context, testID string, test domain.Test) error
}

const (
	createTestUrl        = "/tests"
	updateTestPreviewUrl = "/tests/{test_id}/preview"
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
	router.Methods(http.MethodPost).Path(createTestUrl).HandlerFunc(h.CreateTest)
	router.Methods(http.MethodPut).Path(updateTestPreviewUrl).HandlerFunc(h.UpdateTestPreview)
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

	userInfo, ok := user.GetUserInfoFromContext(r.Context())
	if !ok || (userInfo.ID != int(req.Test.UserID) && !slice.Contains(userInfo.Permissions, user.Admin)) {
		log.Error("failed to get user info or no rights", zap.Any("userInfo", userInfo))
		ahttp.WriteError(w, ahttp.ErrForbidden)
		return
	}

	if err := h.srv.CreateTest(r.Context(), *req.Test.ToDomain()); err != nil {
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

	userInfo, ok := user.GetUserInfoFromContext(r.Context())
	if !ok || (userInfo.ID != req.UserID && !slice.Contains(userInfo.Permissions, user.Admin)) {
		log.Error("failed to get user info or no rights", zap.Any("userInfo", userInfo))
		ahttp.WriteError(w, ahttp.ErrForbidden)
		return
	}

	if err := h.srv.UpdateTest(r.Context(), testID, *req.ToDomain()); err != nil {
		if errors.Is(err, testsservice.ErrNotFound) {
			ahttp.WriteError(w, ahttp.ErrNotFound)
			return
		}
		ahttp.WriteError(w, ahttp.ErrInternal)
		return
	}

	ahttp.WriteResponse(w, http.StatusOK, "test preview was updated")
}
