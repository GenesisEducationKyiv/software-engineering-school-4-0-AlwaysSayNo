package util_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"genesis-currency-api/internal/handler/util"

	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	router          *gin.Engine
	userGetterMock  *mocks.UserGetter
	emailSenderMock *mocks.EmailSender
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.userGetterMock = new(mocks.UserGetter)
	suite.emailSenderMock = new(mocks.EmailSender)

	suite.router.Use(middleware.ErrorHandler())

	utilHandler := util.NewHandler(suite.userGetterMock, suite.emailSenderMock)
	util.RegisterRoutes(suite.router, *utilHandler)
}

func (suite *HandlerTestSuite) TestFindAll_checkResult() {
	// SETUP
	users := []dto.UserResponseDTO{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}
	suite.userGetterMock.On("GetAll").Return(users, nil)

	req, _ := http.NewRequest("GET", "/api/v1/util/users", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)
	var responseBody []dto.UserResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	suite.Require().NoError(err)

	// VERIFY
	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal(len(responseBody), 2)
	suite.Equal(responseBody[0].Email, users[0].Email)
	suite.Equal(responseBody[1].Email, users[1].Email)

	suite.userGetterMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestFindAll_whenError() {
	// SETUP
	suite.userGetterMock.On("GetAll").Return(nil, errors.New("test"))

	req, _ := http.NewRequest("GET", "/api/v1/util/users", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusInternalServerError, resp.Code)

	suite.userGetterMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestSendEmails_checkResult() {
	// SETUP
	suite.emailSenderMock.On("SendEmails").Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/util/emails/send", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal("Emails are successfully sent", strings.ReplaceAll(resp.Body.String(), "\"", ""))

	suite.emailSenderMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestSendEmails_whenError() {
	// SETUP
	suite.emailSenderMock.On("SendEmails").Return(errors.New("test"))

	req, _ := http.NewRequest("POST", "/api/v1/util/emails/send", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusInternalServerError, resp.Code)

	suite.emailSenderMock.AssertExpectations(suite.T())
}
