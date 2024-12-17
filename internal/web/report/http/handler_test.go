package reportHttp_test

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
	report_service_mock "github.com/Negat1v9/work-marketplace/internal/services/report/mock"
	user_service_mock "github.com/Negat1v9/work-marketplace/internal/services/user/mock"
	mongo_mock "github.com/Negat1v9/work-marketplace/internal/storage/mock"
	"github.com/Negat1v9/work-marketplace/internal/web/middleware"
	reportHttp "github.com/Negat1v9/work-marketplace/internal/web/report/http"
	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
	"github.com/Negat1v9/work-marketplace/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
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

func TestCreateReportOnWorker(t *testing.T) {
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)

	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().Auth(gomock.Any(), testUserID).Return(nil)

	reportService := serviceMock.ReportService.(*report_service_mock.MockReportService)
	reportHandler := reportHttp.New(mockCfg.WebConfig, reportService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	handler := reportHttp.RestReportRouter(reportHandler, mw)

	report := &reportmodel.NewReportReq{
		ReportUser:  "id",
		Reason:      "reason",
		Description: "Description",
	}

	reportService.EXPECT().CreateReportOnWorker(gomock.Any(), gomock.Any(), report).Return("generated-id", nil)

	data, err := json.Marshal(&report)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/user/create", strings.NewReader(string(data)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code)
}

func TestCreateReportOnUser(t *testing.T) {
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := sercicesMock.NewServiceMockBuilder(ctrl)

	userServiceMock := serviceMock.UserService.(*user_service_mock.MockUserService)
	userServiceMock.EXPECT().AuthWorker(gomock.Any(), testUserID).Return(nil)

	reportService := serviceMock.ReportService.(*report_service_mock.MockReportService)
	reportHandler := reportHttp.New(mockCfg.WebConfig, reportService)

	mw := middleware.New(mockCfg.WebConfig, serviceMock.UserService)
	handler := reportHttp.RestReportRouter(reportHandler, mw)

	report := &reportmodel.NewReportReq{
		ReportUser:  "id",
		Reason:      "reason",
		Description: "Description",
	}

	reportService.EXPECT().CreateReportOnUser(gomock.Any(), gomock.Any(), report).Return("generated-id", nil)

	data, err := json.Marshal(&report)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/worker/create", strings.NewReader(string(data)))
	assert.NoError(t, err)
	req.Header.Set("Authorization", testJwt)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code)
}
