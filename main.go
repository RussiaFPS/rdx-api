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

var connStr = "user=postgres password=1234 dbname=rdx sslmode=disable"

func main() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	app := gin.Default()
	app.POST("/post/:username/:password", func(context *gin.Context) {
		username := context.Param("username")
		password := context.Param("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Panic(err)
		}

		result, err := db.Exec("insert into \"User\" (username, password, created_at)  values ($1, $2, $3)",
			username, string(hashedPassword), time.Now())
		if err != nil {
			log.Panic(err)
		}

		res, _ := result.RowsAffected()
		if res > 0 {
			context.String(http.StatusOK, "Successful")
		} else {
			context.String(http.StatusBadRequest, "Failed")
		}
	})

	update := app.Group("/update")
	{
		update.PATCH("/name/:old-username/:new-username", func(context *gin.Context) {
			oldUsername := context.Param("old-username")
			newUsername := context.Param("new-username")

			result, err := db.Exec("UPDATE \"User\" SET username = $1,updated_at = $2 WHERE username=$3", newUsername, time.Now(), oldUsername)
			if err != nil {
				log.Panic(err)
			}

			res, _ := result.RowsAffected()
			if res > 0 {
				context.String(http.StatusOK, "Successful")
			} else {
				context.String(http.StatusBadRequest, "Failed")
			}
		})
		update.PATCH("/password/:username/:new-password", func(context *gin.Context) {
			username := context.Param("username")
			newPassword := context.Param("new-password")
			newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				log.Panic(err)
			}
			result, err := db.Exec("UPDATE \"User\" SET password = $1,updated_at = $2 WHERE username = $3", newPasswordHash, time.Now(), username)
			if err != nil {
				log.Panic(err)
			}
			res, _ := result.RowsAffected()
			if res > 0 {
				context.String(http.StatusOK, "Successful")
			} else {
				context.String(http.StatusBadRequest, "Failed")
			}
		})
	}

	app.DELETE("/delete/:id", func(context *gin.Context) {
		id := context.Param("id")
		result, err := db.Exec("DELETE FROM \"User\" WHERE id = $1", id)
		if err != nil {
			log.Panic(err)
		}
		res, _ := result.RowsAffected()
		if res > 0 {
			context.String(http.StatusOK, "Successful")
		} else {
			context.String(http.StatusBadRequest, "Failed")
		}
	})

	app.GET("/get/:id", func(context *gin.Context) {
		id := context.Param("id")
		row := db.QueryRow("select * from \"User\" where id = $1", id)
		if err != nil {
			log.Panic(err)
		}

		u := User{}
		err := row.Scan(&u.id, &u.username, &u.password, &u.createAt, &u.updateAt)
		if err != nil {
			context.String(http.StatusBadRequest, "Failed")
		} else {
			context.JSON(http.StatusOK, gin.H{"id": u.id, "username": u.username, "password": u.password, "createAt": u.createAt, "updateAt": u.updateAt})
		}
	})

	defer db.Close()
	app.Run("localhost:8080")
}
