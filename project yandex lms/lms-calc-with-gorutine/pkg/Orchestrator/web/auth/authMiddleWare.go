package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"project_yandex_lms/lms-calc-with-gorutine/config"
	"project_yandex_lms/lms-calc-with-gorutine/database"
)

func AutorizationMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		var tokenAccsess string

		_, okInvalid := c.Get("accsessInvalid")

		if okInvalid {

			tokenRefresh, err := c.Cookie("Refresh")
			if err == http.ErrNoCookie {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "аутентифицируетесь  еще раз на /api/v1/login"})
				return
			}
			userId, err := PullUserIdByJwtToken(tokenRefresh)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("ошибка при получении user id:%v", err)})
				return
			}

			err = RefreshTokens(userId, c)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка при обновлении jwt токенов:%v", err)})
				return
			}
			c.JSON(http.StatusOK, gin.H{"check": "сделайте запрос еще раз чобы убедиться то что это вы"})
			return
		}

		userId, ok := c.Get("user_id")

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("ошибка при получении id пользователя:%s", userId)})
			return
		}

		RawAccsess := c.MustGet("jwtAccsess")

		tokenAccsess = fmt.Sprintf("%v", RawAccsess)

		parsedToken, err := jwt.Parse(tokenAccsess, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.DefaultJWTConfig().PasswordSigningJwt), nil
		})

		if err != nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("ошибка при парсинге jwt токена:%v", err)})

			return
		}

		if parsedToken.Valid {
			c.Set("user_id", userId)
			c.Next()

			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "ошибка авторизации: юзер не имеет  refresh токен или  он является не валдиным"})

		return

	}

}

func AutentificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenAccsess, err := c.Cookie("Access")

		if err != nil {
			c.Set("jwtAccsess", tokenAccsess)
			c.Set("accsessInvalid", true)
			c.Next()

			return
		}

		userId, err := PullUserIdByJwtToken(tokenAccsess)

		if err != nil {
			c.Set("jwtAccsess", tokenAccsess)
			c.Next()

			return
		}

		ok := database.CheckUserByUserId(fmt.Sprintf("%v", userId))
		if ok {
			c.Set("jwtAccsess", tokenAccsess)
			c.Set("user_id", userId)
			c.Next()

			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": " ошибка аутентификации:"})

		return
	}
}
