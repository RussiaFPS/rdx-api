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

	defer db.Close()
	app.Run("localhost:8080")
}
