package application

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

//go: build e2e

const (
	registerUrl  = "http://localhost:8080/api/v1/register"
	loginUrl     = "http://localhost:8080/api/v1/login"
	calculateUrl = "http://localhost:8080/api/v1/calculate"
)

var accsessToken *http.Cookie
var refreshToken *http.Cookie

func TestMain(m *testing.M) {

	testApp := New()
	testApp.Setup()
	testApp.Run()

	code := m.Run()

	os.Exit(code)
}

func TestUser(t *testing.T) {

	t.Run("a test that check user request correct data and  register  correctly", func(t *testing.T) {

		req := bytes.NewReader([]byte(`{"testName":"testPass"}`))
		resp, err := http.Post(registerUrl, "application/json", req)

		assert.NoError(t, err)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("ошибка при регистрации юзера. ожидался статус код:%d, получил:%d", http.StatusCreated, resp.StatusCode)
		}

	})

	t.Run("a test that check user request incorrect data and  register  incorrectly", func(t *testing.T) {

		req := bytes.NewReader([]byte(`{"invalid name :"invalid password"}`))
		resp, err := http.Post(registerUrl, "application/json", req)

		assert.NoError(t, err)
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusCreated {
			t.Errorf("ошибка при регистрации юзера. ожидался статус код:%d, получил:%d", http.StatusUnprocessableEntity, resp.StatusCode)
		}
	})

	t.Run("a test that checks the user's login and that asks for the correct data and logs in correctly", func(t *testing.T) {

		req := bytes.NewReader([]byte(`{"testName":"testPass"}`))
		resp, err := http.Post(loginUrl, "application/json", req)

		assert.NoError(t, err)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("ошибка при логине юзера. ожидался статус код:%d, получил:%d", http.StatusOK, resp.StatusCode)
		}

		for _, cookie := range resp.Cookies() {
			switch cookie.Name {
			case "Accsess":
				accsessToken = cookie
			case "Refresh":
				refreshToken = cookie

			}
		}
	})

	t.Run("a test that checks the user's login and that asks for the incorrect data and logs in incorrectly", func(t *testing.T) {

		req := bytes.NewReader([]byte(`{"notExsistsLogin"":"notExsistsPassword"}`))
		resp, err := http.Post(loginUrl, "application/json", req)

		assert.NoError(t, err)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("ошибка при логине юзера. ожидался статус код:%d, получил:%d", http.StatusForbidden, resp.StatusCode)
		}
	})

}

func TestCreateExpression(t *testing.T) {
	t.Run("a test  that logined user send correct request", func(t *testing.T) {

		body := bytes.NewReader([]byte(`{"expression":"2+2*2"}`))

		req, err := http.NewRequest("POST", calculateUrl, body)
		req.Header.Set("Content-Type", "application-json")
		req.AddCookie(accsessToken)
		req.AddCookie(refreshToken)

		assert.NoError(t, err)

		client := &http.Client{}

		resp, err := client.Do(req)

		defer resp.Body.Close()

		assert.NoError(t, err)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("ошибка при логине юзера. ожидался статус код:%d, получил:%d", http.StatusCreated, resp.StatusCode)
		}

	})

	t.Run("a test  that unlogined user send correct request", func(t *testing.T) {

		body := bytes.NewReader([]byte(`{"expression":"2+2*2"}`))

		req, err := http.NewRequest("POST", calculateUrl, body)
		req.Header.Set("Content-Type", "application-json")

		assert.NoError(t, err)

		client := &http.Client{}

		resp, err := client.Do(req)

		defer resp.Body.Close()

		assert.NoError(t, err)

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("ошибка при логине юзера. ожидался статус код:%d, получил:%d", http.StatusUnauthorized, resp.StatusCode)
		}

	})

	t.Run("a test  that logined user send uncorrect request", func(t *testing.T) {

		body := bytes.NewReader([]byte(`{"testexpt":"test2+4"}`))

		req, err := http.NewRequest("POST", calculateUrl, body)
		req.Header.Set("Content-Type", "application-json")

		assert.NoError(t, err)

		client := &http.Client{}

		resp, err := client.Do(req)

		assert.NoError(t, err)

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("ошибка при логине юзера. ожидался статус код:%d, получил:%d", http.StatusUnauthorized, resp.StatusCode)
		}

	})

	t.Run("a test  that unlogined user send uncorrect request", func(t *testing.T) {

		body := bytes.NewReader([]byte(`{" ":"2+*"}`))

		req, err := http.NewRequest("POST", calculateUrl, body)
		req.Header.Set("Content-Type", "application-json")

		assert.NoError(t, err)

		client := &http.Client{}

		resp, err := client.Do(req)

		defer resp.Body.Close()

		assert.NoError(t, err)

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("ошибка при логине юзера. ожидался статус код:%d, получил:%d", http.StatusUnauthorized, resp.StatusCode)
		}

	})
}
