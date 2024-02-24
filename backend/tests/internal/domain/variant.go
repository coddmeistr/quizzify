package domain

type SingleChoice struct {
	Fields *[]*CommonField `json:"fields" bson:"fields"`
}

type MultipleChoice struct {
	MaxChoices *int            `json:"max" bson:"max"`
	Fields     *[]*CommonField `json:"fields" bson:"fields"`
}

// CommonField basic answer field, that contains information that shows to user.
// User allowed to choose between CommonFields.
// FieldID contains unique numeric value for each CommonField in VariantModel.
// FieldID must be linked with ids in AnswerModel to have correct answer to the question.
type CommonField struct {
	FieldID int     `json:"id" bson:"id"`
	Text    *string `json:"text" bson:"text"`
	Image   *Image  `json:"image,omitempty" bson:"image"`
}

// VariantsModel contains all structs, that represents specific variants model type for the question.
// For example: question can be with single answer with choices, then it must use SingleChoice struct
// that represents model for such question type with variant fields.
// Each different variants model like single-choices or multiple-choiced should implement it's own struct embedded in this struct.
type VariantsModel struct {
	SingleChoice   *SingleChoice   `json:"single_choice,omitempty" bson:"single_choice"`
	MultipleChoice *MultipleChoice `json:"multiple_choice,omitempty" bson:"multiple_choice"`
}
