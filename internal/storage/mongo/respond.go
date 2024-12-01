package mongoStore

import (
	"context"
	"errors"
	"time"

	respondmodel "github.com/Negat1v9/work-marketplace/model/respond"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoRespond = errors.New("no find respond")
)

type respondRepository struct {
	c *mongo.Collection
}

func newRespondRepo(c *mongo.Collection) *respondRepository {
	return &respondRepository{
		c: c,
	}
}
func (r *respondRepository) Create(ctx context.Context, respond *respondmodel.Respond) (string, error) {
	respond.CreatedAt = time.Now().UTC()
	res, err := r.c.InsertOne(ctx, respond)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}
func (r *respondRepository) Find(ctx context.Context, filter bson.D) (*respondmodel.Respond, error) {
	var respond respondmodel.Respond
	err := r.c.FindOne(ctx, filter).Decode(&respond)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoRespond
	case err != nil:
		return nil, err
	}

	return &respond, nil
}
func (r *respondRepository) FindMany(ctx context.Context, filter bson.D, limit, skip int64) ([]respondmodel.Respond, error) {
	cur, err := r.c.Find(
		ctx,
		filter,
		options.Find().SetLimit(limit),
		options.Find().SetSkip(skip),
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoRespond
	case err != nil:
		return nil, err
	}
	responds := []respondmodel.Respond{}
	err = cur.All(ctx, &responds)
	if len(responds) == 0 {
		return nil, ErrNoRespond
	}

	return responds, err
}
func (r *respondRepository) Delete(ctx context.Context, filter bson.D) error {
	_, err := r.c.DeleteOne(ctx, filter)
	switch {
	case err == mongo.ErrNoDocuments:
		return ErrNoRespond
	case err != nil:
		return err
	}

	return nil
}
