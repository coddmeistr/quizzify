package helpers

const (
	TestTypeForm       = "form"
	TestTypeQuiz       = "quiz"
	TestTypeTest       = "test"
	TestTypeStrictTest = "strict_test"

	QuestionTypeSingleChoice   = "single_choice"
	QuestionTypeMultipleChoice = "multiple_choice"
	QuestionTypeManualInput    = "manual_input"
)

type ApplyTestRequest struct {
	UserAnswers []*UserAnswer `json:"user_answers" validate:"required,dive"`
}

type UserAnswer struct {
	QuestionID int     `json:"question_id" validate:"required,gte=1"`
	ChosenID   *int    `json:"chosen_id"`
	ChosenIDs  *[]int  `json:"chosen_ids"`
	WritedText *string `json:"writed_text"`
}

type UpdateTestPreviewRequest struct {
	Title     *string   `json:"title"`
	ShortText *string   `json:"short_text"`
	LongText  *string   `json:"long_text"`
	MainImage *Image    `json:"main_image"`
	Tags      *[]string `json:"tags"`
}

type CreateTestRequest struct {
	Test
}

type Test struct {
	Title     *string      `json:"title" validate:"required"`
	CreatorID *int         `json:"creator_id" validate:"required"`
	Type      *string      `json:"type" validate:"required"`
	ShortText *string      `json:"short_text" validate:"required"`
	LongText  *string      `json:"long_text" validate:"required"`
	MainImage *Image       `json:"main_image" validate:"required"`
	Questions *[]*Question `json:"questions" validate:"required,gte=1,dive"`
	Tags      *[]string    `json:"tags"`
}

type Question struct {
	ID        int       `json:"id"  validate:"required,gte=1"`
	Type      *string   `json:"type" validate:"required"`
	LongText  *string   `json:"long_text"`
	ShortText *string   `json:"short_text" validate:"required"`
	Required  bool      `json:"required"`
	Points    *int      `json:"points"`
	Variants  *Variants `json:"variants" validate:"required,dive"`
	Answer    *Answer   `json:"answers,omitempty"`
}

type Image struct {
	Name    *string `json:"name"`
	Content *[]byte `json:"content"`
}

type VariantField struct {
	ID    int     `json:"id"`
	Text  *string `json:"text" validate:"required"`
	Image *Image  `json:"image"`
}

type Variants struct {
	VariantSingleChoice   *VariantSingleChoice   `json:"single_choice"`
	VariantMultipleChoice *VariantMultipleChoice `json:"multiple_choice"`
}

type VariantSingleChoice struct {
	SingleChoiceFields *[]*VariantField `json:"fields" validate:"required,dive"`
}

type VariantMultipleChoice struct {
	MaxChoices           *int             `json:"max" validate:"gte=1"`
	MultipleChoiceFields *[]*VariantField `json:"fields" validate:"required,dive"`
}

type Answer struct {
	CorrectID   *int       `json:"correct_id"`
	CorrectIDs  *[]int     `json:"correct_ids"`
	CorrectText *string    `json:"correct_text"`
	Params      *[]*Params `json:"params"`
}

type Params struct {
	FieldID int            `json:"id"`
	Params  *[]*FlexParams `json:"params"`
}

type FlexParams struct {
	Name       *string `json:"name"`
	Effect     *int    `json:"effect"`
	IsNegative *bool   `json:"is_negative"`
}

type UserInfo struct {
	ID          int   `json:"id"`
	Permissions []int `json:"permissions"`
}
