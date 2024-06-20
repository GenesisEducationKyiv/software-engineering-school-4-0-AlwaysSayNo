package currency_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"genesis-currency-api/internal/handler/currency"

	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type CurrencyHandlerTestSuite struct {
	suite.Suite
	router              *gin.Engine
	mockCurrencyService *mocks.CurrencyServiceInterface
}

func TestCurrencyHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyHandlerTestSuite))
}

func (suite *CurrencyHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.mockCurrencyService = new(mocks.CurrencyServiceInterface)

	currencyHandler := currency.NewHandler(suite.mockCurrencyService)
	currency.RegisterRoutes(suite.router, *currencyHandler)
}

func (suite *CurrencyHandlerTestSuite) TestGetLatest_checkResult() {
	// SETUP
	rate := dto.CurrencyResponseDTO{
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
