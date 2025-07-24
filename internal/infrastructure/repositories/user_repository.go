package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/interfaces/repositories"
	"user-management/internal/infrastructure/database"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *database.MongoDB) repositories.IUserRepository {
	return &UserRepository{collection: db.DB.Collection("users")}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	user.ID = bson.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user entities.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*entities.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*entities.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
		"name":  user.Name,
		"email": user.Email,
	}})
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
