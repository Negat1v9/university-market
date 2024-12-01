package storage

import "go.mongodb.org/mongo-driver/mongo"

type Store interface {
	StartSession() (mongo.Session, error)
	User() UserRepository
	Task() TaskRepository
	Payment() PaymentRepository
	TgCmd() TgCommandRepository
	Respond() RespondRepository
	Comment() CommentRepository
}
