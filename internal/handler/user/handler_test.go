package user_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"genesis-currency-api/internal/handler/user"

	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	router    *gin.Engine
	saverMock *mocks.Saver
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	suite.saverMock = new(mocks.Saver)
	suite.router.Use(middleware.ErrorHandler())

	userHandler := user.NewHandler(suite.saverMock)
	user.RegisterRoutes(suite.router, *userHandler)
}

func (suite *HandlerTestSuite) TestAdd_checkResult() {
	// SETUP
	saveDto := dto.UserSaveRequestDTO{
		Email: "test@example.com",
	}

	suite.saverMock.On("Save", saveDto).Return(&dto.UserResponseDTO{
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

	suite.saverMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestAdd_whenError() {
	// SETUP
	saveDto := dto.UserSaveRequestDTO{
		Email: "exist@example.com",
	}

	suite.saverMock.On("Save", saveDto).Return(nil, errors.NewUserWithEmailExistsError())

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

	suite.saverMock.AssertExpectations(suite.T())
}
