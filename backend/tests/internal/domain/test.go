package domain

type Test struct {
	ID        string      `json:"id" bson:"_id"`
	Title     string      `json:"title" bson:"title"`
	Type      string      `json:"type" bson:"type"`
	UserID    int32       `json:"user_id" bson:"user_id"`
	ShortText string      `json:"short_text" bson:"short_text"`
	LongText  string      `json:"long_text" bson:"long_text"`
	MainImage *Image      `json:"main_image" bson:"main_image"`
	Questions []*Question `json:"questions" bson:"questions"`
	Tags      []string    `json:"tags" bson:"tags"`
}
