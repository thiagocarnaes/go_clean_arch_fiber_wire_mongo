package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/interfaces/repositories"
	"user-management/internal/infrastructure/database"
)

type GroupRepository struct {
	collection *mongo.Collection
}

func NewGroupRepository(db *database.MongoDB) repositories.IGroupRepository {
	return &GroupRepository{collection: db.DB.Collection("groups")}
}

func (r *GroupRepository) Create(ctx context.Context, group *entities.Group) error {
	_, err := r.collection.InsertOne(ctx, group)
	return err
}

func (r *GroupRepository) GetByID(ctx context.Context, id string) (*entities.Group, error) {
	var group entities.Group
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) List(ctx context.Context) ([]*entities.Group, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var groups []*entities.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *GroupRepository) Update(ctx context.Context, group *entities.Group) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": group.ID}, bson.M{"$set": group})
	return err
}

func (r *GroupRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *GroupRepository) AddUserToGroup(ctx context.Context, groupID, userID string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": groupID}, bson.M{"$addToSet": bson.M{"members": userID}})
	return err
}

func (r *GroupRepository) RemoveUserFromGroup(ctx context.Context, groupID, userID string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": groupID}, bson.M{"$pull": bson.M{"members": userID}})
	return err
}
