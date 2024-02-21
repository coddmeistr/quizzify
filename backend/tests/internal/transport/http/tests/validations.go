package testshandlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/transport/http"
)

type TestValidator struct {
	val *validator.Validate
	cfg *config.Config
}

func NewTestValidator(cfg *config.Config, val *validator.Validate) *TestValidator {
	return &TestValidator{
		val: val,
		cfg: cfg,
	}
}

func (v *TestValidator) Validator() *validator.Validate {
	return v.val
}

func (v *TestValidator) Register() {
	v.val.RegisterStructValidation(v.ImageStructLevelValidation, Image{})
	v.val.RegisterStructValidation(v.TestStructLevelValidation, Test{})
	v.val.RegisterStructValidation(v.QuestionStructLevelValidation, Question{})
	v.val.RegisterStructValidation(v.UpdateTestPreviewStructLevelValidation, UpdateTestPreviewRequest{})
}

func (v *TestValidator) ImageStructLevelValidation(sl validator.StructLevel) {
	image := sl.Current().Interface().(Image)

	if image.Name == "" {
		sl.ReportError(image, "Name", "Name", "required", "")
	}

	if image.Content == nil || len(image.Content) == 0 {
		sl.ReportError(image, "Content", "Content", "required", "")
	}

	if len(image.Content) > int(v.cfg.Service.Tests.MainImageByteSize) {
		sl.ReportError(image, "Content", "Content", http.ErrTagHigherThanMaxLimit, "")
	}
}

// UpdateTestPreviewStructLevelValidation TODO: add title validation for validations
func (v *TestValidator) UpdateTestPreviewStructLevelValidation(sl validator.StructLevel) {
	test := sl.Current().Interface().(UpdateTestPreviewRequest)

	// Check text values for short and long text
	if test.LongText != "" && len(test.LongText) > int(v.cfg.Service.Tests.LongTextMaxLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagHigherThanMaxLimit, "")
	}

	if test.LongText != "" && len(test.LongText) < int(v.cfg.Service.Tests.LongTextMinLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagLowerThanMinLimit, "")
	}

	if test.ShortText != "" && len(test.ShortText) > int(v.cfg.Service.Tests.ShortTextMaxLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagHigherThanMaxLimit, "")
	}

	if test.ShortText != "" && len(test.ShortText) < int(v.cfg.Service.Tests.ShortTextMinLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagLowerThanMinLimit, "")
	}
}

func (v *TestValidator) TestStructLevelValidation(sl validator.StructLevel) {
	test := sl.Current().Interface().(Test)

	// Check text values for short and long text
	if len(test.LongText) > int(v.cfg.Service.Tests.LongTextMaxLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagHigherThanMaxLimit, "")
	}

	if len(test.LongText) < int(v.cfg.Service.Tests.LongTextMinLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagLowerThanMinLimit, "")
	}

	if len(test.ShortText) > int(v.cfg.Service.Tests.ShortTextMaxLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagHigherThanMaxLimit, "")
	}

	if len(test.ShortText) < int(v.cfg.Service.Tests.ShortTextMinLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagLowerThanMinLimit, "")
	}

	if len(test.Questions) > int(v.cfg.Service.Questions.MaxForCommonUser) {
		sl.ReportError(test.Questions, "Questions", "Questions", http.ErrTagHigherThanMaxLimit, "")
	}

}

func (v *TestValidator) QuestionStructLevelValidation(sl validator.StructLevel) {
	q := sl.Current().Interface().(Question)

	// Check text values for short and long text
	if len(q.LongText) > int(v.cfg.Service.Questions.LongTextMaxLength) {
		sl.ReportError(q.LongText, "LongText", "LongText", http.ErrTagHigherThanMaxLimit, "")
	}

	if len(q.LongText) < int(v.cfg.Service.Questions.LongTextMinLength) {
		sl.ReportError(q.LongText, "LongText", "LongText", http.ErrTagLowerThanMinLimit, "")
	}

	if len(q.ShortText) > int(v.cfg.Service.Questions.ShortTextMaxLength) {
		sl.ReportError(q.ShortText, "ShortText", "ShortText", http.ErrTagHigherThanMaxLimit, "")
	}

	if len(q.ShortText) < int(v.cfg.Service.Questions.ShortTextMinLength) {
		sl.ReportError(q.ShortText, "ShortText", "ShortText", http.ErrTagLowerThanMinLimit, "")
	}

	// TODO: do other validations, mb after some logic would be implemented

}
