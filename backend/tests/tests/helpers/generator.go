package helpers

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/numbers"
	p "github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/pointer"
	"time"
)

const (
	TagsMax      = 10
	TagsMin      = 0
	QuestionsMax = 20
	QuestionsMin = 1
	PointsMin    = 10
	PointsMax    = 50
	VariantsMin  = 2
	VariantsMax  = 10
)

func GenerateRandomTest(t string) Test {

	faker := gofakeit.New(uint64(time.Now().Unix()))

	title := faker.ProductName() + " Test"
	longDescr := "Test description " + faker.Quote()
	shortDescr := "Short description " + faker.ProductDescription()
	tagsCount := numbers.RandomInt(TagsMin, TagsMax)
	tags := make([]string, 0, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags = append(tags, faker.ProductCategory())
	}

	qtypes := []string{QuestionTypeSingleChoice, QuestionTypeMultipleChoice, QuestionTypeManualInput}
	questionsCount := numbers.RandomInt(QuestionsMin, QuestionsMax)
	questions := make([]*Question, 0, questionsCount)
	for i := 0; i < questionsCount; i++ {
		var q Question
		var a *Answer
		if t == TestTypeTest || t == TestTypeStrictTest {
			a = &Answer{}
		}
		id := i + 1
		qtype := qtypes[numbers.RandomInt(0, len(qtypes)-1)]
		qLongDescr := "Question description " + faker.Quote()
		qShortDescr := faker.Question()
		required := faker.Bool()
		var points *int
		if t == TestTypeStrictTest {
			points = p.Int(numbers.RandomInt(PointsMin, PointsMax))
		}

		var variants Variants
		switch qtype {
		case QuestionTypeSingleChoice:
			fieldsCount := numbers.RandomInt(VariantsMin, VariantsMax)
			fields := make([]*VariantField, 0, fieldsCount)
			for j := 0; j < fieldsCount; j++ {
				text := faker.ProductDescription()
				field := VariantField{
					ID:    j + 1,
					Text:  &text,
					Image: nil,
				}
				fields = append(fields, &field)
			}
			if a != nil {
				a.CorrectID = p.Int(numbers.RandomInt(1, fieldsCount))
			}
			variants.VariantSingleChoice = &VariantSingleChoice{
				SingleChoiceFields: &fields,
			}
		case QuestionTypeMultipleChoice:
			fieldsCount := numbers.RandomInt(VariantsMin, VariantsMax)
			fields := make([]*VariantField, 0, fieldsCount)
			correctIDs := make([]int, 0)
			for j := 0; j < fieldsCount; j++ {
				text := faker.ProductDescription()
				field := VariantField{
					ID:    j + 1,
					Text:  &text,
					Image: nil,
				}
				fields = append(fields, &field)
				if faker.Bool() {
					correctIDs = append(correctIDs, j+1)
				}
			}
			if len(correctIDs) == 0 {
				correctIDs = append(correctIDs, 1)
			}
			if a != nil {
				a.CorrectIDs = &correctIDs
			}
			variants.VariantMultipleChoice = &VariantMultipleChoice{
				MaxChoices:           &fieldsCount,
				MultipleChoiceFields: &fields,
			}
		case QuestionTypeManualInput:
			textAnswer := faker.Word()
			if a != nil {
				a.CorrectText = &textAnswer
			}
		default:
			panic("invalid question type")
		}

		q = Question{
			ID:        id,
			Type:      &qtype,
			LongText:  &qLongDescr,
			ShortText: &qShortDescr,
			Required:  required,
			Answer:    a,
			Points:    points,
			Variants:  &variants,
		}
		questions = append(questions, &q)
	}

	test := Test{
		Title:     &title,
		LongText:  &longDescr,
		ShortText: &shortDescr,
		Tags:      &tags,
		Questions: &questions,
		MainImage: nil,
		Type:      &t,
		CreatorID: nil,
	}

	return test
}
