package main

import (
	"sync"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	users    []User
	usersMtx sync.Mutex
	apiKey   string
)
