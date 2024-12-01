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

// TEST
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

// TEST
func (s *PaymentServiceImpl) SuccessBalancePayment(ctx context.Context, sp *paymentmodel.SuccessPayment) error {
	payment, err := s.store.Payment().Find(ctx, filters.New().Add(filters.PaymentByID(sp.PaymentID)).Filters())
	if err != nil {
		s.log.Error("success balance payment", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	filter := filters.New().Add(filters.UserByID(payment.UserID))

	worker, err := s.store.User().FindProj(ctx, filter.Filters(), usermodel.OnlyWorkerInfo)
	if err != nil {
		s.log.Error("success balance payment", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}

	if worker.WorkerInfo != nil {
		worker.WorkerInfo.StarsBalance += payment.Amount
	} else {
		s.log.Error("success balance payment", slog.String("err", usermodel.ErrWorkerInfoIsNil(worker.ID).Error()))
		return httpresponse.ServerError()
	}

	updPayment := &paymentmodel.Payment{
		ID:          payment.ID, // set only for transaction func
		TgPaymentID: sp.TgPaymentID,
		Status:      paymentmodel.Success,
	}
	err = s.successBalanePayTrx(ctx, worker, updPayment)
	if err != nil {
		s.log.Error("success balance payment", slog.String("err", err.Error()))
		return httpresponse.ServerError()
	}
	return nil
}
