package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

type JWT interface {
	//psrse имитирует метод jwt.Parse только не парсит jwt токен а просто ищет подстроку passwordSigning в токене
	Parse(token string, passwordSigning string) error
	//думаю уже и так понятно
	IsAutorizate(token string, passwordSigning string) bool
}

type MockJWT struct {
}

func (m *MockJWT) Parse(token string, passwordSigning string) error {

	if bytes.HasSuffix([]byte(token), []byte(passwordSigning)) {
		return nil
	}
	return fmt.Errorf("token has no valid")
}

func (m *MockJWT) IsAutorizate(token string, passwordSigning string) bool {
	err := m.Parse(token, passwordSigning)
	if err != nil {
		return false
	}
	return true
}

type DB interface {
	//имитирует реальное поведение базы данных  разве что вместо url тут []string и нужное число должно быть на 1-вой позиции
	Query(db []string, id int) (string, int)
}

type MockDB struct {
}

func (m *MockDB) Query(db []string, id int) (string, int) {
	realId, err := strconv.Atoi(db[1])
	//допустим у нас есть только 4 выражения в бд
	if realId != id || id > 4 || id < 0 || err != nil {
		return "", http.StatusNotFound
	}

	return `{"id": "0", "Status": "ready","Result":"2"}`, http.StatusOK
}

func TestGetASeparateExpression(t *testing.T) {

	testsdb := []struct {
		name               string
		mock               *MockDB
		idExpr             int
		db                 []string
		exceptedStatusCode int
		exceptedResp       string
	}{
		{
			name:               "id exsist",
			mock:               &MockDB{},
			db:                 []string{"1", "3", "4"},
			idExpr:             3,
			exceptedStatusCode: http.StatusOK,
			exceptedResp:       `{"id": "0", "Status": "ready","Result":"2"`,
		},
		{
			name:               "id not exsist",
			mock:               &MockDB{},
			db:                 []string{"1", "78", "5"},
			idExpr:             6,
			exceptedStatusCode: http.StatusNotFound,
			exceptedResp:       "",
		},
		{
			name:               "id  is negatuive",
			mock:               &MockDB{},
			db:                 []string{"1", "-3", "4"},
			idExpr:             -3,
			exceptedStatusCode: http.StatusNotFound,
			exceptedResp:       "",
		},
	}

	for _, tt := range testsdb {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			str, code := tt.mock.Query(tt.db, tt.idExpr)
			if tt.exceptedStatusCode != code || tt.exceptedResp != str {
				t.Errorf("ошибка, ожидался статус  код:%d и ответ:%s но  получил:%d,%s", tt.exceptedStatusCode, tt.exceptedResp, code, str)
			}
		})
	}

	testsJWT := []struct {
		name                 string
		mock                 *MockJWT
		token                string
		signingPassword      string
		exceptedIsAutorizate bool
	}{
		{
			name:                 "valid token",
			mock:                 &MockJWT{},
			token:                "passSomething",
			signingPassword:      "pass",
			exceptedIsAutorizate: true,
		},
		{
			name:                 "invalid token",
			mock:                 &MockJWT{},
			token:                "passSomething",
			signingPassword:      "pass12345",
			exceptedIsAutorizate: false,
		},
	}

	for _, tt := range testsJWT {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.mock.IsAutorizate(tt.token, tt.signingPassword)
			if ok != tt.exceptedIsAutorizate {
				t.Errorf("ожидалось что юзер авторизован: %v, а в итоге:%v", tt.exceptedIsAutorizate, ok)
			}
		})
	}

}
