package repositories

import (
	"context"
	"fmt"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/interfaces/repositories"
	"user-management/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	// CollectionNameUsers é o nome da coleção de usuários no MongoDB
	mongoRegex   = "$regex"
	mongoOptions = "$options"
)

type UserRepository struct {
	*BaseRepository
	collection *mongo.Collection
}

func NewUserRepository(db *database.MongoDB) (repositories.IUserRepository, error) {
	// Check if MongoDB database is valid
	if db == nil || db.DB == nil {
		return nil, fmt.Errorf("failed to get MongoDB collection for users: database connection is nil")
	}

	// Get MongoDB collection
	collection := db.DB.Collection("users")
	if collection == nil {
		return nil, fmt.Errorf("failed to get MongoDB collection for users")
	}

	return &UserRepository{
		BaseRepository: NewBaseRepository(collection),
		collection:     collection,
	}, nil
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

func (r *UserRepository) List(ctx context.Context, offset int64, limit int64) ([]*entities.User, error) {
	cursor, err := r.FindWithPagination(ctx, bson.M{}, offset, limit)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // IMPORTANTE: fecha o cursor para liberar recursos

	var users []*entities.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Search(ctx context.Context, searchTerm string, offset int64, limit int64) ([]*entities.User, error) {
	// Cria filtro de busca usando regex para buscar no nome e email
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{mongoRegex: searchTerm, mongoOptions: "i"}},
			{"email": bson.M{mongoRegex: searchTerm, mongoOptions: "i"}},
		},
	}

	cursor, err := r.FindWithPagination(ctx, filter, offset, limit)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*entities.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) CountSearch(ctx context.Context, searchTerm string) (int64, error) {
	// Cria filtro de busca usando regex para buscar no nome e email
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{mongoRegex: searchTerm, mongoOptions: "i"}},
			{"email": bson.M{mongoRegex: searchTerm, mongoOptions: "i"}},
		},
	}

	return r.CountWithFilter(ctx, filter)
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
		"name":  user.Name,
		"email": user.Email,
	}})
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.DeleteByID(ctx, id)
}
