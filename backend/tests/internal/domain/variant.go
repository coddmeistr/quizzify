package domain

type VariantSingleChoice struct {
	SingleChoiceFields *[]*VariantField `json:"fields" bson:"fields"`
}

type VariantMultipleChoice struct {
	MaxChoices           *int             `json:"max" bson:"max"`
	MultipleChoiceFields *[]*VariantField `json:"fields" bson:"fields"`
}

type VariantField struct {
	Text         *string       `json:"text" bson:"text"`
	Image        *Image        `json:"image,omitempty" bson:"image"`
	AnswerSimple *AnswerSimple `json:"answer_simple,omitempty" bson:"answer_simple,omitempty"`
}

type Variants struct {
	VariantSingleChoice   *VariantSingleChoice   `json:"single_choice,omitempty" bson:"single_choice"`
	VariantMultipleChoice *VariantMultipleChoice `json:"multiple_choice,omitempty" bson:"multiple_choice"`
}
