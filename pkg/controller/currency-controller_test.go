package controller_test

import (
	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/controller"
	"genesis-currency-api/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type CurrencyControllerImplTestSuite struct {
	suite.Suite
	router              *gin.Engine
	mockCurrencyService *mocks.CurrencyService
}

func TestCurrencyControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyControllerImplTestSuite))
}

func (suite *CurrencyControllerImplTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.mockCurrencyService = new(mocks.CurrencyService)
	controller.RegisterCurrencyRoutes(suite.router, suite.mockCurrencyService)
}

func (suite *CurrencyControllerImplTestSuite) TestGetLatest_checkResult() {
	// SETUP
	rate := dto.CurrencyResponseDto{
		Number: 39.35,
	}
	suite.mockCurrencyService.On("GetCurrencyRate").Return(rate, nil)

	req, _ := http.NewRequest("GET", "/api/rate/", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)
	actual, _ := strconv.ParseFloat(resp.Body.String(), 64)

	// VERIFY
	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(rate.Number, actual)

	suite.mockCurrencyService.AssertExpectations(suite.T())
}
