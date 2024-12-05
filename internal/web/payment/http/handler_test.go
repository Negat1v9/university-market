package paymentHttp_test

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
	payment_service_mock "github.com/Negat1v9/work-marketplace/internal/services/payment/mock"
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	mongoStore "github.com/Negat1v9/work-marketplace/internal/storage/mongo"
	botmock "github.com/Negat1v9/work-marketplace/internal/tgBot/bot/mock"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	paymentHttp "github.com/Negat1v9/work-marketplace/internal/web/payment/http"
	paymentmodel "github.com/Negat1v9/work-marketplace/model/payment"
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

func TestCreateInvoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := sercicesMock.NewServiceMockBuilder(ctrl)
	userService := services.UserService.(*user_service_mock.MockUserService)
	userService.EXPECT().AuthWorker(gomock.Any(), gomock.Any()).Return(nil)

	paymentService := services.PaymentService.(*payment_service_mock.MockPaymentService)

	paymentHandler := paymentHttp.New(mockCfg.WebConfig, paymentService)

	mw := middleware.New(mockCfg.WebConfig, userService)

	handler := paymentHttp.RestPaymentRouter(paymentHandler, mw)

	paymentLink := "https://example.com"
	paymentService.EXPECT().CreateInvoiceLink(gomock.Any(), gomock.Any(), gomock.Any()).Return(paymentLink, nil)

	rr := httptest.NewRecorder()

	data := paymentmodel.PaymentLinkReq{
		Amount: 100,
	}

	body, err := json.Marshal(&data)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/create", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
