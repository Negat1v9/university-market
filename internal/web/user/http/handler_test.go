package userHttp_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Negat1v9/work-marketplace/internal/config"
	sercicesMock "github.com/Negat1v9/work-marketplace/internal/services/mock"
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	botmock "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/mock"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	userHttp "github.com/Negat1v9/work-marketplace/internal/web/user/http"
	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

var testJwt string
var mockLogger = slog.Logger{}
var storeMock *mongoStore.StoreMock
var mockCfg *config.Config
var botClientMock = &botmock.WebTgClientMock{}
var testUserID string

func TestMain(m *testing.M) {
	ctrl := gomock.NewController(nil)

	defer ctrl.Finish()

	mockCfg = config.NewConfigMock()
	storeMock = mongoStore.NewMockStore(ctrl)
	testUserID = primitive.NewObjectID().Hex()

	claims := &utils.Claims{
		UserID: testUserID,
	}
	testJwt, _ = utils.GenerateJwtToken(claims, mockCfg.WebConfig.JwtSecret)

	code := m.Run()
	os.Exit(code)
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)
	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	userHandler := userHttp.New(mockCfg.WebConfig, serviceMock.UserService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)

	handler := userHttp.RestUserRouter(userHandler, mw)

	req, err := http.NewRequest("GET", "/info/id", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	user := &usermodel.User{}
	userServiceMock.EXPECT().User(gomock.Any(), gomock.Any()).Return(user, nil)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
