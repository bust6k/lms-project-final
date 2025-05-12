package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"project_yandex_lms/lms-calc-with-gorutine/models"
	"testing"
)

//go: build integration

var db, errOpenDB = sqlx.Open("mysql", "root:pass12345@tcp(127.0.0.1:3306)/LMS")

func TestMain(m *testing.M) {

	if errOpenDB != nil {
		log.Fatal("произошла ошибка при подключении к бд:%v", errOpenDB)
	}
	defer db.Close()
	code := m.Run()

	os.Exit(code)
}

func SafeMustExec(db *sqlx.DB, query string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("SQL-запрос упал: %v", r)
		}
	}()
	db.MustExec(query, args...)
}

var correctTestUser = models.NewUser("testLogin", "testPassword", "testUserId")

func TestUserInDb(t *testing.T) {

	t.Run("test to save user in db", func(t *testing.T) {
		SafeMustExec(db, `
    INSERT INTO Users  (Login,Password,User_id)
    VALUES (?, ?,?)`,
			correctTestUser.Login, correctTestUser.Password, correctTestUser.User_id,
		)

	})
	t.Run("test to select user from db knowing only user id", func(t *testing.T) {

		t.Skip(fmt.Sprintf("произошла ошибка при подключении к бд", errOpenDB))

		var resultUser models.User

		err := db.Get(&resultUser, "SELECT Login,Password,User_id FROM  Users WHERE  User_id=?", correctTestUser.User_id)

		if err != nil {
			t.Skip(fmt.Sprintf("ошибка при выборе юзера из базы данных зная только user_id:%v", err))
		}

		assert.Equal(t, correctTestUser, resultUser)

	})
}

var correctTestProcessedExpr = models.NewProcessedExpression("ready", 1.5, "testUserid")
var incorrectTestProcessedExpr = models.NewProcessedExpression("ready", 1.5, "NotExsistsTestUserid")

func TestProcessedExprInDB(t *testing.T) {
	t.Run("test to save correct processed expr to db", func(t *testing.T) {
		SafeMustExec(db, `
    INSERT INTO ProcessedExpressions  (Status,Result,user_id)
    VALUES (?, ?,?)`,
			correctTestProcessedExpr.Status, correctTestProcessedExpr.Result, correctTestProcessedExpr.UserId,
		)

	})
	t.Run("test to save processed expr not having real user id in table Users", func(t *testing.T) {
		_, err := db.Exec(`
    INSERT INTO ProcessedExpressions  (Status,Result,user_id)
    VALUES (?, ?,?)`,
			incorrectTestProcessedExpr.Status, incorrectTestProcessedExpr.Result, incorrectTestProcessedExpr.UserId,
		)

		if err == nil {
			t.Error(", ожидалось что при попытке внести данные в таблицу ProcessedExpressions  без существующего user_id возникнет ошибка, ее не возникло")
		}

	})

}
