package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/database"

	"strconv"
)

func GetASeparateExpression(c *gin.Context) {

	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	userId := c.MustGet("user_id")

	value := c.Param("id")

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "ошибка при преобразовании id в целочисленный тип ")
		logger.Warn("ошибка при преобразовании id в целочисленный тип", zap.String("id", value))
		return
	}

	expr, err := database.GetSeparateProcessedExprInDB(valueInt, fmt.Sprintf("%v", userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("ошибка при получении выражения:%v", err)})
		logger.Warn("ошибка при получении выражения")
		return
	}
	c.JSON(http.StatusOK, &expr)
}
