package mongoStore

import (
	"context"
	"errors"
	"time"

	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoPayment = errors.New("no find payment by payer id")
)

type paymentRepository struct {
	c *mongo.Collection
}

func newPaymentRepo(c *mongo.Collection) *paymentRepository {
	return &paymentRepository{
		c: c,
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *paymentmodel.Payment) (string, error) {
	payment.CreatedAt = time.Now().UTC()
	res, err := r.c.InsertOne(ctx, payment)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), err
}

func (r *paymentRepository) Find(ctx context.Context, filter bson.D) (*paymentmodel.Payment, error) {

	var payment paymentmodel.Payment

	err := r.c.FindOne(ctx, filter).Decode(&payment)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNoPayment
	case err != nil:
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) Edit(ctx context.Context, filter bson.D, payment *paymentmodel.Payment) error {
	paymentID := payment.ID
	if payment.ID != "" {
		payment.ID = ""
	}
	update := bson.D{
		{Key: "$set", Value: payment},
	}
	_, err := r.c.UpdateOne(ctx, filter, update)
	switch {
	case err == mongo.ErrNoDocuments:
		return ErrNoPayment
	case err != nil:
		return err
	}

	payment.ID = paymentID

	return nil
}

func (r *paymentRepository) Delete(ctx context.Context, filter bson.D) error {
	_, err := r.c.DeleteOne(ctx, filter)
	switch {
	case err == mongo.ErrNoDocuments:
		return ErrNoPayment
	case err != nil:
		return err
	}

	return nil
}
