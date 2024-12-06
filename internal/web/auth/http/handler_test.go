package authHttp_test

import (
	"os"
	"testing"

	"github.com/Negat1v9/work-marketplace/internal/config"
	mongo_mock "github.com/Negat1v9/work-marketplace/internal/storage/mock"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
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
