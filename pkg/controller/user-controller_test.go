package controller_test

import (
	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/controller"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type UserControllerImplTestSuite struct {
	suite.Suite
	router           *gin.Engine
	mockUserService  *mocks.UserService
	mockEmailService *mocks.EmailService
}

func TestUserControllerImplTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerImplTestSuite))
}

func (suite *UserControllerImplTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.mockUserService = new(mocks.UserService)
	suite.mockEmailService = new(mocks.EmailService)

	suite.router.Use(middleware.ErrorHandler())
	controller.RegisterUserRoutes(suite.router, suite.mockUserService, suite.mockEmailService)
}

func (suite *UserControllerImplTestSuite) TestAdd_checkResult() {
	// SETUP
	saveDto := dto.UserSaveRequestDTO{
		Email: "test@example.com",
	}

	suite.mockUserService.On("Save", saveDto).Return(dto.UserResponseDTO{
		ID:    1,
		Email: "test@example.com",
	}, nil)

	form := url.Values{}
	form.Add("email", "test@example.com")
	req, _ := http.NewRequest("POST", "/api/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// Перевірка результату
	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "E-mail додано")

	suite.mockUserService.AssertExpectations(suite.T())
}

func (suite *UserControllerImplTestSuite) TestAdd_whenError() {
	// SETUP
	saveDto := dto.UserSaveRequestDTO{
		Email: "exist@example.com",
	}

	suite.mockUserService.On("Save", saveDto).Return(dto.UserResponseDTO{},
		errors.NewUserWithEmailExistsError())

	form := url.Values{}
	form.Add("email", "exist@example.com")
	req, _ := http.NewRequest("POST", "/api/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// ACT
	suite.router.ServeHTTP(resp, req)

	// VERIFY
	suite.Equal(http.StatusBadRequest, resp.Code)
	suite.Contains(resp.Body.String(), "Повертати, якщо e-mail вже є в базі даних")

	suite.mockUserService.AssertExpectations(suite.T())
}
