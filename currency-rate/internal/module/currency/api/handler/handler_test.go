package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/middleware"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/currency/api/handler"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/mocks"
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

	currencyHandler := handler.NewHandler(suite.raterMock)
	handler.RegisterRoutes(suite.router, currencyHandler)
}

func (suite *HandlerTestSuite) TestGetLatest_checkResult() {
	// SETUP
	rate := sharcurrdto.ResponseDTO{
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
		Return(sharcurrdto.ResponseDTO{}, fmt.Errorf("test error"))

	req, _ := http.NewRequest("GET", "/api/v1/rate/", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusInternalServerError, resp.Code)
	suite.Contains(resp.Body.String(), "unknown error")

	suite.raterMock.AssertExpectations(suite.T())
}
