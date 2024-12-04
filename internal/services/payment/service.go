package paymentservice

import (
	"context"
	"log/slog"
	"time"

	"github.com/Negat1v9/work-marketplace/internal/storage"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	filters "github.com/Negat1v9/work-marketplace/internal/storage/mongo/filter"
	tgpayment "github.com/Negat1v9/work-marketplace/internal/tgBot/payment"
	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

var (
	// the minimum time that can pass after creating the first payment, which must pass
	minPaymentTime time.Duration = time.Minute * 1
	// the amount of telegram stars that users are entitled to if they participate in the referral program
	bonusReferral int = 200
)

type PaymentServiceImpl struct {
	log        *slog.Logger
	store      storage.Store
	tgBotToken string
}

func NewServicePayment(log *slog.Logger, tgBotToken string, store storage.Store) PaymentService {
	return &PaymentServiceImpl{
		log:        log,
		store:      store,
		tgBotToken: tgBotToken,
	}
}

func (s *PaymentServiceImpl) CreateInvoiceLink(ctx context.Context, workerID string, data *paymentmodel.PaymentLinkReq) (string, error) {
	filter := filters.New().
		Add(filters.PaymentByUser(workerID)).
		Add(filters.PaymentByStatus(paymentmodel.Pending))
	payment, err := s.store.Payment().Find(ctx, filter.Filters())

	switch {
	// payment exist
	case payment != nil:
		// checking that the previous payment was created no earlier than minPaymentTime ago
		if ok := checkTimeCreatedPayment(minPaymentTime, payment.CreatedAt); !ok {
			return "", httpresponse.NewError(425, "")
		}
		// deleting a previous failed payment
		err = s.store.Payment().Delete(ctx, filters.New().Add(filters.PaymentByID(payment.ID)).Filters())
		if err != nil {
			s.log.Error("create invoice link", slog.String("err", err.Error()))
			httpresponse.ServerError()
		}
	// return only if store error
	case err != nil && err != mongoStore.ErrNoPayment:
		s.log.Error("create invoice link", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	payment = paymentmodel.NewPayment(data.Amount, workerID)
	paymentID, err := s.store.Payment().Create(ctx, payment)
	if err != nil {
		s.log.Error("create invoice link", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	link, err := tgpayment.CreateInvoiceLink(ctx, paymentID, s.tgBotToken, data.Amount)
	if err != nil {
		s.log.Error("create invoice link", slog.String("err", err.Error()))
		return "", httpresponse.ServerError()
	}

	return link, nil
}

func (s *PaymentServiceImpl) SuccessBalancePayment(ctx context.Context, sp *paymentmodel.SuccessPayment) error {
	payment, err := s.store.Payment().Find(ctx, filters.New().Add(filters.PaymentByID(sp.PaymentID)).Filters())
	if err != nil {
		s.log.Error("success balance payment", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	filter := filters.New().Add(filters.UserByID(payment.UserID))

	worker, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.ProjSuccessPayment)
	if err != nil {
		s.log.Error("success balance payment", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	if needBonus := s.updateReferral(ctx, worker); needBonus {
		worker.Balance.StarsBalance += bonusReferral
	}
	if worker.Balance != nil {
		worker.Balance.NumberPayments += 1
		worker.Balance.StarsBalance += payment.Amount
	} else {
		s.log.Error("success balance payment", slog.String("err", usermodel.ErrBalanceIsNil(worker.ID).Error()))
		return httpresponse.ServerError()
	}

	updPayment := &paymentmodel.Payment{
		ID:          payment.ID, // set only for transaction func
		TgPaymentID: sp.TgPaymentID,
		Status:      paymentmodel.Success,
	}
	err = s.successPaymentTrx(ctx, worker, updPayment)
	if err != nil {
		s.log.Error("success balance payment", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}
	return nil
}

// returns true if you need a bonus for the referral program and false if the bonus is not credited
func (s *PaymentServiceImpl) updateReferral(ctx context.Context, replenishment *usermodel.User) bool {
	switch {
	case replenishment.ReferralID == 0:
		return false
	case replenishment.Balance.NumberPayments > 0:
		return false
	}
	referral, err := s.store.User().FindProj(
		ctx,
		filters.New().Add(filters.UserByTgID(replenishment.ReferralID)).Filters(),
		usermodel.ProjOnlyBalance,
	)
	if err != nil {
		s.log.Error("find referral", slog.Int64("referralID", referral.ReferralID))
		return false
	}

	referral.Balance.StarsBalance += bonusReferral
	_, err = s.store.User().Edit(
		ctx,
		filters.New().Add(filters.UserByID(referral.ID)).Filters(),
		referral,
	)
	if err != nil {
		s.log.Error("update referral", slog.Int64("referralID", referral.ReferralID))
		return true
	}
	return true
}
