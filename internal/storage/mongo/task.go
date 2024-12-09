package mongoStore

import (
	"context"
	"errors"
	"time"

	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoTask = errors.New("no task")
)

type taskRepository struct {
	c *mongo.Collection
}

func newTaskRepo(coll *mongo.Collection) *taskRepository {
	return &taskRepository{
		c: coll,
	}
}

func (r *taskRepository) Create(ctx context.Context, task *taskmodel.Task) (string, error) {
	task.CreatedAt = time.Now().UTC()

	res, err := r.c.InsertOne(ctx, task)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (r *taskRepository) Find(ctx context.Context, filter bson.D) (*taskmodel.Task, error) {
	var task taskmodel.Task
	err := r.c.FindOne(ctx, filter).Decode(&task)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoTask
	case err != nil:
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) FindMany(ctx context.Context, filter bson.D, proj bson.M, limit, skip int64) ([]taskmodel.Task, error) {

	cur, err := r.c.Find(
		ctx,
		filter,
		options.Find().SetProjection(proj),
		options.Find().SetLimit(limit),
		options.Find().SetSkip(skip),
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoTask
	case err != nil:
		return nil, err
	}

	tasks := []taskmodel.Task{}

	err = cur.All(ctx, &tasks)
	if len(tasks) == 0 {
		return nil, ErrNoTask
	}

	return tasks, err
}

func (r *taskRepository) FindProj(ctx context.Context, filter bson.D, proj bson.M) (*taskmodel.Task, error) {
	var task taskmodel.Task
	err := r.c.FindOne(
		ctx,
		filter,
		options.FindOne().SetProjection(proj)).Decode(&task)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoTask
	case err != nil:
		return nil, err
	}

	return &task, nil
}
func (r *taskRepository) Count(ctx context.Context, filter bson.D) (int64, error) {
	number, err := r.c.CountDocuments(ctx, filter)
	switch {
	case err == mongo.ErrNoDocuments:
		return 0, nil
	case err != nil:
		return 0, err
	}

	return number, nil
}
func (r *taskRepository) Update(ctx context.Context, filter bson.D, task *taskmodel.Task) (*taskmodel.Task, error) {
	task.UpdatedAt = time.Now().UTC()
	taskID := task.ID
	if task.ID != "" {
		task.ID = ""
	}
	update := bson.D{
		{Key: "$set", Value: task},
	}

	var afterTask taskmodel.Task
	err := r.c.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&afterTask)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoTask
	case err != nil:
		return nil, err
	}
	task.ID = taskID

	return &afterTask, err

}
func (r *taskRepository) Delete(ctx context.Context, filter bson.D) error {
	_, err := r.c.DeleteOne(ctx, filter)

	switch {
	case err == mongo.ErrNoDocuments:
		return ErrNoTask
	case err != nil:
		return err
	}
	return nil
}
