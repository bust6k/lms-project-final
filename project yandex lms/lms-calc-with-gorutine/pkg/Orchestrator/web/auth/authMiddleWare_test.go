package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDBService struct {
	mock.Mock
}

func (m *MockDBService) VerifyUserExists(userID string) bool {
	args := m.Called(userID)
	return args.Bool(0)
}

type MockTokenParser struct {
	mock.Mock
}

func (m *MockTokenParser) ExtractUserID(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0), args.Error(1)
}

type MockTokenRefresher struct {
	mock.Mock
}

func (m *MockTokenRefresher) RefreshUserTokens(userID string, c *gin.Context) {
	m.Called(userID, c)
}

type MiddlewareTestSuite struct {
	suite.Suite
	dbMock       *MockDBService
	tokenMock    *MockTokenParser
	refreshMock  *MockTokenRefresher
	testRouter   *gin.Engine
	testRecorder *httptest.ResponseRecorder
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.dbMock = new(MockDBService)
	s.tokenMock = new(MockTokenParser)
	s.refreshMock = new(MockTokenRefresher)
	s.testRecorder = httptest.NewRecorder()
	gin.SetMode(gin.ReleaseMode)
	s.testRouter = gin.Default()

}

func (s *MiddlewareTestSuite) TestAuthMiddleware_ValidToken() {
	c, _ := gin.CreateTestContext(s.testRecorder)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("jwtAccsess", "valid_token")
	c.Set("user_id", "user123")

	s.tokenMock.On("ExtractUserID", "valid_token").Return("user123", nil)
	s.dbMock.On("VerifyUserExists", "user123").Return(true)

	handler := AutorizationMiddleWare()
	handler(c)

	s.False(c.IsAborted())
	s.tokenMock.AssertExpectations(s.T())
	s.dbMock.AssertExpectations(s.T())
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	c, _ := gin.CreateTestContext(s.testRecorder)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("jwtAccsess", "invalid_token")

	s.tokenMock.On("ExtractUserID", "invalid_token").Return("", errors.New("invalid token"))

	handler := AutorizationMiddleWare()
	handler(c)

	s.True(c.IsAborted())
	s.Equal(http.StatusUnauthorized, s.testRecorder.Code)
	s.tokenMock.AssertExpectations(s.T())
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_RefreshFlow() {
	c, _ := gin.CreateTestContext(s.testRecorder)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "Refresh", Value: "refresh_token"})
	c.Set("accsessInvalid", true)

	s.tokenMock.On("ExtractUserID", "refresh_token").Return("user123", nil)
	s.refreshMock.On("RefreshUserTokens", "user123", c)

	handler := AutorizationMiddleWare()
	handler(c)

	s.False(c.IsAborted())
	s.tokenMock.AssertExpectations(s.T())
	s.refreshMock.AssertExpectations(s.T())
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_MissingUserID() {
	c, _ := gin.CreateTestContext(s.testRecorder)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("jwtAccsess", "valid_token")

	handler := AutorizationMiddleWare()
	handler(c)

	s.True(c.IsAborted())
	s.Equal(http.StatusUnauthorized, s.testRecorder.Code)
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_UserNotInDB() {
	c, _ := gin.CreateTestContext(s.testRecorder)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("jwtAccsess", "valid_token")
	c.Set("user_id", "user123")

	s.tokenMock.On("ExtractUserID", "valid_token").Return("user123", nil)
	s.dbMock.On("VerifyUserExists", "user123").Return(false)

	handler := AutorizationMiddleWare()
	handler(c)

	s.True(c.IsAborted())
	s.Equal(http.StatusUnauthorized, s.testRecorder.Code)
	s.tokenMock.AssertExpectations(s.T())
	s.dbMock.AssertExpectations(s.T())
}

func (s *MiddlewareTestSuite) TestAuthMiddleware_ExpiredAccessToken() {
	c, _ := gin.CreateTestContext(s.testRecorder)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("jwtAccsess", "expired_token")
	c.Set("user_id", "user123")

	s.tokenMock.On("ExtractUserID", "expired_token").Return("", jwt.ErrTokenExpired)

	handler := AutorizationMiddleWare()
	handler(c)

	s.True(c.IsAborted())
	s.Equal(http.StatusUnauthorized, s.testRecorder.Code)
	s.tokenMock.AssertExpectations(s.T())
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
