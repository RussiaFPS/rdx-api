// Package main Основной пакет для работы с Api.
package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//Данные для подключения в БД
var connStr = "user=postgres password=1234 dbname=rdx sslmode=disable"

//Основаня функция
func main() {
	db, err := sql.Open("postgres", connStr) //Подключение к Бд
	if err != nil {
		log.Panic(err)
	}

	app := gin.Default()
	app.POST("/post/:username/:password", func(context *gin.Context) { //Обработка post запроса
		username := context.Param("username") //Получение параметра username
		password := context.Param("password") //Получение параметра password

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //Хэширование пароля
		if err != nil {
			log.Panic(err)
		}

		result, err := db.Exec("insert into \"User\" (username, password, created_at)  values ($1, $2, $3)",
			username, string(hashedPassword), time.Now()) //Запрос в Бд
		if err != nil {
			log.Panic(err)
		}

		res, _ := result.RowsAffected() //Проверка на успех в запросе
		if res > 0 {
			context.String(http.StatusOK, "Successful")
		} else {
			context.String(http.StatusBadRequest, "Failed")
		}
	})

	update := app.Group("/update") //Создание группы update запросов
	{
		update.PATCH("/name/:old-username/:new-username", func(context *gin.Context) { //Обработка изменения username
			oldUsername := context.Param("old-username") //Получение парметра старого username
			newUsername := context.Param("new-username") //Получение параметра нового username

			result, err := db.Exec("UPDATE \"User\" SET username = $1,updated_at = $2 WHERE username=$3", newUsername, time.Now(), oldUsername) //Запрос в БД
			if err != nil {
				log.Panic(err)
			}

			res, _ := result.RowsAffected() //Проверка на успех запроса
			if res > 0 {
				context.String(http.StatusOK, "Successful")
			} else {
				context.String(http.StatusBadRequest, "Failed")
			}
		})
		update.PATCH("/password/:username/:new-password", func(context *gin.Context) { //Обработка запроса на изменение пароля
			username := context.Param("username")                                                        //Получение параметра username
			newPassword := context.Param("new-password")                                                 //Получение параметра новый пароль
			newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost) //Хэширование нового проля
			if err != nil {
				log.Panic(err)
			}
			result, err := db.Exec("UPDATE \"User\" SET password = $1,updated_at = $2 WHERE username = $3", newPasswordHash, time.Now(), username) //Запрос в Бд
			if err != nil {
				log.Panic(err)
			}
			res, _ := result.RowsAffected() //Проверка на успех в запросе
			if res > 0 {
				context.String(http.StatusOK, "Successful")
			} else {
				context.String(http.StatusBadRequest, "Failed")
			}
		})
	}

	app.DELETE("/delete/:id", func(context *gin.Context) { //Обработка удаления данных
		id := context.Param("id")                                        //Получение параметра id
		result, err := db.Exec("DELETE FROM \"User\" WHERE id = $1", id) //Запрос в БД
		if err != nil {
			context.String(http.StatusBadRequest, "Failed")
			return
		}
		res, _ := result.RowsAffected() //Проверка на успех запроса
		if res > 0 {
			context.String(http.StatusOK, "Successful")
		} else {
			context.String(http.StatusBadRequest, "Failed")
		}
	})

	app.GET("/get/:id", func(context *gin.Context) { //Обработка запроса на получение Get
		id := context.Param("id")                                      //Получение параметра id
		row := db.QueryRow("select * from \"User\" where id = $1", id) //Запрос в Бд
		if err != nil {
			log.Panic(err)
		}

		u := User{}
		err := row.Scan(&u.id, &u.username, &u.password, &u.createAt, &u.updateAt) //Запись данных в структуру
		if err != nil {                                                            //Проверка успеха запроса
			context.String(http.StatusBadRequest, "Failed")
		} else {
			context.JSON(http.StatusOK, gin.H{"id": u.id, "username": u.username, "password": u.password, "createAt": u.createAt, "updateAt": u.updateAt})
		}
	})

	defer db.Close()          //Отложенное закрытие подключения к БД
	app.Run("localhost:8080") //Запуск сервиса на порту 8080
}
