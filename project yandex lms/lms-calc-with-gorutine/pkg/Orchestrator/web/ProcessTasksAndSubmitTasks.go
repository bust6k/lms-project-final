package web

import (
	"bytes"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"project_yandex_lms/lms-calc-with-gorutine/variables"
)

func ProcessTasksAndSubmitTasks(c *gin.Context) {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()
	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	var slicceOfTasks []entites.Task

	err = c.ShouldBindJSON(&slicceOfTasks)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "ошибка при чтении запроса"})
		logger.Warn("ошибка при чтении запроса", zap.Error(err))
		return
	}
	if len(slicceOfTasks) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "слайс задач сейчас пуст"})
		logger.Debug("слайс задач сейчас пуст")
		return
	}

	variables.TheTasks = slicceOfTasks

	for len(variables.TheTasks) > 0 {
		firstElement := variables.TheTasks[0]
		variables.CurrentTask = firstElement

		variables.TheTasks = variables.TheTasks[1:]

		var bytesjson bytes.Buffer
		bytesElement, err := sonic.Marshal(firstElement)

		bytesjson.Write(bytesElement)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "ошибка при сериализации ответа"})
			logger.Warn("ошибка при сериализации ответа")

			return
		}

		bytestopost := bytes.NewReader(bytesjson.Bytes())

		_, err = http.Post("http://localhost:8080/internal/task", "application/json", bytestopost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке задачи на сервер"})
			logger.Warn("Ошибка при отправке задачи на сервер", zap.Error(err))
			return
		}
	}

	c.Status(http.StatusOK)
	return

}
