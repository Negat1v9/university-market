package mongoStore

import (
	"context"
	"errors"
	"time"

	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoEvent = errors.New("no event")
)

type eventRepository struct {
	c *mongo.Collection
}

func newEventRepo(c *mongo.Collection) *eventRepository {
	return &eventRepository{
		c: c,
	}
}

func (r *eventRepository) Create(ctx context.Context, event *eventmodel.Event) (string, error) {
	event.CreatedAt = time.Now().UTC()
	res, err := r.c.InsertOne(ctx, event)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (r *eventRepository) Update(ctx context.Context, filter bson.D, event *eventmodel.Event) (*eventmodel.Event, error) {
	var afterEvent eventmodel.Event
	update := bson.D{
		{Key: "$set", Value: event},
	}
	err := r.c.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&afterEvent)

	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoEvent
	case err != nil:
		return nil, err
	}

	return &afterEvent, nil
}

func (r *eventRepository) FindOne(ctx context.Context, filter bson.D) (*eventmodel.Event, error) {
	var event eventmodel.Event

	err := r.c.FindOne(ctx, filter).Decode(&event)

	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoEvent
	case err != nil:
		return nil, err
	}

	return &event, nil
}
