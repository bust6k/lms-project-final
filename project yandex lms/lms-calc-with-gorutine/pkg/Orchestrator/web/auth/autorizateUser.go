package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/database"
	"project_yandex_lms/lms-calc-with-gorutine/models"
	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/web/auth/security"
)

func AutorizateUser(c *gin.Context) {
	RegistrationDetailsByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("ошибка при чтении тела запроса:%v", err)})
		return
	}

	log, pass, err := models.UnmarshalRegistrationDetailsFromJSON(RegistrationDetailsByte)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("ошибка формата данных:%v", err)})
		return
	}

	userId, hashedPassword, err := database.PullUserIdAndPasswordByLogin(log)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный логин или пароль"})
		return
	}

	err = security.Compare(pass, []byte(hashedPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный логин или пароль"})
		return
	}

	signedAccess, signedRefresh, err := CreateNewSignedJwtTokens(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка генерации токенов:%v", err)})
		return
	}

	c.SetCookie("Access", signedAccess, 15*60, "/", "", true, true)
	c.SetCookie("Refresh", signedRefresh, 7*24*60*60, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"success": "аутентификация прошла успешно"})
}
