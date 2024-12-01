package paymentservice

import (
	"context"

	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *PaymentServiceImpl) successBalanePayTrx(ctx context.Context, user *usermodel.User, payment *paymentmodel.Payment) error {
	session, err := s.store.StartSession()
	if err != nil {
		return err
	}

	err = session.StartTransaction()
	if err != nil {
		return err
	}
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		_, err = s.store.User().Edit(sc, filters.New().Add(filters.UserByID(user.ID)).Filters(), user)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		s.store.Payment().Edit(sc, filters.New().Add(filters.PaymentByID(payment.ID)).Filters(), payment)
		if err != nil {
			session.AbortTransaction(ctx)
			return err
		}

		return session.CommitTransaction(ctx)
	})
	return err
}
