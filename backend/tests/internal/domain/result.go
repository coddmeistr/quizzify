package domain

type Result struct {
	TestID      string            `json:"test_id" bson:"test_id"`           // Test id
	UserID      int               `json:"user_id" bson:"user_id"`           // User id
	UserAnswers []UserAnswerModel `json:"user_answers" bson:"user_answers"` // Storing all user choses
	ResultID    *int              `json:"result_id" bson:"result_id"`       // For test. Test contains result with this id
	Percentage  *int              `json:"percentage" bson:"percentage"`     // For strict-test. Percents of right answers
}
