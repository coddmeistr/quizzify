package domain

type Image struct {
	Name    string `json:"name" bson:"name"`
	Content []byte `json:"content" bson:"content"`
}
