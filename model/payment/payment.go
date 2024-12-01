package paymentmodel

import "time"

type PaymentStatus string

var (
	// status pending payment
	Pending PaymentStatus = "pnd"
	// status success payment
	Success PaymentStatus = "scc"
)

type Payment struct {
	ID          string        `bson:"_id,omitempty"`        // mongo id
	TgPaymentID string        `bson:"tg_payment_id"`        // payment id get from telegram, unique (only after success payment before is empty)
	Amount      int           `bson:"amount,omitempty"`     // amount telegram star
	UserID      string        `bson:"user_id,omitempty"`    // user who pay
	Status      PaymentStatus `bson:"status,omitempty"`     // status of payment from pending to success
	CreatedAt   time.Time     `bson:"created_at,omitempty"` // time when payment was created UTC format
}

type SuccessPayment struct {
	PaymentID   string
	Amount      int
	TgPaymentID string
}

func NewPayment(amount int, userID string) *Payment {
	return &Payment{
		Amount:    amount,
		UserID:    userID,
		Status:    Pending,
		CreatedAt: time.Now(),
	}
}

type PaymentLinkReq struct {
	Amount int `json:"amount"`
}

type PaymentLinkRes struct {
	Link string `json:"link"`
}
