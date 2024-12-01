package mongoStore

import (
	"github.com/Negat1v9/work-marketplace/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection    = "user"
	taskCollection    = "task"
	paymentCollection = "payment"
	respondCollection = "respond"
	tgCmdCollection   = "tg_command"
	commentCollection = "comment"
)

type Store struct {
	db *mongo.Database
	*userRepository
	*taskRepository
	*paymentRepository
	*tgCmdRepository
	*respondRepository
	*commentRepository
}

func New(db *mongo.Database) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) StartSession() (mongo.Session, error) {
	session, err := s.db.Client().StartSession()

	return session, err
}

func (s *Store) User() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = newUserRepo(s.db.Collection(userCollection))
	return s.userRepository
}

func (s *Store) Task() storage.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}
	s.taskRepository = newTaskRepo(s.db.Collection(taskCollection))
	return s.taskRepository
}

func (s *Store) Payment() storage.PaymentRepository {
	if s.paymentRepository != nil {
		return s.paymentRepository
	}

	s.paymentRepository = newPaymentRepo(s.db.Collection(paymentCollection))

	return s.paymentRepository
}

func (s *Store) TgCmd() storage.TgCommandRepository {
	if s.tgCmdRepository != nil {
		return s.tgCmdRepository
	}
	s.tgCmdRepository = newTgCmdRepo(s.db.Collection(tgCmdCollection))
	return s.tgCmdRepository
}

func (s *Store) Respond() storage.RespondRepository {
	if s.respondRepository != nil {
		return s.respondRepository
	}
	s.respondRepository = newRespondRepo(s.db.Collection(respondCollection))
	return s.respondRepository
}

func (s *Store) Comment() storage.CommentRepository {
	if s.commentRepository != nil {
		return s.commentRepository
	}

	s.commentRepository = newCommentRepo(s.db.Collection(commentCollection))
	return s.commentRepository
}
