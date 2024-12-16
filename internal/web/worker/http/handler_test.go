package workerHttp_test

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
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	worker_service_mock "github.com/Negat1v9/work-marketplace/internal/services/worker/mock"
	mongo_mock "github.com/Negat1v9/work-marketplace/internal/storage/mock"
	botmock "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/mock"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	workerHttp "github.com/Negat1v9/work-marketplace/internal/web/worker/http"
	respondmodel "github.com/Negat1v9/work-marketplace/model/respond"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

var testJwt string
var mockLogger = slog.Logger{}
var storeMock *mongo_mock.MockStore
var mockCfg *config.Config
var botClientMock = &botmock.WebTgClientMock{}
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

func TestNewWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	dataNewWorker := &usermodel.WorkerCreate{
		Response: "fdfkdskfsodfksdok",
	}

	body, err := json.Marshal(&dataNewWorker)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/user/new", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	worker := &usermodel.User{}
	workerServMock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(worker, nil)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code)
}

func TestIsWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), gomock.Any()).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("POST", "/user/isworker", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	workerServMock.EXPECT().IsWorker(gomock.Any(), gomock.Any()).Return(true, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestPublicWorkerInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	// userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)
	userServiceMock.EXPECT().Auth(gomock.Any(), gomock.Any()).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("GET", "/user/worker/id", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	workerPublic := &usermodel.WorkerInfoWithTaskRes{}
	workerServMock.EXPECT().WorkerPublicInfo(gomock.Any(), "id").Return(workerPublic, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestWorkerProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("GET", "/profile", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	worker := &usermodel.User{}
	workerServMock.EXPECT().Worker(gomock.Any(), gomock.Any()).Return(worker, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
func TestUpdateInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	data := usermodel.WorkerInfo{
		Education: "NY UNIVERSITY",
	}
	body, err := json.Marshal(&data)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/edit/info", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	worker := &usermodel.User{}
	workerServMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(worker, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestAvailableTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("GET", "/task/all", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	tasks := []taskmodel.Task{}
	workerServMock.EXPECT().AvailableTasks(gomock.Any(), gomock.Any(), gomock.Any()).Return(tasks, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestTaskInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("GET", "/task/info/id", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	task := &taskmodel.InfoTaskRes{}
	workerServMock.EXPECT().TaskInfo(gomock.Any(), gomock.Any(), "id").Return(task, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestSendFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("PUT", "/task/files/id", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	workerServMock.EXPECT().SendTaskFiles(gomock.Any(), testUserID, "id").Return(nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestRespondOnTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("POST", "/task/respond/id", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	workerServMock.EXPECT().RespondOnTask(gomock.Any(), gomock.Any(), "id").Return(nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestTaskByResponded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("GET", "/task/responds", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	tasks := []taskmodel.Task{}

	workerServMock.EXPECT().TasksResponded(gomock.Any(), gomock.Any(), gomock.Any()).Return(tasks, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func TestGetResponds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	workerServMock := serviceMock.WorkerService.(*worker_service_mock.MockWorkerService)
	workerHandler := workerHttp.New(mockCfg.WebConfig, workerServMock)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := workerHttp.RestWorkerRouter(workerHandler, mw)

	req, err := http.NewRequest("GET", "/responds", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()

	responds := []respondmodel.Respond{}

	workerServMock.EXPECT().Responds(gomock.Any(), gomock.Any(), gomock.Any()).Return(responds, nil)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
