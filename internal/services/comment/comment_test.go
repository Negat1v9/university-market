package commentservice

import (
	"context"
	"log/slog"
	"net/url"
	"testing"

	mongo_mock "github.com/Negat1v9/work-marketplace/internal/storage/mock"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCommentServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mongo_mock.NewMockStore(ctrl)
	mockTaskRepo := mongo_mock.NewMockTaskRepository(ctrl)
	mockCommentRepo := mongo_mock.NewMockCommentRepository(ctrl)
	mockLogger := slog.Logger{}
	mockSession := &mongo_mock.MockSession{}

	commentService := &CommentServiceImpl{
		log:   &mockLogger,
		store: mockStore,
	}

	ctx := context.Background()
	userID := "user123"
	taskID := "task123"
	comment := &commentmodel.Comment{
		TaskID:      taskID,
		TaskType:    "task_type",
		CreatorID:   userID,
		WorkerID:    "worker_id",
		IsLike:      true,
		Description: "description",
	}

	t.Run("Validation error in beforeCreate", func(t *testing.T) {
		temp := comment.TaskType
		comment.TaskType = ""
		commentID, err := commentService.Create(ctx, userID, comment)
		assert.Empty(t, commentID)
		assert.Error(t, err)
		comment.TaskType = temp
	})

	t.Run("Task not found", func(t *testing.T) {
		mockStore.EXPECT().Task().Return(mockTaskRepo).AnyTimes()
		mockTaskRepo.EXPECT().FindProj(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, mongoStore.ErrNoTask)

		commentID, err := commentService.Create(ctx, userID, comment)

		assert.Empty(t, commentID)
		assert.EqualError(t, err, httpresponse.NewError(404, mongoStore.ErrNoTask.Error()).Error())
	})

	t.Run("Task alredy commented", func(t *testing.T) {
		task := &taskmodel.Task{
			IsComment: true,
		}
		mockStore.EXPECT().Task().Return(mockTaskRepo).AnyTimes()
		mockTaskRepo.EXPECT().FindProj(gomock.Any(), gomock.Any(), gomock.Any()).Return(task, nil)

		commentID, err := commentService.Create(ctx, userID, comment)
		assert.Empty(t, commentID)
		assert.EqualError(t, err, httpresponse.NewError(409, "alredy commented").Error())
	})

	t.Run("Task assigned_to mismatch", func(t *testing.T) {
		task := &taskmodel.Task{
			IsComment:  false,
			AssignedTo: "anotherWorkerID",
			Status:     taskmodel.Completed,
		}
		mockStore.EXPECT().Task().Return(mockTaskRepo).AnyTimes()
		mockTaskRepo.EXPECT().FindProj(gomock.Any(), gomock.Any(), gomock.Any()).Return(task, nil)

		commentID, err := commentService.Create(ctx, userID, comment)
		assert.Empty(t, commentID)
		assert.EqualError(t, err, httpresponse.NewError(409, "task.assigned_to is not "+comment.WorkerID).Error())
	})

	t.Run("Successful creating", func(t *testing.T) {
		task := &taskmodel.Task{
			IsComment:  false,
			AssignedTo: comment.WorkerID,
			Status:     taskmodel.Completed,
		}
		mockStore.EXPECT().Task().Return(mockTaskRepo).AnyTimes()
		mockTaskRepo.EXPECT().FindProj(gomock.Any(), gomock.Any(), gomock.Any()).Return(task, nil)

		mockStore.EXPECT().StartSession().Return(mockSession, nil)

		mockStore.EXPECT().Comment().Return(mockCommentRepo).AnyTimes()
		mockCommentRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return("comment123", nil)

		mockTaskRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

		commentID, err := commentService.Create(ctx, userID, comment)
		assert.Equal(t, commentID, "comment123")
		assert.NoError(t, err)
	})
}

func TestUserComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mongo_mock.NewMockStore(ctrl)
	mockCommentRepo := mongo_mock.NewMockCommentRepository(ctrl)
	mockLogger := slog.Logger{}

	commentService := &CommentServiceImpl{
		log:   &mockLogger,
		store: mockStore,
	}
	ctx := context.Background()
	userID := "user_id"
	t.Run("No User Comments", func(t *testing.T) {
		mockStore.EXPECT().Comment().Return(mockCommentRepo).AnyTimes()

		mockCommentRepo.EXPECT().FindMany(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mongoStore.ErrNoComment)

		comments, err := commentService.UserComments(ctx, userID, url.Values{})
		assert.Len(t, comments, 0)
		assert.EqualError(t, err, httpresponse.NewError(404, mongoStore.ErrNoComment.Error()).Error())
	})
	t.Run("Successful get User comments", func(t *testing.T) {
		mockStore.EXPECT().Comment().Return(mockCommentRepo).AnyTimes()

		mockCommentRepo.EXPECT().FindMany(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]commentmodel.Comment, 5), nil)

		comments, err := commentService.UserComments(ctx, userID, url.Values{})
		assert.Len(t, comments, 5)
		assert.NoError(t, err)
	})
}

func TestWorkerComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mongo_mock.NewMockStore(ctrl)
	mockCommentRepo := mongo_mock.NewMockCommentRepository(ctrl)
	mockLogger := slog.Logger{}

	commentService := &CommentServiceImpl{
		log:   &mockLogger,
		store: mockStore,
	}
	ctx := context.Background()
	workerID := "worker_id"
	t.Run("No Worker Comments", func(t *testing.T) {
		mockStore.EXPECT().Comment().Return(mockCommentRepo).AnyTimes()

		mockCommentRepo.EXPECT().FindMany(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mongoStore.ErrNoComment)

		comments, err := commentService.WorkerComments(ctx, workerID, url.Values{})
		assert.Len(t, comments, 0)
		assert.EqualError(t, err, httpresponse.NewError(404, mongoStore.ErrNoComment.Error()).Error())
	})
	t.Run("Successful get Worker comments", func(t *testing.T) {
		mockStore.EXPECT().Comment().Return(mockCommentRepo).AnyTimes()

		mockCommentRepo.EXPECT().FindMany(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]commentmodel.Comment, 10), nil)

		comments, err := commentService.WorkerComments(ctx, workerID, url.Values{})
		assert.Len(t, comments, 10)
		assert.NoError(t, err)
	})
}
