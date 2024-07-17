package handler_test

//
//import (
//	"errors"
//	handler2 "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/api/handler/handler"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//
//	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/middleware"
//	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/mocks"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/suite"
//)
//
//type HandlerTestSuite struct {
//	suite.Suite
//	router          *gin.Engine
//	emailSenderMock *mocks.EmailSender
//}
//
//func TestHandlerTestSuite(t *testing.T) {
//	suite.Run(t, new(HandlerTestSuite))
//}
//
//func (suite *HandlerTestSuite) SetupTest() {
//	gin.SetMode(gin.TestMode)
//	suite.router = gin.New()
//
//	suite.emailSenderMock = new(mocks.EmailSender)
//
//	suite.router.Use(middleware.ErrorHandler())
//
//	utilHandler := handler2.NewHandler(suite.emailSenderMock)
//	handler2.RegisterRoutes(suite.router, utilHandler)
//}
//
//func (suite *HandlerTestSuite) TestSendEmails_checkResult() {
//	// SETUP
//	suite.emailSenderMock.On("SendEmails").Return(nil)
//
//	req, _ := http.NewRequest("POST", "/api/v1/emails/send", nil)
//	resp := httptest.NewRecorder()
//
//	// ACT
//	suite.router.ServeHTTP(resp, req)
//
//	// VERIFY
//	suite.Equal(http.StatusOK, resp.Code)
//	suite.Equal("Emails are successfully sent", strings.ReplaceAll(resp.Body.String(), "\"", ""))
//
//	suite.emailSenderMock.AssertExpectations(suite.T())
//}
//
//func (suite *HandlerTestSuite) TestSendEmails_whenError() {
//	// SETUP
//	suite.emailSenderMock.On("SendEmails").Return(errors.New("test"))
//
//	req, _ := http.NewRequest("POST", "/api/v1/emails/send", nil)
//	resp := httptest.NewRecorder()
//
//	// ACT
//	suite.router.ServeHTTP(resp, req)
//
//	// VERIFY
//	suite.Equal(http.StatusInternalServerError, resp.Code)
//
//	suite.emailSenderMock.AssertExpectations(suite.T())
//}
