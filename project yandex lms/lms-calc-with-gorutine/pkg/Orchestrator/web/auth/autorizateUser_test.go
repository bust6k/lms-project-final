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

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) pullUserIdByLoginAndPassword(login, password string) (string, error) {
	args := m.Called(login, password)
	return args.String(0), args.Error(1)
}

func (m *mockDatabase) CheckUserInDB(user models.User) (error, bool) {
	args := m.Called(user)
	return args.Error(0), args.Bool(1)
}

type mockJWT struct {
	mock.Mock
}

func (m *mockJWT) CreateNewSignedJwtTokens(userId string) (string, string, error) {
	args := m.Called(userId)
	return args.String(0), args.String(1), args.Error(2)
}

type authTestSuite struct {
	suite.Suite
	dbMock  *mockDatabase
	jwtMock *mockJWT
	router  *gin.Engine
}

func (s *authTestSuite) SetupTest() {
	s.dbMock = new(mockDatabase)
	s.jwtMock = new(mockJWT)
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.Default()

	s.router.POST("/auth", AutorizateUser)
}

func (s *authTestSuite) TestSuccessfulAuth() {
	userId := uuid.NewString()
	body := bytes.NewBufferString(`{"login":"test","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/auth", body)
	w := httptest.NewRecorder()

	s.dbMock.On("PullUserIdByLoginAndPassword", "test", "pass").Return(userId, nil)
	expectedUser := models.NewUser("test", "pass", userId)
	s.dbMock.On("CheckUserInDB", expectedUser).Return(nil, true)
	s.jwtMock.On("CreateNewSignedJwtTokens", userId).Return("access", "refresh", nil)

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	cookies := w.Result().Cookies()
	s.Len(cookies, 2)
	s.Equal("Access", cookies[0].Name)
	s.Equal("access", cookies[0].Value)
	s.Equal("Refresh", cookies[1].Name)
	s.Equal("refresh", cookies[1].Value)
	s.dbMock.AssertExpectations(s.T())
	s.jwtMock.AssertExpectations(s.T())
}

func (s *authTestSuite) TestInvalidJSON() {
	body := bytes.NewBufferString(`invalid json`)
	req, _ := http.NewRequest("POST", "/auth", body)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
}

func (s *AuthTestSuite) TestUserNotFoundInDB() {
	body := bytes.NewBufferString(`{"login":"test","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/auth", body)
	w := httptest.NewRecorder()

	s.dbMock.On("PullUserIdByLoginAndPassword", "test", "pass").Return("", errors.New("user not found"))

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.dbMock.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestUserNotAuthenticated() {
	userId := uuid.NewString()
	body := bytes.NewBufferString(`{"login":"test","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/auth", body)
	w := httptest.NewRecorder()

	s.dbMock.On("PullUserIdByLoginAndPassword", "test", "pass").Return(userId, nil)
	expectedUser := models.NewUser("test", "pass", userId)
	s.dbMock.On("CheckUserInDB", expectedUser).Return(nil, false)

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.dbMock.AssertExpectations(s.T())
}

func (s *authTestSuite) TestJWTCreationFailed() {
	userId := uuid.NewString()
	body := bytes.NewBufferString(`{"login":"test","password":"pass"}`)
	req, _ := http.NewRequest("POST", "/auth", body)
	w := httptest.NewRecorder()

	s.dbMock.On("PullUserIdByLoginAndPassword", "test", "pass").Return(userId, nil)
	expectedUser := models.NewUser("test", "pass", userId)
	s.dbMock.On("CheckUserInDB", expectedUser).Return(nil, true)
	s.jwtMock.On("CreateNewSignedJwtTokens", userId).Return("", "", errors.New("jwt error"))

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.dbMock.AssertExpectations(s.T())
	s.jwtMock.AssertExpectations(s.T())
}

func TestAuthSuiteStart(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
