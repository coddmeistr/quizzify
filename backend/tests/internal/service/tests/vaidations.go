package testsservice

import (
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/slice"
	"go.uber.org/zap"
)

type Validation struct {
	cfg *config.Config
	log *zap.Logger
}

func NewValidation(cfg *config.Config, log *zap.Logger) *Validation {
	return &Validation{
		cfg: cfg,
		log: log,
	}
}

func (val *Validation) validateUserAnswers(q domain.Question, a domain.UserAnswerModel) bool {
	const op = "testsservice.validation.validateUserAnswers"
	log := val.log.With(zap.String("op", op), zap.String("qtype", *q.Type))

	switch *q.Type {
	case domain.QuestionTypeSingleChoice:
		if a.ChosenID == nil {
			log.Error("no chosen id")
			return false
		}
		for _, v := range *q.Variants.SingleChoice.Fields {
			if v.FieldID == *a.ChosenID {
				return true
			}
		}
		log.Error("chosen id is not in variants")
		return false
	case domain.QuestionTypeMultipleChoice:
		if a.ChosenIDs == nil {
			log.Error("no chosen ids")
			return false
		}
		if slice.ContainsRepeated(*a.ChosenIDs) {
			log.Error("repeated chosen ids")
			return false
		}
		met := make(map[int]struct{})
		for _, v := range *q.Variants.MultipleChoice.Fields {
			met[v.FieldID] = struct{}{}
		}
		for _, v := range *a.ChosenIDs {
			if _, ok := met[v]; !ok {
				log.Error("chosen id is not in variants")
				return false
			}
		}
		return true
	case domain.QuestionTypeManualInput:
		if a.WritedText == nil {
			log.Error("no writed text")
			return false
		}
		return true
	default:
		log.Error("invalid question type")
		return false
	}
}

func (val *Validation) validateForm(test domain.Test) bool {
	for _, q := range *test.Questions {
		if !val.validateQuestion(*q, false) {
			return false
		}
	}

	return true
}

func (val *Validation) validateQuiz(test domain.Test) bool {
	for _, q := range *test.Questions {
		if !val.validateQuestion(*q, false) {
			return false
		}
	}

	return true
}

func (val *Validation) validateTest(test domain.Test) bool {
	for _, q := range *test.Questions {
		if !val.validateQuestion(*q, true) {
			return false
		}
	}

	return true
}

func (val *Validation) validateStrictTest(test domain.Test) bool {
	const op = "testsservice.validation.validateStrictTest"
	log := val.log.With(zap.String("op", op))

	for _, q := range *test.Questions {
		if q.Points == nil || *q.Points <= 0 {
			log.Error("no points")
			return false
		}
		if !val.validateQuestion(*q, true) {
			return false
		}
	}

	return true
}

func (val *Validation) validateQuestion(q domain.Question, checkAnswers bool) bool {
	const op = "testsservice.validation.validateQuestion"
	log := val.log.With(zap.String("op", op), zap.String("qtype", *q.Type))

	switch *q.Type {
	case domain.QuestionTypeSingleChoice:
		var (
			v = q.Variants.SingleChoice
		)
		if v == nil {
			log.Error("no single choice structure")
			return false
		}
		if v.Fields == nil {
			log.Error("no single choice variants slice")
			return false
		}
		if len(*v.Fields) <= 0 {
			log.Error("zero variants")
			return false
		}
		if !val.validateFields(*v.Fields) {
			log.Error("failed to validate fields")
			return false
		}

		if checkAnswers {
			var a = q.Answers
			if !val.validateAnswers(a) {
				log.Error("coulnd't validate answers")
				return false
			}
			if a.CorrectID == nil {
				log.Error("no required CorrectID value for single choice model")
				return false
			}

			for _, v := range *v.Fields {
				if v.FieldID == *a.CorrectID {
					return true
				}
			}

			log.Error("CorrectID doesn't pointing to some of the FieldID in fields slice")
			return false
		}
	case domain.QuestionTypeMultipleChoice:
		var (
			v = q.Variants.MultipleChoice
		)
		if v == nil {
			log.Error("no multiple choice structure")
			return false
		}
		if v.Fields == nil {
			log.Error("no multiple choice variants")
			return false
		}
		if v.MaxChoices == nil {
			log.Error("no max choices")
			return false
		}
		if !(len(*v.Fields) > 0 &&
			*v.MaxChoices > 0) {
			log.Error("no fields or zero max choices", zap.Any("struct", v))
			return false
		}
		if !val.validateFields(*v.Fields) {
			log.Error("failed to validate fields")
			return false
		}

		if checkAnswers {
			var a = q.Answers
			if !val.validateAnswers(a) {
				log.Error("coulnd't validate answers")
				return false
			}
			if a.CorrectIDs == nil {
				log.Error("no required CorrectIDs slice for multiple choice model")
				return false
			}

			count := 0
			for _, v := range *v.Fields {
				if slice.Contains(*a.CorrectIDs, v.FieldID) {
					count++
				}
			}
			if count != len(*a.CorrectIDs) {
				log.Error("CorrectIDs contains id that is not in fields slice")
				return false
			}

			if count > *v.MaxChoices {
				log.Error("Correct answers more than max choices")
				return false
			}

			if count == 0 {
				log.Error("no correct answers")
				return false
			}

			return true
		}
	case domain.QuestionTypeManualInput:
		if checkAnswers {
			var a = q.Answers
			if !val.validateAnswers(a) {
				log.Error("coulnd't validate answers")
				return false
			}
			if a.CorrectText == nil {
				log.Error("no required CorrectText for manualInput model")
				return false
			}
			return true
		}
	default:
		log.Error("unknown question type")
		return false
	}

	log.Info("question was validated successfully")
	return true
}

func (val *Validation) validateFields(fs []*domain.CommonField) bool {
	const op = "testsservice.validation.validateFields"
	log := val.log.With(zap.String("op", op))

	met := make(map[int]struct{})
	for _, v := range fs {
		if v == nil {
			log.Error("field is null")
			return false
		}
		if v.FieldID <= 0 {
			log.Error("FieldID is 0 or less")
			return false
		}
		if _, ok := met[v.FieldID]; ok {
			log.Error("repeated FieldID in fields array")
			return false
		}
		met[v.FieldID] = struct{}{}
	}

	return true
}

func (val *Validation) validateAnswers(a *domain.AnswerModel) bool {
	const op = "testsservice.validation.validateAnswers"
	log := val.log.With(zap.String("op", op))
	if a == nil {
		log.Error("answers struct is nil")
		return false
	}

	if a.CorrectID != nil && *a.CorrectID <= 0 {
		log.Error("CorrectID 0 or less")
		return false
	}

	if a.CorrectIDs != nil && len(*a.CorrectIDs) == 0 {
		log.Error("no correct answers in ids slice, len is zero")
		return false
	}

	if a.CorrectIDs != nil {
		met := make(map[int]struct{})
		for _, v := range *a.CorrectIDs {
			if _, ok := met[v]; ok {
				log.Error("repeated ids in CorrectIDs slice")
				return false
			}
			met[v] = struct{}{}
		}
	}

	if a.CorrectText != nil && *a.CorrectText == "" {
		log.Error("CorrectText presented but empty string")
		return false
	}

	return true
}
