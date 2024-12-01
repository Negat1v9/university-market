package mongoStore

import (
	"context"
	"errors"

	tgbotmodel "github.com/Negat1v9/work-marketplace/model/tgBot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoTgCmd    = errors.New("no task")
	ErrTgCmdExist = errors.New("tgCmd is exist")
)

type tgCmdRepository struct {
	c *mongo.Collection
}

func newTgCmdRepo(coll *mongo.Collection) *tgCmdRepository {
	return &tgCmdRepository{
		c: coll,
	}
}

func (t *tgCmdRepository) Create(ctx context.Context, cmd *tgbotmodel.UserCommand) error {
	_, err := t.c.InsertOne(ctx, cmd)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrTgCmdExist
		}
		return err
	}
	return nil
}
func (t *tgCmdRepository) CreateDelete(ctx context.Context, cmd *tgbotmodel.UserCommand) error {
	_, err := t.c.InsertOne(ctx, cmd)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			err = t.Delete(ctx, cmd.ID)
			if err != nil {
				return err
			}
			err = t.Create(ctx, cmd)
			return err
		}
		return err
	}
	return nil
}

func (t *tgCmdRepository) Find(ctx context.Context, userID int64) (*tgbotmodel.UserCommand, error) {
	var cmd tgbotmodel.UserCommand
	filter := bson.M{"_id": userID}

	err := t.c.FindOne(ctx, filter).Decode(&cmd)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoTgCmd
	case err != nil:
		return nil, err
	}
	return &cmd, nil
}

func (t *tgCmdRepository) Delete(ctx context.Context, userID int64) error {
	filter := bson.M{"_id": userID}
	_, err := t.c.DeleteOne(ctx, filter)
	switch {
	case err == mongo.ErrNoDocuments:
		return ErrNoTgCmd
	case err != nil:
		return err
	}
	return nil
}
