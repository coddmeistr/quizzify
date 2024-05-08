package mongo

import "go.mongodb.org/mongo-driver/bson"

// getProjectionForAnswers return projection for all existing bson answer fields in domain entities
// It used to get quick bson to exclude all answer fields from bson document
func getProjectionForAnswers() bson.D {
	return bson.D{{"questions.answers", 0}}
}
