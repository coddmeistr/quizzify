package domain

type UserAnswerModel struct {
	QuestionID int     `json:"question_id" bson:"question_id"`
	ChosenID   *int    `json:"chosen_id" bson:"chosen_id"`
	ChosenIDs  *[]int  `json:"chosen_ids" bson:"chosen_ids"`
	WritedText *string `json:"writed_text" bson:"writed_text"`
}
