package commentHttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Negat1v9/work-marketplace/internal/config"
	comment_service_mock "github.com/Negat1v9/work-marketplace/internal/services/comment/mock"
	sercicesMock "github.com/Negat1v9/work-marketplace/internal/services/mock"
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	mongo_mock "github.com/Negat1v9/work-marketplace/internal/storage/mock"
	commentHttp "github.com/Negat1v9/work-marketplace/internal/web/comment"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	commentmodel "github.com/Negat1v9/work-marketplace/model/comment"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

var testJwt string
var storeMock *mongo_mock.MockStore
var mockCfg *config.Config
var testUserID string

func TestMain(m *testing.M) {
	ctrl := gomock.NewController(nil)

	defer ctrl.Finish()

	mockCfg = config.NewConfigMock()
	storeMock = mongo_mock.NewMockStore(ctrl)
	testUserID = primitive.NewObjectID().Hex()

	claims := &utils.Claims{
		UserID: testUserID,
	}
	testJwt, _ = utils.GenerateJwtToken(claims, mockCfg.WebConfig.JwtSecret)

	code := m.Run()
	os.Exit(code)
}

func TestCreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := sercicesMock.NewServiceMockBuilder(ctrl)
	userService := services.UserService.(*user_service_mock.MockUserService)

	userService.EXPECT().Auth(gomock.Any(), gomock.Any()).Return(nil)

	commentService := services.CommentService.(*comment_service_mock.MockCommentService)

	commentHandler := commentHttp.New(mockCfg.WebConfig, commentService)

	mw := middleware.New(mockCfg.WebConfig, userService)

	handler := commentHttp.RestPaymentRouter(commentHandler, mw)

	commentID := "comment_id"
	data := &commentmodel.Comment{
		ID:          commentID,
		TaskID:      "id",
		TaskType:    "type",
		CreatorID:   testUserID,
		WorkerID:    testUserID,
		IsLike:      true,
		Description: "like",
	}

	body, err := json.Marshal(&data)
	assert.NoError(t, err)

	commentService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(commentID, nil)

	req, err := http.NewRequest("POST", "/user/create", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code)

}

func TestUserComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := sercicesMock.NewServiceMockBuilder(ctrl)
	userService := services.UserService.(*user_service_mock.MockUserService)

	userService.EXPECT().Auth(gomock.Any(), gomock.Any()).Return(nil)

	commentService := services.CommentService.(*comment_service_mock.MockCommentService)

	commentHandler := commentHttp.New(mockCfg.WebConfig, commentService)

	mw := middleware.New(mockCfg.WebConfig, userService)

	handler := commentHttp.RestPaymentRouter(commentHandler, mw)

	comments := []commentmodel.Comment{}
	commentService.EXPECT().UserComments(gomock.Any(), gomock.Any(), gomock.Any()).Return(comments, nil)

	req, err := http.NewRequest("GET", "/user/my", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestUserWorkerComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := sercicesMock.NewServiceMockBuilder(ctrl)
	userService := services.UserService.(*user_service_mock.MockUserService)

	userService.EXPECT().Auth(gomock.Any(), gomock.Any()).Return(nil)

	commentService := services.CommentService.(*comment_service_mock.MockCommentService)

	commentHandler := commentHttp.New(mockCfg.WebConfig, commentService)

	mw := middleware.New(mockCfg.WebConfig, userService)

	handler := commentHttp.RestPaymentRouter(commentHandler, mw)

	comments := []commentmodel.Comment{}
	commentService.EXPECT().WorkerComments(gomock.Any(), gomock.Any(), gomock.Any()).Return(comments, nil)

	req, err := http.NewRequest("GET", "/user/worker/id", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
func TestWokrerComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := sercicesMock.NewServiceMockBuilder(ctrl)
	userService := services.UserService.(*user_service_mock.MockUserService)

	userService.EXPECT().AuthWorker(gomock.Any(), gomock.Any()).Return(nil)

	commentService := services.CommentService.(*comment_service_mock.MockCommentService)

	commentHandler := commentHttp.New(mockCfg.WebConfig, commentService)

	mw := middleware.New(mockCfg.WebConfig, userService)

	handler := commentHttp.RestPaymentRouter(commentHandler, mw)

	comments := []commentmodel.Comment{}
	commentService.EXPECT().WorkerComments(gomock.Any(), gomock.Any(), gomock.Any()).Return(comments, nil)

	req, err := http.NewRequest("GET", "/worker/my", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
