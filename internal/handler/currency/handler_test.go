package currency_test

import (
	"fmt"
	"genesis-currency-api/internal/middleware"
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

type HandlerTestSuite struct {
	suite.Suite
	router    *gin.Engine
	raterMock *mocks.Rater
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.raterMock = new(mocks.Rater)
	suite.router.Use(middleware.ErrorHandler())

	currencyHandler := currency.NewHandler(suite.raterMock)
	currency.RegisterRoutes(suite.router, *currencyHandler)
}

func (suite *HandlerTestSuite) TestGetLatest_checkResult() {
	// SETUP
	rate := dto.CurrencyResponseDTO{
		Number: 39.35,
	}
	suite.raterMock.On("GetCurrencyRate").Return(rate, nil)

	req, _ := http.NewRequest("GET", "/api/v1/rate/", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)
	actual, _ := strconv.ParseFloat(resp.Body.String(), 64)

	// VERIFY
	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(rate.Number, actual)

	suite.raterMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetLatest_whenReturnError() {
	// SETUP
	suite.raterMock.On("GetCurrencyRate").
		Return(dto.CurrencyResponseDTO{}, fmt.Errorf("test error"))

	req, _ := http.NewRequest("GET", "/api/v1/rate/", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusInternalServerError, resp.Code)
	suite.Contains(resp.Body.String(), "unknown error")

	suite.raterMock.AssertExpectations(suite.T())
}
