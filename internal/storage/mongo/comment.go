package mongoStore

import (
	"context"
	"errors"
	"time"

	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoComment = errors.New("no find comment")
)

type commentRepository struct {
	c *mongo.Collection
}

func newCommentRepo(c *mongo.Collection) *commentRepository {
	return &commentRepository{
		c: c,
	}
}

func (r *commentRepository) createIndexes(ctx context.Context) ([]string, error) {
	indexed := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "creator_id", Value: 1}},
			Options: options.Index().SetName("comment_creator_id_idx"),
		},
		{
			Keys:    bson.D{{Key: "worker_id", Value: 1}},
			Options: options.Index().SetName("comment_worker_id_idx"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("comment_created_at_idx"),
		},
	}
	res, err := r.c.Indexes().CreateMany(ctx, indexed)

	return res, err
}
func (r *commentRepository) Create(ctx context.Context, comment *commentmodel.Comment) (string, error) {
	comment.CreatedAt = time.Now().UTC()
	res, err := r.c.InsertOne(ctx, comment)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (r *commentRepository) FindMany(ctx context.Context, filter bson.D, limit, skip int64) ([]commentmodel.Comment, error) {
	cur, err := r.c.Find(
		ctx,
		filter,
		options.Find().SetSort(bson.M{"created_at": -1}),
		options.Find().SetLimit(limit),
		options.Find().SetSkip(skip),
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoComment
	case err != nil:
		return nil, err
	}
	comments := []commentmodel.Comment{}
	err = cur.All(ctx, &comments)
	if len(comments) == 0 {
		return nil, ErrNoComment
	}
	return comments, err
}

func (r *commentRepository) CountWorkerLikesDislikes(ctx context.Context, workerID string) (*commentmodel.CountLikeDislikeWorker, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "worker_id", Value: workerID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$is_like"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
	cur, err := r.c.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	result := commentmodel.CountLikeDislikeWorker{}
	for cur.Next(ctx) {
		var tempResult struct {
			ID    bool `bson:"_id"`
			Count int  `bson:"count"`
		}
		if err := cur.Decode(&tempResult); err != nil {
			return nil, err
		}
		switch tempResult.ID {
		case true:
			result.Likes = tempResult.Count
		case false:
			result.Dislike = tempResult.Count
		}
	}

	return &result, nil
}

func (r *commentRepository) Update(ctx context.Context, filter bson.D, upd *commentmodel.Comment) error {
	commentID := upd.ID
	if upd.ID != "" {
		upd.ID = ""
	}
	update := bson.D{
		{Key: "$set", Value: upd},
	}
	_, err := r.c.UpdateOne(ctx, filter, update)
	switch {
	case err == mongo.ErrNoDocuments:
		return ErrNoComment
	case err != nil:
		return err
	}
	upd.ID = commentID

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, filter bson.D) error {
	_, err := r.c.DeleteOne(ctx, filter)
	switch {
	case err == mongo.ErrNoDocuments:
		return ErrNoComment
	case err != nil:
		return err
	}

	return nil
}
