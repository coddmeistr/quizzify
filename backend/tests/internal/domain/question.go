package domain

type Question struct {
	ID        string    `json:"id" bson:"_id"`
	Type      string    `json:"type" bson:"type"`
	LongText  string    `json:"long_text" bson:"long_text"`
	ShortText string    `json:"short_text" bson:"short_text"`
	Required  bool      `json:"required" bson:"required"`
	Variants  *Variants `json:"variants" bson:"variants"`
}
