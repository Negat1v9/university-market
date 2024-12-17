package mongoStore

import (
	"context"
	"log"

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
	reportCollection  = "report"
)

type Store struct {
	db *mongo.Database
	*userRepository
	*taskRepository
	*paymentRepository
	*tgCmdRepository
	*respondRepository
	*commentRepository
	*reportRepository
}

func New(db *mongo.Database) *Store {
	s := &Store{
		db: db,
	}
	s.userIndex()
	s.taskIndex()
	s.commentIndex()
	s.respondIndex()
	s.reportIndex()
	return s
}

func (s *Store) StartSession() (mongo.Session, error) {
	session, err := s.db.Client().StartSession()

	return session, err
}

func (s *Store) userIndex() {
	s.userRepository = newUserRepo(s.db.Collection(userCollection))
	indxs, err := s.userRepository.createIndexes(context.Background())
	if err != nil {
		log.Printf("error create user indexes %v", err)
	} else {
		log.Printf("create user indexes %v", indxs)
	}

}
func (s *Store) taskIndex() {
	s.taskRepository = newTaskRepo(s.db.Collection(taskCollection))
	indxs, err := s.taskRepository.createIndexes(context.Background())
	if err != nil {
		log.Printf("error create task indexes %v", err)
	} else {
		log.Printf("create task indexes %v", indxs)
	}
}
func (s *Store) commentIndex() {
	s.commentRepository = newCommentRepo(s.db.Collection(commentCollection))
	indxs, err := s.commentRepository.createIndexes(context.Background())
	if err != nil {
		log.Printf("error create comment indexes %v", err)
	} else {
		log.Printf("create comment indexes %v", indxs)
	}
}

func (s *Store) respondIndex() {
	s.respondRepository = newRespondRepo(s.db.Collection(respondCollection))
	indxs, err := s.respondRepository.createIndexes(context.Background())
	if err != nil {
		log.Printf("error create respond indexes %v", err)
	} else {
		log.Printf("create respond indexes %v", indxs)
	}
}
func (s *Store) reportIndex() {
	s.reportRepository = newReportRepo(s.db.Collection(reportCollection))
	indxs, err := s.reportRepository.createIndexes(context.Background())
	if err != nil {
		log.Printf("error create report indexes %v", err)
	} else {
		log.Printf("create report indexes %v", indxs)
	}
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

func (s *Store) Report() storage.ReportRepository {
	if s.reportRepository != nil {
		return s.reportRepository
	}
	s.reportRepository = newReportRepo(s.db.Collection(reportCollection))
	return s.reportRepository
}
