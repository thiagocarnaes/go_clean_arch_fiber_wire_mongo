package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Email    string        `bson:"email"`
	IsActive bool          `bson:"is_active"`
}
