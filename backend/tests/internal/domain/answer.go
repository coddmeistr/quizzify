package domain

type AnswerSimple struct {
	IsCorrect bool `json:"is_correct" bson:"is_correct"`
}
