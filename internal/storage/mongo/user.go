package mongoStore

import (
	"context"
	"errors"

	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoUser    = errors.New("no user")
	ErrExistUser = errors.New("user is exist")
)

type userRepository struct {
	c *mongo.Collection
}

func newUserRepo(coll *mongo.Collection) *userRepository {
	return &userRepository{
		c: coll,
	}
}

func (r *userRepository) createIndexes(ctx context.Context) ([]string, error) {
	indexed := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "telegram_id", Value: 1}},
			Options: options.Index().SetName("user_telegram_id_idx"),
		},
	}
	res, err := r.c.Indexes().CreateMany(ctx, indexed)

	return res, err
}

func (r *userRepository) Create(ctx context.Context, user *usermodel.User) (string, error) {
	res, err := r.c.InsertOne(ctx, user)
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (r *userRepository) Find(ctx context.Context, filter bson.D) (*usermodel.User, error) {
	var user usermodel.User
	err := r.c.FindOne(ctx, filter).Decode(&user)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoUser
	case err != nil:
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindProj(ctx context.Context, filter bson.D, proj bson.M) (*usermodel.User, error) {
	var user usermodel.User
	err := r.c.FindOne(ctx, filter, options.FindOne().SetProjection(proj)).Decode(&user)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoUser
	case err != nil:
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) Edit(ctx context.Context, filter bson.D, user *usermodel.User) (*usermodel.User, error) {
	userID := user.ID
	if user.ID != "" {
		user.ID = ""
	}
	var afterUser usermodel.User
	err := r.c.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$set": user},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&afterUser)
	// _, err := r.c.UpdateOne(ctx, filter, bson.M{"$set": user}, options.Update().set)

	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoUser
	case err != nil:
		return nil, err
	}
	user.ID = userID

	return &afterUser, nil
}
