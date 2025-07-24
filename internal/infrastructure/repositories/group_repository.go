package repositories

import (
	"context"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/interfaces/repositories"
	"user-management/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GroupRepository struct {
	*BaseRepository
	collection *mongo.Collection
}

func NewGroupRepository(dbManager *database.DatabaseManager) repositories.IGroupRepository {
	baseRepo := NewBaseRepository(dbManager)

	// Get MongoDB collection
	collection, err := baseRepo.GetMongoCollection("groups")
	if err != nil {
		// For now, we'll panic. In production, you might want to handle this differently
		panic("Failed to initialize group repository: " + err.Error())
	}

	return &GroupRepository{
		BaseRepository: baseRepo,
		collection:     collection,
	}
}

func (r *GroupRepository) Create(ctx context.Context, group *entities.Group) error {
	group.ID = bson.NewObjectID()
	_, err := r.collection.InsertOne(ctx, group)
	return err
}

func (r *GroupRepository) GetByID(ctx context.Context, id string) (*entities.Group, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var group entities.Group
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&group)
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
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": group.ID}, bson.M{"$set": bson.M{
		"name":    group.Name,
		"members": group.Members,
	}})
	return err
}

func (r *GroupRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *GroupRepository) AddUserToGroup(ctx context.Context, groupID, userID string) error {
	groupObjectID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": groupObjectID}, bson.M{"$addToSet": bson.M{"members": userID}})
	return err
}

func (r *GroupRepository) RemoveUserFromGroup(ctx context.Context, groupID, userID string) error {
	groupObjectID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": groupObjectID}, bson.M{"$pull": bson.M{"members": userID}})
	return err
}
