package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/database"
	"project_yandex_lms/lms-calc-with-gorutine/models"
)

func RegistrateUser(c *gin.Context) {
	logger, err := zap.NewDevelopment()

	defer logger.Sync()

	if err != nil {
		log.Printf("ошибка при создании логгера zap с ошибкой: %v", err)

		return
	}

	src, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("ошибка при попытке прочитать тело ответа с ошибкой:%v", err)})
		logger.Warn("при попытке порчитать данные юзера произошла ошибка", zap.Error(err))

		return
	}

	login, pass, err := models.UnmarshalRegistrationDetailsFromJSON(src)
	userId := uuid.NewString()

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("ошибка при попытке получить имя и пароль с ошибкой:%v", err)})
		logger.Warn("юзер не смог зарегистрироваться так как произошла ошибка при декодировании его данных", zap.String("user login", login), zap.String("user password", pass), zap.Error(err))

		return

	}

	user := models.NewUser(login, pass, userId)
	err = database.SaveUserInDB(user)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка  при сохранении пользователя в нашу базу данных:%v", err)})
		logger.Warn("при попытке сохранить юзера в бд произошла ошибка", zap.String("user login", login), zap.String("user password", pass))

		return
	}

	accsessToken, refreshToken, err := CreateNewSignedJwtTokens(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка  при создании jwt токенов:%v", err)})
		logger.Warn("при попытке создать jwt токены для юзера произошла ошибка", zap.String("accsess token", accsessToken), zap.String("refresh token", refreshToken))

		return
	}

	c.JSON(http.StatusCreated, gin.H{"succsess": "регистрация прошла успешно"})

	return
}
