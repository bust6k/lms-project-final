package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
)

func GetOrPostTask(c *gin.Context) {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()
	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
		return
	}

	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodPost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "метод не разрешен"})
		return
	}

	if c.Request.Method == http.MethodPost {
		err := c.ShouldBindJSON(&variables.CurrentTask)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "ошибка при чтении тела запроса"})
			logger.Warn("ошибка при чтении тела запроса")
			return

		}
		c.Status(http.StatusOK)

	} else if c.Request.Method == http.MethodGet {

		c.JSON(http.StatusOK, variables.CurrentTask)

		variables.CurrentTask = entites.Task{0, 0, 0, "", variables.CurrentTask.Operation_time}
		return

	}
}
