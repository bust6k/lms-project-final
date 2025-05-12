//go:build unit

package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/web"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateAUserExpressionInSystem(t *testing.T) {

	tests := []struct {
		name         string
		inputData    interface{}
		expectedCode int
		setupMock    func()
	}{
		{
			name:         "user send correct expression",
			inputData:    map[string]string{"expression": "2+2*2"},
			expectedCode: http.StatusOK,
		},
		{
			name:         " user send invalid data",
			inputData:    map[string]string{"invalid": "data"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := gin.Default()
			gin.SetMode(gin.ReleaseMode)
			router.POST("/calculate", web.CreateAUserExpressionInSystem)

			jsonData, _ := json.Marshal(tt.inputData)
			body := bytes.NewBuffer(jsonData)

			req, _ := http.NewRequest("POST", "/calculate", body)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "expression_id")
			}
		})
	}
}
