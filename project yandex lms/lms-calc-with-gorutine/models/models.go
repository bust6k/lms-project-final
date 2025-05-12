package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type User struct {
	Id       int    `json:"id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
	User_id  string `json:"user_Id"`
}

func NewUser(login string, password string, userId string) *User {
	return &User{
		Id:       0,
		Login:    login,
		Password: password,
		User_id:  userId,
	}
}

func UnmarshalRegistrationDetailsFromJSON(src []byte) (login, password string, err error) {

	var data map[string]json.RawMessage
	if err := json.Unmarshal(src, &data); err != nil {
		return "", "", fmt.Errorf("invalid JSON format: %w", err)
	}

	
	if len(data) != 1 {
		return "", "", errors.New("JSON must contain exactly one key-value pair")
	}

	
	for key, val := range data {
		login = key

		// Парсим пароль из RawMessage
		var pass string
		if err := json.Unmarshal(val, &pass); err != nil {
			return "", "", fmt.Errorf("invalid password format: %w", err)
		}
		password = pass
	}

	return login, password, nil
}

type ProcessedExpression struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
	UserId string  `json:"user_id"`
}

func NewProcessedExpression(status string, result float64, userId string) *ProcessedExpression {

	return &ProcessedExpression{
		Status: status,
		Result: result,
		UserId: userId,
	}
}
