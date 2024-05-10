package http

import (
	"errors"
	"fmt"
	"github.com/coddmeistr/quizzify/backend/tests/internal/config"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api"
	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"
	"net/http"
)

const (
	ErrTagHigherThanMaxLimit = "higher_than_max_limit"
	ErrTagLowerThanMinLimit  = "lower_than_min_limit"
	ErrTagNotUnique          = "not_unique"
)

type ValidatorProvider interface {
	Validator() *validator.Validate
}

type Validator struct {
	val ValidatorProvider
	log *zap.Logger
	cfg *config.Config
}

func NewValidator(log *zap.Logger, cfg *config.Config, valProv ValidatorProvider) *Validator {
	return &Validator{
		val: valProv,
		log: log,
		cfg: cfg,
	}
}

func (v *Validator) Validate(w http.ResponseWriter, s interface{}) bool {
	const op = "http.validator.Validate"
	log := v.log.With(zap.String("op", op))

	if err := v.val.Validator().Struct(s); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			log.Error("failed validation", zap.Error(err))
			WriteError(w, ErrFailedValidation)
			return false
		}

		log.Error("failed validation on fields", zap.Error(err))
		var valErrs validator.ValidationErrors
		errors.As(err, &valErrs)
		errs := make([]api.Error, 0, len(valErrs))
		for _, err := range valErrs {
			var errCode string
			switch err.Tag() {
			case ErrTagHigherThanMaxLimit:
				errCode = ErrorCode(ErrMaxLimit)
			case ErrTagLowerThanMinLimit:
				errCode = ErrorCode(ErrMinLimit)
			case ErrTagNotUnique:
				errCode = ErrorCode(ErrUniqueConstraint)
			case "required":
				errCode = ErrorCode(ErrNoRequiredValue)
			default:
				errCode = ErrorCode(ErrFailedValidation)
			}
			errs = append(errs, api.Error{
				Code:    errCode,
				Message: fmt.Sprintf("failed validation on field '%s' with tag '%s'", err.StructField(), err.Tag()),
			})
		}
		WriteErrorManual(w, http.StatusBadRequest, api.Error{
			Code:         ErrorCode(ErrFailedValidation),
			Message:      ErrFailedValidation.Error(),
			NestedErrors: &errs,
		})

		return false
	}

	return true
}
