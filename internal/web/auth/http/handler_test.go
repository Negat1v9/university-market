package authHttp_test

import (
	"os"
	"testing"

	"github.com/Negat1v9/work-marketplace/internal/config"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testJwt string
var storeMock *mongoStore.StoreMock
var mockCfg *config.Config
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
