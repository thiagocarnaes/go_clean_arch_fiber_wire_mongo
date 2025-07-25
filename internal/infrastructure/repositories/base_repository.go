package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// BaseRepository contém funcionalidades comuns para todos os repositórios
type BaseRepository struct {
	collection *mongo.Collection
}

// NewBaseRepository cria uma nova instância do BaseRepository
func NewBaseRepository(collection *mongo.Collection) *BaseRepository {
	return &BaseRepository{
		collection: collection,
	}
}

// Count retorna o total de documentos na coleção
func (r *BaseRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountWithFilter retorna o total de documentos que atendem ao filtro
func (r *BaseRepository) CountWithFilter(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteByID remove um documento pelo ID
func (r *BaseRepository) DeleteByID(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// ExistsByID verifica se um documento existe pelo ID
func (r *BaseRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindWithPagination encontra documentos com paginação
func (r *BaseRepository) FindWithPagination(ctx context.Context, filter bson.M, offset int64, limit int64) (*mongo.Cursor, error) {
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	return r.collection.Find(ctx, filter, opts)
}
