package usermodel

import (
	"time"

	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
)

type UserType string

const (
	RegularUser UserType = "regular"
	Worker      UserType = "worker"
)

type WorkerInfo struct {
	Ban          bool   `bson:"ban,omitempty" json:"ban,omitempty"`                     // does the worker have a ban?
	FullName     string `bson:"full_name,omitempty" json:"full_name,omitempty"`         // real name
	Education    string `bson:"education,omitempty" json:"education,omitempty"`         //level of education worker
	Experience   string `bson:"experience,omitempty" json:"experience,omitempty"`       // experience
	StarsBalance int    `bson:"stars_balance,omitempty" json:"stars_balance,omitempty"` // balance of telegram stars
	Description  string `bson:"description,omitempty" json:"description,omitempty"`     // Description of the worker
}

type NotificationSettings struct {
	SendSilent bool `bson:"send_silent,omitempty" json:"send_silent,omitempty"`
}

type User struct {
	ID           string                `bson:"_id,omitempty" json:"id"`
	TelegramID   int64                 `bson:"telegram_id,omitempty" json:"-"`                       // unique identifier from telegram
	Username     string                `bson:"username,omitempty" json:"username,omitempty"`         // nickname name from telegram
	PhoneNumber  string                `bson:"phone_number,omitempty" json:"phone_number,omitempty"` // phone number
	ReferralID   int64                 `bson:"referral_id,omitempty" json:"referral_id,omitempty"`   // ID of the user who invited him
	Notification *NotificationSettings `bson:"notification,omitempty" json:"notification,omitempty"`
	Role         UserType              `bson:"role,omitempty" json:"-"`                            // user role
	WorkerInfo   *WorkerInfo           `bson:"worker_info,omitempty" json:"worker_info,omitempty"` // worker information if the user is an worker
	CreatedAt    time.Time             `bson:"created_at,omitempty" json:"created_at"`             // Creation date
	UpdatedAt    time.Time             `bson:"updated_at,omitempty" json:"updated_at"`             // Updated date
}

func NewUser(tgId int64, username, fullname string, referall int64) *User {
	return &User{
		TelegramID: tgId,
		Username:   username,
		ReferralID: referall,
		Notification: &NotificationSettings{
			SendSilent: false,
		},
		Role:       RegularUser,
		WorkerInfo: NewWorkerInfo(0, fullname),
		CreatedAt:  time.Now().UTC(),
	}
}

func NewWorkerInfo(starsBalance int, fullName string) *WorkerInfo {
	return &WorkerInfo{
		Ban:          false,
		FullName:     fullName,
		StarsBalance: starsBalance,
		Description:  "",
	}
}

// data from tg bot
type UserCreate struct {
	TgID         int64
	TgReferralID int64
}

type WorkerCreate struct {
	PhoneNumber string `json:"phone_number"`
}

type UserLoginReq struct {
	InitData string `json:"init_data"`
}

type LoginRes struct {
	TokenType string `json:"token_type"`
	Token     string `json:"token"`
}

type WorkerInfoWithTaskRes struct {
	ID          string                               `json:"id"`
	FullName    string                               `json:"full_name"`
	Rating      *commentmodel.CountLikeDislikeWorker `json:"rating"`
	Education   string                               `json:"education,omitempty"`
	Experience  string                               `json:"experience,omitempty"`
	Description string                               `json:"description"`
}

type IsWorkerRes struct {
	UserID   string `json:"user_id"`
	IsWorker bool   `json:"is_worker"`
}
