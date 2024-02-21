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

func (v *Validation) validateForm(test domain.Test) bool {
	for _, q := range test.Questions {
		if !v.validateQuestion(*q, false) {
			return false
		}
	}

	return true
}

func (v *Validation) validateQuiz(test domain.Test) bool {
	for _, q := range test.Questions {
		if !v.validateQuestion(*q, false) {
			return false
		}
	}

	return true
}

func (v *Validation) validateTest(test domain.Test) bool {
	for _, q := range test.Questions {
		if !v.validateQuestion(*q, true) {
			return false
		}
	}

	return true
}

func (v *Validation) validateStrictTest(test domain.Test) bool {
	for _, q := range test.Questions {
		if !v.validateQuestion(*q, true) {
			return false
		}
	}

	return true
}

func (v *Validation) validateQuestion(q domain.Question, checkAnswers bool) bool {

	switch q.Type {
	case QuestionTypeSingleChoice:
		if !(q.Variants.VariantSingleChoice != nil &&
			len(q.Variants.VariantSingleChoice.SingleChoiceFields) > 0) {
			v.log.Error("no single choice structure or no variants")
			return false
		}

		if checkAnswers {
			count := 0
			for _, field := range q.Variants.VariantSingleChoice.SingleChoiceFields {
				if field.AnswerSimple.IsCorrect {
					count++
				}
			}
			if count != 1 {
				v.log.Error("single choice variants has more than 1 correct answer")
				return false
			}
		}
	case QuestionTypeMultipleChoice:
		if !(q.Variants.VariantMultipleChoice != nil &&
			len(q.Variants.VariantMultipleChoice.MultipleChoiceFields) > 0 &&
			q.Variants.VariantMultipleChoice.MaxChoices > 0) {
			v.log.Error("no multiple choice structure or no variants", zap.Any("struct", q.Variants.VariantMultipleChoice))
			return false
		}

		if checkAnswers {
			count := 0
			for _, field := range q.Variants.VariantMultipleChoice.MultipleChoiceFields {
				if field.AnswerSimple.IsCorrect {
					count++
				}
			}
			if count > q.Variants.VariantMultipleChoice.MaxChoices || count <= 0 {
				v.log.Error("multiple choice variants has more than max choices or less than 1 correct answer")
				return false
			}
		}
	default:
		v.log.Error("unknown question type")
		return false
	}

	v.log.Info("question was validated successfully")
	return true
}
