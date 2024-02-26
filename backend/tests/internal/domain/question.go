package domain

import (
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/slice"
)

type Question struct {
	ID        int            `json:"id" bson:"id"`
	Type      *string        `json:"type" bson:"type"`
	LongText  *string        `json:"long_text" bson:"long_text"`
	ShortText *string        `json:"short_text" bson:"short_text"`
	Required  bool           `json:"required" bson:"required"`
	Variants  *VariantsModel `json:"variants" bson:"variants"`
	Answers   *AnswerModel   `json:"answers,omitempty" bson:"answers"`

	// Strict Test
	Points *int `json:"points" bson:"points"`
}

func (q *Question) ComparePreciseResults(ua UserAnswerModel) int {
	qa := *q.Answers
	switch *q.Type {
	case QuestionTypeSingleChoice:
		if *ua.ChosenID == *qa.CorrectID {
			return 100
		}
		return 0
	case QuestionTypeMultipleChoice:
		count := 0
		for _, v := range *qa.CorrectIDs {
			if slice.Contains(*ua.ChosenIDs, v) {
				count++
			}
		}
		return int((float64(count) / float64(len(*qa.CorrectIDs))) * 100)
	case QuestionTypeManualInput:
		if *ua.WritedText == *qa.CorrectText {
			return 100
		}
		return 0
	default:
		return 0
	}
}
