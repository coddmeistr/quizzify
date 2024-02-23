package domain

type AnswerSimple struct {
	IsCorrect *bool `bson:"is_correct"`
}
