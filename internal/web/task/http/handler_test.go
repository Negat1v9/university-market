package taskHttp_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Negat1v9/work-marketplace/internal/config"
	sercicesMock "github.com/Negat1v9/work-marketplace/internal/services/mock"
	task_service_mock "github.com/Negat1v9/work-marketplace/internal/services/task/mock"
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	mongo_mock "github.com/Negat1v9/work-marketplace/internal/storage/mock"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	taskHttp "github.com/Negat1v9/work-marketplace/internal/web/task/http"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

var testUserID string
var testJwt string
var mockLogger = slog.Logger{}
var storeMock *mongo_mock.MockStore
var mockCfg *config.Config
var ctrl *gomock.Controller

func TestMain(m *testing.M) {
	ctrl = gomock.NewController(nil)

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

func TestGetUsersAllTasks(t *testing.T) {
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)

	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	taskservice := serviceMock.TaskService.(*task_service_mock.MockTaskService)
	taskHandler := taskHttp.New(mockCfg.WebConfig, taskservice)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	hander := taskHttp.RestTaskRouter(taskHandler, mw)

	userTasks := []taskmodel.Task{}

	taskservice.EXPECT().FindUserTasks(gomock.Any(), gomock.Any(), gomock.Any()).Return(userTasks, nil)
	req, err := http.NewRequest("GET", "/my", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

}

func TestFindOne(t *testing.T) {
	ctrl = gomock.NewController(t)

	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	taskservice := serviceMock.TaskService.(*task_service_mock.MockTaskService)
	taskHandler := taskHttp.New(mockCfg.WebConfig, serviceMock.TaskService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	hander := taskHttp.RestTaskRouter(taskHandler, mw)

	task := &taskmodel.InfoTaskRes{}
	taskservice.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(task, nil)

	req, err := http.NewRequest("GET", "/find/"+primitive.NewObjectID().Hex(), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)
	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestCreateTask(t *testing.T) {
	ctrl = gomock.NewController(t)

	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	taskservice := serviceMock.TaskService.(*task_service_mock.MockTaskService)
	taskHandler := taskHttp.New(mockCfg.WebConfig, serviceMock.TaskService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	hander := taskHttp.RestTaskRouter(taskHandler, mw)

	newTaskMeta := &taskmodel.TaskMeta{
		WithFiles:     false,
		MaxDaysWork:   0,
		MinPrice:      0,
		MaxPrice:      0,
		FormEducation: "back",
		University:    "University",
		Subject:       "subject",
		TaskType:      "type",
		Description:   "Description",
	}
	taskservice.EXPECT().Create(gomock.Any(), testUserID, newTaskMeta).Return("generated-id", nil)

	data, err := json.Marshal(&newTaskMeta)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/create", strings.NewReader(string(data)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code)
}

func TestUpdateMetaTask(t *testing.T) {
	ctrl = gomock.NewController(t)

	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	taskservice := serviceMock.TaskService.(*task_service_mock.MockTaskService)
	taskHandler := taskHttp.New(mockCfg.WebConfig, serviceMock.TaskService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	hander := taskHttp.RestTaskRouter(taskHandler, mw)

	updTask := &taskmodel.TaskMeta{
		WithFiles:     false,
		MaxDaysWork:   1,
		MinPrice:      100,
		MaxPrice:      1000,
		FormEducation: "updated",
		University:    "updated",
		Subject:       "updated",
		TaskType:      "updated",
		Description:   "updated",
	}

	task := &taskmodel.Task{}
	taskservice.EXPECT().UpdateTaskMeta(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(task, nil)

	data, err := json.Marshal(&updTask)
	assert.NoError(t, err)
	req, err := http.NewRequest("PUT", "/edit/meta/"+primitive.NewObjectID().Hex(), strings.NewReader(string(data)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestSelectWorker(t *testing.T) {
	ctrl = gomock.NewController(t)

	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	taskservice := serviceMock.TaskService.(*task_service_mock.MockTaskService)
	taskHandler := taskHttp.New(mockCfg.WebConfig, serviceMock.TaskService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	hander := taskHttp.RestTaskRouter(taskHandler, mw)

	workerID := primitive.NewObjectID().Hex()
	taskID := primitive.NewObjectID().Hex()
	taskInfo := &taskmodel.InfoTaskRes{}
	taskservice.EXPECT().SelectWorker(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(taskInfo, nil)

	req, err := http.NewRequest("PUT", "/"+taskID+"/select/worker/"+workerID, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestCompleteTask(t *testing.T) {
	ctrl = gomock.NewController(t)

	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	taskservice := serviceMock.TaskService.(*task_service_mock.MockTaskService)
	taskHandler := taskHttp.New(mockCfg.WebConfig, serviceMock.TaskService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	hander := taskHttp.RestTaskRouter(taskHandler, mw)

	taskID := primitive.NewObjectID().Hex()
	task := &taskmodel.InfoTaskRes{}
	taskservice.EXPECT().CompleteTask(gomock.Any(), taskID, testUserID).Return(task, nil)

	req, err := http.NewRequest("PUT", "/complete/"+taskID, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestDeleteTask(t *testing.T) {
	ctrl = gomock.NewController(t)

	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	taskservice := serviceMock.TaskService.(*task_service_mock.MockTaskService)
	taskHandler := taskHttp.New(mockCfg.WebConfig, serviceMock.TaskService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	hander := taskHttp.RestTaskRouter(taskHandler, mw)

	taskID := primitive.NewObjectID().Hex()
	taskservice.EXPECT().DeleteTask(gomock.Any(), taskID, testUserID).Return(nil)

	req, err := http.NewRequest("DELETE", "/delete/"+taskID, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
