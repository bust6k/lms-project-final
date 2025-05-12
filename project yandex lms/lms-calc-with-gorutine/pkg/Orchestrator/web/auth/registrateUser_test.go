package auth

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"project_yandex_lms/lms-calc-with-gorutine/models"
	"testing"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) pullUserIdByLoginAndPassword(login, password string) (string, error) {
	args := m.Called(login, password)
	return args.String(0), args.Error(1)
}

func (m *MockDatabase) CheckUserInDB(user models.User) (error, bool) {
	args := m.Called(user)
	return args.Error(0), args.Bool(1)
}

type mockJWTGenerator struct {
	mock.Mock
}

func (m *mockJWTGenerator) CreateNewSignedJwtTokens(userId string) (string, string, error) {
	args := m.Called(userId)
	return args.String(0), args.String(1), args.Error(2)
}

type AuthTestSuite struct {
	suite.Suite
	dbMock  *MockDatabase
	jwtMock *mockJWTGenerator
	router  *gin.Engine
}

func (s *AuthTestSuite) SetupTest() {
	s.dbMock = new(MockDatabase)
	s.jwtMock = new(mockJWTGenerator)
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.Default()

	s.router.POST("/register", AutorizateUser)
}

func (s *AuthTestSuite) TestSuccessfulAuth() {
	userId := uuid.NewString()
	body := bytes.NewBufferString(`{"login":"test","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/register", body)
	w := httptest.NewRecorder()

	s.dbMock.On("PullUserIdByLoginAndPassword", "test", "pass").Return(userId, nil)
	s.dbMock.On("CheckUserInDB", mock.Anything).Return(nil, true)
	s.jwtMock.On("CreateNewSignedJwtTokens", userId).Return("access", "refresh", nil)

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.dbMock.AssertExpectations(s.T())
	s.jwtMock.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestInvalidJSON() {
	body := bytes.NewBufferString(`invalid json`)
	req, _ := http.NewRequest("POST", "/register", body)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *AuthTestSuite) TestUserNotFound() {
	body := bytes.NewBufferString(`{"login":"test","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/register", body)
	w := httptest.NewRecorder()

	s.dbMock.On("PullUserIdByLoginAndPassword", "test", "pass").Return("", errors.New("not found"))

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.dbMock.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestJWTCreationFailed() {
	userId := uuid.NewString()
	body := bytes.NewBufferString(`{"login":"test","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/register", body)
	w := httptest.NewRecorder()

	s.dbMock.On("PullUserIdByLoginAndPassword", "test", "pass").Return(userId, nil)
	s.dbMock.On("CheckUserInDB", mock.Anything).Return(nil, true)
	s.jwtMock.On("CreateNewSignedJwtTokens", userId).Return("", "", errors.New("jwt error"))

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.dbMock.AssertExpectations(s.T())
	s.jwtMock.AssertExpectations(s.T())
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
