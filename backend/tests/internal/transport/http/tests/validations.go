package testshandlers

import (
	"github.com/coddmeistr/quizzify/backend/tests/internal/config"
	"github.com/coddmeistr/quizzify/backend/tests/internal/transport/http"
	"github.com/go-playground/validator/v10"
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

	if image.Name == nil || *image.Name == "" {
		sl.ReportError(image, "Name", "Name", "required", "")
	}

	if image.Content == nil || len(*image.Content) == 0 {
		sl.ReportError(image, "Content", "Content", "required", "")
		return
	}

	if len(*image.Content) > int(v.cfg.Service.Tests.MainImageByteSize) {
		sl.ReportError(image, "Content", "Content", http.ErrTagHigherThanMaxLimit, "")
	}
}

func (v *TestValidator) UpdateTestPreviewStructLevelValidation(sl validator.StructLevel) {
	test := sl.Current().Interface().(UpdateTestPreviewRequest)

	if test.Title != nil && len(*test.Title) > 100 {
		sl.ReportError(test.Title, "Title", "Title", http.ErrTagHigherThanMaxLimit, "")
	}

	if test.Title != nil && len(*test.Title) < 3 {
		sl.ReportError(test.Title, "Title", "Title", http.ErrTagLowerThanMinLimit, "")
	}

	if test.LongText != nil && len(*test.LongText) > int(v.cfg.Service.Tests.LongTextMaxLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagHigherThanMaxLimit, "")
	}

	if test.LongText != nil && len(*test.LongText) < int(v.cfg.Service.Tests.LongTextMinLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagLowerThanMinLimit, "")
	}

	if test.ShortText != nil && len(*test.ShortText) < int(v.cfg.Service.Tests.ShortTextMinLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagLowerThanMinLimit, "")
	}

	if test.ShortText != nil && len(*test.ShortText) > int(v.cfg.Service.Tests.ShortTextMaxLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagHigherThanMaxLimit, "")
	}
}

func (v *TestValidator) TestStructLevelValidation(sl validator.StructLevel) {
	test := sl.Current().Interface().(Test)

	if test.LongText != nil && len(*test.LongText) > int(v.cfg.Service.Tests.LongTextMaxLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagHigherThanMaxLimit, "")
	}

	if test.LongText != nil && len(*test.LongText) < int(v.cfg.Service.Tests.LongTextMinLength) {
		sl.ReportError(test.LongText, "LongText", "LongText", http.ErrTagLowerThanMinLimit, "")
	}

	if test.ShortText != nil && len(*test.ShortText) > int(v.cfg.Service.Tests.ShortTextMaxLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagHigherThanMaxLimit, "")
	}

	if test.ShortText != nil && len(*test.ShortText) < int(v.cfg.Service.Tests.ShortTextMinLength) {
		sl.ReportError(test.ShortText, "ShortText", "ShortText", http.ErrTagLowerThanMinLimit, "")
	}

	if len(*test.Questions) > int(v.cfg.Service.Questions.MaxForCommonUser) {
		sl.ReportError(test.Questions, "Questions", "Questions", http.ErrTagHigherThanMaxLimit, "")
	}

	met := make(map[int]struct{})
	for _, v := range *test.Questions {
		if _, ok := met[v.ID]; ok {
			sl.ReportError(test.Questions, "Questions.ID", "Questions.ID", http.ErrTagNotUnique, "")
		}
		met[v.ID] = struct{}{}
	}

}

func (v *TestValidator) QuestionStructLevelValidation(sl validator.StructLevel) {
	q := sl.Current().Interface().(Question)

	// Check text values for short and long text
	if q.LongText != nil && len(*q.LongText) > int(v.cfg.Service.Questions.LongTextMaxLength) {
		sl.ReportError(q.LongText, "LongText", "LongText", http.ErrTagHigherThanMaxLimit, "")
	}

	if q.LongText != nil && len(*q.LongText) < int(v.cfg.Service.Questions.LongTextMinLength) {
		sl.ReportError(q.LongText, "LongText", "LongText", http.ErrTagLowerThanMinLimit, "")
	}

	if q.ShortText != nil && len(*q.ShortText) > int(v.cfg.Service.Questions.ShortTextMaxLength) {
		sl.ReportError(q.ShortText, "ShortText", "ShortText", http.ErrTagHigherThanMaxLimit, "")
	}

	if q.ShortText != nil && len(*q.ShortText) < int(v.cfg.Service.Questions.ShortTextMinLength) {
		sl.ReportError(q.ShortText, "ShortText", "ShortText", http.ErrTagLowerThanMinLimit, "")
	}

}
