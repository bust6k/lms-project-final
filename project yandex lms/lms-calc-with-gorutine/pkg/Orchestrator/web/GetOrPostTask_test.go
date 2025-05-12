//go:build unit

package web

import (
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"net/http/httptest"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"strings"
	"testing"
)

var mockCurrentTask entites.Task

func mockGetOrPostTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if err := sonic.Unmarshal(body, &mockCurrentTask); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	byteTask, err := sonic.Marshal(mockCurrentTask)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if _, err := w.Write(byteTask); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func TestGetOrPostTask(t *testing.T) {
	tests := []struct {
		name               string
		inputData          string
		method             string
		contentType        string
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name:               "valid data",
			inputData:          `{"id":0,"arg1":2,"arg2":2,"operation":"+","operation_time":100}`,
			method:             http.MethodPost,
			contentType:        "application/json",
			expectedBody:       "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid data",
			method:             http.MethodPost,
			contentType:        "application/json",
			inputData:          `{"id":"0","arg1":"2",ar34":"8","operation":"-","Operation_tim/e""100"}`,
			expectedBody:       "",
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "not allowed method",
			inputData:          `{"id":0,"arg1":80,"arg2":90,"operation":"*","Operation_time":300}`,
			expectedBody:       "",
			method:             http.MethodPut,
			contentType:        "application/json",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid content type",
			inputData:          "id: 0\narg1: 80\narg2: 90\noperation: \"*\"\nOperation_time: 300",
			method:             http.MethodPost,
			contentType:        "application/yaml",
			expectedBody:       "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "get current task",
			method:             http.MethodGet,
			contentType:        "application/json",
			expectedBody:       `{"id":0,"arg1":2,"arg2":2,"operation":"+","operation_time":100}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(mockGetOrPostTask)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/internal/task", strings.NewReader(tt.inputData))
			req.Header.Set("Content-Type", tt.contentType)

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("ожидался статус код: %d, получил: %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("ожидалось тело ответа: %s, получил: %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}
