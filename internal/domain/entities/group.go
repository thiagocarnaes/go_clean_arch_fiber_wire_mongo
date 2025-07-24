package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type Group struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	Name    string        `bson:"name"`
	Members []string      `bson:"members"`
}
