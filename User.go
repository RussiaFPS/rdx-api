package main

import (
	"time"
)

type User struct {
	id       int
	username string
	password string
	createAt time.Time
	updateAt interface{}
}
