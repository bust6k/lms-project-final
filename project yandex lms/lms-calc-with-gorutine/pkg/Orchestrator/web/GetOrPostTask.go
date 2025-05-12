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
	
	switch c.Request.Method {
	case http.MethodPost:
		handlePostRequest(c,logger)
	case http.MethodGet:
		handleGetRequest(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "метод не разрешен"})
	}
}


func handlePostRequest(c *gin.Context,logger *zap.Logger){
	
	err := c.ShouldBindJSON(&variables.CurrentTask)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "ошибка при чтении тела запроса"})
		logger.Warn("ошибка при чтении тела запроса",zap.Error(err))
		return

	}
	c.Status(http.StatusOK)
}


func handleGetRequest(c *gin.Context){
	c.JSON(http.StatusOK, variables.CurrentTask)

	variables.CurrentTask = entites.Task{0, 0, 0, "", variables.CurrentTask.Operation_time}
	return
}
