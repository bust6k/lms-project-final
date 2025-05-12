package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"project_yandex_lms/lms-calc-with-gorutine/config"

	"time"
)

func RefreshTokens(userId string, c *gin.Context) error {

	signedAcsess, signedRefresh, err := CreateNewSignedJwtTokens(userId)

	if err != nil {
		return fmt.Errorf("ошибка при созданиии jwt токенов:%v", err)
	}

	c.SetCookie("Access", signedAcsess, 15*60, "/", "", true, true)
	c.SetCookie("Refresh", signedRefresh, 7*24*60*60, "/", "", true, true)

	return nil
}

func CreateNewSignedJwtTokens(userId string) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
		"nbf":     time.Now().Unix(),
		"iat":     time.Now().Unix(),
		"user_id": userId,
		"type":    "acsess",
	})

	signedAcsess, err := token.SignedString([]byte(config.DefaultJWTConfig().PasswordSigningJwt))
	if err != nil {
		return "", "", fmt.Errorf("ошибка при подписи accsess токена:%v", err)
	}

	tokenref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
		"nbf":     time.Now().Unix(),
		"iat":     time.Now().Unix(),
		"user_id": userId,
		"type":    "refresh",
	})

	signedRefresh, err := tokenref.SignedString([]byte(config.DefaultJWTConfig().PasswordSigningJwt))
	if err != nil {
		return "", "", fmt.Errorf("ошибка при подписи refresh токена:%v", err)
	}

	return signedAcsess, signedRefresh, nil
}
