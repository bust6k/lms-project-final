package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/config"
	"project_yandex_lms/lms-calc-with-gorutine/database"
	"project_yandex_lms/lms-calc-with-gorutine/entites"
	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/algoritms/AST"
	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/algoritms/RPN"

	"project_yandex_lms/lms-calc-with-gorutine/pkg/agent"

	"strconv"
	"time"
)

var cfg = config.DefaultEnvConfig()

func CreateAUserExpressionInSystem(c *gin.Context) {

	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)
	}

	logger.Debug("------------Start Pipeline of project------------------")

	userId := c.MustGet("user_id")

	var newExpression entites.UnprocessedUserExpression
	err = c.ShouldBindJSON(&newExpression)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "тело запроса не было считано"})
		logger.Warn("тело запроса не было считано и не было распарщено в структуру", zap.String("expression", newExpression.Expression))
		return
	}

	currentIdInDB, err := database.GetCurrentIdInDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка при попытке получить текущий id в базе данных:%v", err)})
	}
	newExpression.Id = currentIdInDB
	logger.Info("в системе создано новое выражение",
		zap.String("expression id", strconv.Itoa(newExpression.Id)),
		zap.String("expression value", newExpression.Expression))

	resp := struct {
		Id int `json:"id"`
	}{Id: newExpression.Id}

	RpnExpression, err := RPN.InfixToRPN(newExpression.Expression)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при преобразовании выражения в rpn"})
		logger.Warn("ошибка при преобразовании выражения в rpn")
		return

	}
	logger.Debug("вот как выглядит выражение в rpn:", zap.Strings("rpns", RpnExpression))
	BuildedAst, err := AST.BuildAST(RpnExpression)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при преобразовании rpn выражения в ast"})
		logger.Warn("ошибка при преобразовании rpn выражения в ast")
		return
	}
	logger.Debug("вот как выглядит ast дерево", zap.Reflect("ast", BuildedAst))

	tasks := AST.SplitAST(BuildedAst)

	for i := 0; i < len(tasks); i++ {

		if tasks[i].Operation == "+" {

			tasks[i].Operation_time = time.Duration(cfg.TIME_ADDITION_MS) * time.Millisecond

		} else if tasks[i].Operation == "-" {

			tasks[i].Operation_time = time.Duration(cfg.TIME_SUBSTRACTION_MS) * time.Millisecond
		} else if tasks[i].Operation == "*" {

			tasks[i].Operation_time = time.Duration(cfg.TIME_MULTIPLICATIONS_MS) * time.Millisecond
		} else if tasks[i].Operation == "/" {

			tasks[i].Operation_time = time.Duration(cfg.TIME_DIVISIONS_MS) * time.Millisecond
		}
	}

	AST.PostTasksToServer(tasks)

	agent := agent.NewAgent(cfg.COMPUTING_POWER)
	agent.CalculateTasks()
	err = agent.PostTaskResultsToServer(fmt.Sprintf("%v", userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("произошла ошибка при попытке отправить результаты на сервер:%s", userId)})
	}

	c.JSON(http.StatusCreated, resp)

}
