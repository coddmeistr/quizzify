package domain

// AnswerModel represents struct which contains all possible
// answer variations
// If question with specific type need only one correct answer it gets CorrectID
// If question has many correct answer it gets CorrectIDs from slice
// If question requires user input, then correct answer should be stored in CorrectText field e.t.
type AnswerModel struct {
	CorrectID   *int
	CorrectIDs  *[]int
	CorrectText *string
	// Add other possible answer fields
}
