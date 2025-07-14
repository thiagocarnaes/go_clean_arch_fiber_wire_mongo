package entities

type Group struct {
	ID      string   `bson:"_id"`
	Name    string   `bson:"name"`
	Members []string `bson:"members"`
}
