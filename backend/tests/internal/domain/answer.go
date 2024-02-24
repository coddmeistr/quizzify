package domain

// AnswerModel represents struct which contains all possible
// answer variations
// If question with specific type need only one correct answer it gets CorrectID
// If question has many correct answer it gets CorrectIDs from slice
// If question requires user input, then correct answer should be stored in CorrectText field e.t.
type AnswerModel struct {
	CorrectID   *int                `json:"correct_id" bson:"correct_id"`
	CorrectIDs  *[]int              `json:"correct_ids" bson:"correct_ids"`
	CorrectText *string             `json:"correct_text" bson:"correct_text"`
	FlexParams  *[]*FlexParamsModel `json:"params,omitempty" bson:"params"`
	// Add other possible answer fields
}

// FlexParamsModel represents flex parameters for test type 'test'
// Each field in variants must have it's params to calculate answers based on user's choices
// Each field numbericly increasing or decreasing one or more parameters and
// At the end we get final numeric representation for each parameter and based on that we calculating final result
type FlexParamsModel struct {
	FieldID int           `json:"id" bson:"id"`
	Params  *[]*FlexParam `json:"params" bson:"params"`
}

// FlexParam represents flex parameters for test type 'test'
type FlexParam struct {
	ParamName  *string `json:"name" bson:"name"`               // Name of parameter
	Effect     *int    `json:"effect" bson:"effect"`           // Numeric effect of on this parameter
	IsNegative bool    `json:"is_negative" bson:"is_negative"` // If true, then effect is negative, otherwise positive
}
