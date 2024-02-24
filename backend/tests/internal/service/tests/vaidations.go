package testsservice

import (
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
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
	for _, q := range *test.Questions {
		if !val.validateQuestion(*q, true) {
			return false
		}
	}

	return true
}

func (val *Validation) validateQuestion(q domain.Question, checkAnswers bool) bool {
	const op = "testsservice.validation.validateQuestion"
	log := val.log.With(zap.String("op", op))

	switch *q.Type {
	case QuestionTypeSingleChoice:
		var (
			v = q.Variants.VariantSingleChoice
		)
		if v == nil {
			log.Error("no single choice structure")
			return false
		}
		if v.SingleChoiceFields == nil {
			log.Error("no single choice variants")
			return false
		}

		if len(*v.SingleChoiceFields) <= 0 {
			log.Error("no variants")
			return false
		}

		if checkAnswers {
			count := 0
			for _, field := range *v.SingleChoiceFields {
				if field.AnswerSimple == nil {
					log.Error("answer is nil")
					return false
				}
				if field.AnswerSimple.IsCorrect {
					count++
				}
			}
			if count != 1 {
				log.Error("single choice variants has more than 1 correct answer or 0 correct answers")
				return false
			}
		}
	case QuestionTypeMultipleChoice:
		var (
			v = q.Variants.VariantMultipleChoice
		)
		if v == nil {
			log.Error("no multiple choice structure")
			return false
		}
		if v.MultipleChoiceFields == nil {
			log.Error("no multiple choice variants")
			return false
		}
		if v.MaxChoices == nil {
			log.Error("no max choices")
			return false
		}

		if !(len(*v.MultipleChoiceFields) > 0 &&
			*v.MaxChoices > 0) {
			log.Error("no fields or zero max choices", zap.Any("struct", q.Variants.VariantMultipleChoice))
			return false
		}

		if checkAnswers {
			count := 0
			for _, field := range *v.MultipleChoiceFields {
				if field.AnswerSimple == nil {
					log.Error("answer is nil")
					return false
				}
				if field.AnswerSimple.IsCorrect {
					count++
				}
			}
			if count > *v.MaxChoices || count <= 0 {
				log.Error("multiple choice variants has more than max choices or less than 1 correct answer")
				return false
			}
		}
	default:
		log.Error("unknown question type")
		return false
	}

	log.Info("question was validated successfully")
	return true
}
