package testshandlers

import (
	"github.com/google/uuid"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/domain"
)

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
	Type      *string   `json:"type" validate:"required"`
	LongText  *string   `json:"long_text"`
	ShortText *string   `json:"short_text" validate:"required"`
	Required  *bool     `json:"required"`
	Variants  *Variants `json:"variants" validate:"required,dive"`
}

type Image struct {
	Name    *string `json:"name"`
	Content *[]byte `json:"content"`
}

type VariantField struct {
	Text         *string       `json:"text" validate:"required"`
	Image        *Image        `json:"image"`
	SimpleAnswer *SimpleAnswer `json:"simple_answer"`
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

type SimpleAnswer struct {
	IsCorrect bool `json:"is_correct"`
}

func (t *Test) ToDomain() *domain.Test {

	id := uuid.New().String()

	domainQuestions := make([]*domain.Question, 0, len(*t.Questions))
	for _, v := range *t.Questions {
		domainQuestions = append(domainQuestions, v.ToDomain())
	}

	var domainImage *domain.Image
	if t.MainImage != nil {
		domainImage = t.MainImage.ToDomain()
	}

	return &domain.Test{
		ID:        &id,
		Title:     t.Title,
		UserID:    t.CreatorID,
		Type:      t.Type,
		ShortText: t.ShortText,
		LongText:  t.LongText,
		MainImage: domainImage,
		Tags:      t.Tags,
		Questions: &domainQuestions,
	}
}

func (i *Image) ToDomain() *domain.Image {
	return &domain.Image{
		Name:    i.Name,
		Content: i.Content,
	}
}

func (q *Question) ToDomain() *domain.Question {

	id := uuid.New().String()

	var domainVariants *domain.Variants
	if q.Variants != nil {
		domainVariants = q.Variants.ToDomain()
	}

	return &domain.Question{
		ID:        &id,
		Type:      q.Type,
		LongText:  q.LongText,
		ShortText: q.ShortText,
		Required:  q.Required,
		Variants:  domainVariants,
	}
}

func (a *Variants) ToDomain() *domain.Variants {
	var domainSingleChoice *domain.VariantSingleChoice
	var domainMultipleChoice *domain.VariantMultipleChoice

	if a.VariantSingleChoice != nil {
		domainSingleChoice = a.VariantSingleChoice.ToDomain()
	}
	if a.VariantMultipleChoice != nil {
		domainMultipleChoice = a.VariantMultipleChoice.ToDomain()
	}

	return &domain.Variants{
		VariantSingleChoice:   domainSingleChoice,
		VariantMultipleChoice: domainMultipleChoice,
	}
}

func (a VariantMultipleChoice) ToDomain() *domain.VariantMultipleChoice {
	domainFields := make([]*domain.VariantField, 0, len(*a.MultipleChoiceFields))
	for _, v := range *a.MultipleChoiceFields {
		domainFields = append(domainFields, v.ToDomain())
	}

	return &domain.VariantMultipleChoice{
		MaxChoices:           a.MaxChoices,
		MultipleChoiceFields: &domainFields,
	}
}

func (a VariantSingleChoice) ToDomain() *domain.VariantSingleChoice {

	domainFields := make([]*domain.VariantField, 0, len(*a.SingleChoiceFields))
	for _, v := range *a.SingleChoiceFields {
		domainFields = append(domainFields, v.ToDomain())
	}

	return &domain.VariantSingleChoice{
		SingleChoiceFields: &domainFields,
	}
}

func (a *VariantField) ToDomain() *domain.VariantField {

	var domainImage *domain.Image
	if a.Image != nil {
		domainImage = a.Image.ToDomain()
	}

	var domainAnswer *domain.AnswerSimple
	if a.SimpleAnswer != nil {
		domainAnswer = &domain.AnswerSimple{IsCorrect: a.SimpleAnswer.IsCorrect}
	}

	return &domain.VariantField{
		Text:         a.Text,
		Image:        domainImage,
		AnswerSimple: domainAnswer,
	}
}

func (t *UpdateTestPreviewRequest) ToDomain() *domain.Test {

	var domainImage *domain.Image
	if t.MainImage != nil {
		domainImage = t.MainImage.ToDomain()
	}

	return &domain.Test{
		Title:     t.Title,
		ShortText: t.ShortText,
		LongText:  t.LongText,
		MainImage: domainImage,
		Tags:      t.Tags,
	}
}
