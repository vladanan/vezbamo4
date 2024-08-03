// Package routes služi da obrađuje zahvete iz main
package routes

import (
	"os"

	"github.com/gorilla/sessions"

	"github.com/joho/godotenv"
)

// var store string = ""

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

var godotevnErr = godotenv.Load(".env")

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	// key   = []byte("super-secret-key")
	key   = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)
)
