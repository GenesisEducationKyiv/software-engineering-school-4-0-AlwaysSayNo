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

type UtilControllerImplTestSuite struct {
	suite.Suite
	router           *gin.Engine
	mockUserService  *mocks.UserServiceInterface
	mockEmailService *mocks.EmailServiceInterface
}

func TestUtilControllerImplTestSuite(t *testing.T) {
	suite.Run(t, new(UtilControllerImplTestSuite))
}

func (suite *UtilControllerImplTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.mockUserService = new(mocks.UserServiceInterface)
	suite.mockEmailService = new(mocks.EmailServiceInterface)

	suite.router.Use(middleware.ErrorHandler())

	utilHandler := util.NewHandler(suite.mockUserService, suite.mockEmailService)
	util.RegisterRoutes(suite.router, *utilHandler)
}

func (suite *UtilControllerImplTestSuite) TestFindAll_checkResult() {
	// SETUP
	users := []dto.UserResponseDTO{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}
	suite.mockUserService.On("GetAll").Return(users, nil)

	req, _ := http.NewRequest("GET", "/api/util/users", nil)
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

	suite.mockUserService.AssertExpectations(suite.T())
}

func (suite *UtilControllerImplTestSuite) TestFindAll_whenError() {
	// SETUP
	suite.mockUserService.On("GetAll").Return(nil, errors.New("test"))

	req, _ := http.NewRequest("GET", "/api/util/users", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusInternalServerError, resp.Code)

	suite.mockUserService.AssertExpectations(suite.T())
}

func (suite *UtilControllerImplTestSuite) TestSendEmails_checkResult() {
	// SETUP
	suite.mockEmailService.On("SendEmails").Return(nil)

	req, _ := http.NewRequest("POST", "/api/util/emails/send", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusOK, resp.Code)
	suite.Equal("Emails are successfully sent", strings.ReplaceAll(resp.Body.String(), "\"", ""))

	suite.mockEmailService.AssertExpectations(suite.T())
}

func (suite *UtilControllerImplTestSuite) TestSendEmails_whenError() {
	// SETUP
	suite.mockEmailService.On("SendEmails").Return(errors.New("test"))

	req, _ := http.NewRequest("POST", "/api/util/emails/send", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusInternalServerError, resp.Code)

	suite.mockEmailService.AssertExpectations(suite.T())
}
