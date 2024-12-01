package paymentservice

import (
	"context"

	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
)

type PaymentService interface {
	CreateInvoiceLink(ctx context.Context, workerID string, data *paymentmodel.PaymentLinkReq) (string, error)
	SuccessBalancePayment(ctx context.Context, sp *paymentmodel.SuccessPayment) error
}
