package domain

type Question struct {
	ID        int            `json:"id" bson:"id"`
	Type      *string        `json:"type" bson:"type"`
	LongText  *string        `json:"long_text" bson:"long_text"`
	ShortText *string        `json:"short_text" bson:"short_text"`
	Required  *bool          `json:"required" bson:"required"`
	Variants  *VariantsModel `json:"variants" bson:"variants"`
	Answers   *AnswerModel   `json:"answers,omitempty" bson:"answers"`
}
