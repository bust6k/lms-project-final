package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"project_yandex_lms/lms-calc-with-gorutine/config"
)

func PullUserIdByJwtToken(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("пустой токен")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
		}
		return []byte(config.DefaultJWTConfig().PasswordSigningJwt), nil
	})

	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге jwt токена: %v", err)
	}

	if !parsedToken.Valid {
		return "", fmt.Errorf("невалидный токен")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("неверный формат claims")
	}

	if claims["user_id"] == nil {

		return "", fmt.Errorf("токен не содержит user_id")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {

		return "", fmt.Errorf("user_id должен быть строкой ")
	}

	return userId, nil
}
