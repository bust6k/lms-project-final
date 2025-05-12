package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"project_yandex_lms/lms-calc-with-gorutine/config"
	"project_yandex_lms/lms-calc-with-gorutine/models"
	"project_yandex_lms/lms-calc-with-gorutine/pkg/Orchestrator/web/auth/security"
)

const dsn = "root:pass12345@tcp(127.0.0.1:3306)/LMS"

func getDB() (*sqlx.DB, error) {
	return sqlx.Open("mysql", dsn)
}

func SaveUserInDB(user *models.User) error {
	db, err := getDB()
	if err != nil {
		return fmt.Errorf("ошибка при открытии соединения с бд:%v", err)
	}
	defer db.Close()

	hashed, err := security.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("ошибка при генерации хэша для пароля:%v", err)
	}

	_, err = db.Exec("INSERT INTO Users (Login,Password,User_id) VALUES (?,?,?)", user.Login, hashed, user.User_id)
	return err
}

func SaveProcessedExprToDB(expr models.ProcessedExpression) error {
	db, err := getDB()
	if err != nil {
		return fmt.Errorf("ошибка при открытии соединения с бд:%v", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO ProcessedExpressions (Status,Result,user_id) VALUES (?,?,?)", expr.Status, expr.Result, expr.UserId)
	return err
}

func GetProcessedExprsInDB() ([]models.ProcessedExpression, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var expressions []models.ProcessedExpression
	err = db.SelectContext(ctx, &expressions, "SELECT id, status, result, user_id as userid FROM ProcessedExpressions")
	return expressions, err
}

func GetSeparateProcessedExprInDB(exprId int, userId string) (models.ProcessedExpression, error) {
	db, err := getDB()
	if err != nil {
		return models.ProcessedExpression{}, fmt.Errorf("ошибка при открытии соединения с бд:%v", err)
	}
	defer db.Close()

	var expr models.ProcessedExpression
	err = db.Get(&expr, "SELECT id, status, result, user_id as userid FROM ProcessedExpressions WHERE Id=? AND user_id=?", exprId, userId)
	return expr, err
}

func GetCurrentIdInDB() (int, error) {
	db, err := getDB()
	if err != nil {
		return config.DefaultDBConfig().CnnectionRefused, fmt.Errorf("ошибка при открытии соединения с бд:%v", err)
	}
	defer db.Close()

	var lastId int
	err = db.Get(&lastId, "SELECT Id FROM ProcessedExpressions ORDER BY id DESC LIMIT 1")
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return lastId, err
}

func CheckUserByUserId(userId string) bool {
	db, err := getDB()
	if err != nil {
		return false
	}
	defer db.Close()

	var exists bool
	err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM Users WHERE user_id=?)", userId)
	return err == nil && exists
}

func PullUserIdAndPasswordByLogin(login string) (string, string, error) {
	if login == "" {
		return "", "", fmt.Errorf("логин не может быть пустым")
	}

	db, err := getDB()
	if err != nil {
		return "", "", fmt.Errorf("ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	var result struct {
		UserID   string `db:"User_id"`
		Password string `db:"Password"`
	}

	err = db.Get(&result, "SELECT User_id, Password FROM Users WHERE Login = ? LIMIT 1", login)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("пользователь с логином '%s' не найден", login)
		}
		return "", "", fmt.Errorf("ошибка запроса: %v", err)
	}

	return result.UserID, result.Password, nil
}
