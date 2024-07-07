package handler_test

import (
	"encoding/json"
	"errors"
	"genesis-currency-api/internal/user/api/handler"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/dto"
	myerrors "genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	userService *mocks.UserService
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.userService = new(mocks.UserService)
	suite.router.Use(middleware.ErrorHandler())

	userHandler := handler.NewHandler(suite.userService)
	handler.RegisterRoutes(suite.router, *userHandler)
}

func (suite *HandlerTestSuite) TestAdd_checkResult() {
	// SETUP
	saveDto := dto.UserSaveRequestDTO{
		Email: "test@example.com",
	}

	suite.userService.On("Save", saveDto).Return(&dto.UserResponseDTO{
		ID:    1,
		Email: "test@example.com",
	}, nil)

	form := url.Values{}
	form.Add("email", "test@example.com")
	req, _ := http.NewRequest("POST", "/api/v1/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "E-mail додано")

	suite.userService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestAdd_whenError() {
	// SETUP
	saveDto := dto.UserSaveRequestDTO{
		Email: "exist@example.com",
	}

	suite.userService.On("Save", saveDto).Return(nil, myerrors.NewUserWithEmailExistsError())

	form := url.Values{}
	form.Add("email", "exist@example.com")
	req, _ := http.NewRequest("POST", "/api/v1/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Contains(resp.Body.String(), "Повертати, якщо e-mail вже є в базі даних")

	suite.userService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestFindAll_checkResult() {
	// SETUP
	users := []dto.UserResponseDTO{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}
	suite.userService.On("GetAll").Return(users, nil)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
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

	suite.userService.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestFindAll_whenError() {
	// SETUP
	suite.userService.On("GetAll").Return(nil, errors.New("test"))

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusInternalServerError, resp.Code)

	suite.userService.AssertExpectations(suite.T())
}
