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

func mockProcessAndSubmitTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []entites.Task

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytesreq, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = sonic.Unmarshal(bytesreq, &tasks)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func TestProcessTasksAndSubmitTasks(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		contentType        string
		inputData          string
		exceptedStatusCode int
	}{
		{
			name:               "valid data",
			inputData:          `[{"id":0,"arg1":2,"arg2":2,"operation":"+","Operation_time":100}]`,
			method:             http.MethodPost,
			contentType:        "application/json",
			exceptedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid json syntax",
			inputData:          `[{"id":0,"arg1":2,"ar2":2,"operation":+,"Operation_time":100}]`,
			method:             http.MethodPost,
			exceptedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid method",
			method:             http.MethodGet,
			contentType:        "application/json",
			exceptedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid content type",
			method:             http.MethodPost,
			contentType:        "application/xml",
			exceptedStatusCode: http.StatusBadRequest,
			inputData: `
<Tasks>
    <Task>
        <id>0</id>
        <arg1>2</arg1>
        <ar2>2</ar2>
        <operation>+</operation>
        <Operation_time>100</Operation_time>
    </Task>
</Tasks>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/internal", strings.NewReader(tt.inputData))
			req.Header.Set("Content-Type", tt.contentType)

			handler := http.HandlerFunc(mockProcessAndSubmitTasks)
			handler.ServeHTTP(w, req)

			if w.Code != tt.exceptedStatusCode {
				t.Errorf("ожидался статус код:%d, получил:%d", tt.exceptedStatusCode, w.Code)
			}
		})
	}
}
