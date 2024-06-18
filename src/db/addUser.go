package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	// "encoding/json"
	// "github.com/vladanan/vezbamo4/src/models"
)

// type User struct {
// 	Hash_lozinka string `db:"hash_lozinka"`
// 	Email        string `db:"email"`
// 	Test         string `db:"test"`
// }

// func to_struct(user []byte) []models.User {
// 	var p []models.User
// 	err := json.Unmarshal(user, &p)
// 	if err != nil {
// 		fmt.Printf("Json error: %v", err)
// 	}
// 	return p
// }

func AddUser(email, user_name, password_str string) bool {
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	password := []byte(password_str)

	ciphertext, err := bcrypt.GenerateFromPassword(password, 5) //df
	// _, err := bcrypt.GenerateFromPassword(password, 5) //df
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from bcrypt encryption: %s\n", err)
		return false
	}

	fmt.Println("Ciphertexts: ", string(ciphertext))

	err2 := godotenv.Load(".env")
	if err2 != nil {
		fmt.Printf("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// U_id           int       `db:"u_id"`
	// Created_at     time.Time `db:"created_at"`
	// Hash_lozinka   string    `db:"hash_lozinka"`
	// Email          string    `db:"email"`
	// User_name      string    `db:"user_name"`
	// Mode           string    `db:"user_mode"`
	// Level          string    `db:"user_level"`
	// Basic          bool      `db:"basic"`
	// Js             bool      `db:"js"`
	// C              bool      `db:"c"`
	// Payment_date   time.Time `db:"payment_date"`
	// Payment_expire time.Time `db:"payment_expire"`
	// Verified_email bool      `db:"verified_email"`

	commandTag, err := conn.Exec(context.Background(), "INSERT INTO mi_users (hash_lozinka, email, user_name, user_mode, user_level, basic, js, c, verified_email, payment_date, payment_expire) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);", ciphertext, email, user_name, "user", 0, true, true, false, false, time.DateTime, time.DateTime)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return false
	}
	fmt.Printf("insert result: %v\n", commandTag)

	// err = bcrypt.CompareHashAndPassword([]byte("$2a$05$HYej4fyvWYnC5LvrSGEOD.bztJzcYn45t62etOTrN8d59BkoD7fhy"), password)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error from bcrypt Dencryption: %s\n", err)
	// 	return
	// } else {
	// 	fmt.Printf("\nProšlo je!")
	// }

	// PROVERA DA LI NEMA VEĆ USER-A SA ISTIM MEJLOM I USER_NAME

	// rows, err2 := conn.Query(context.Background(), "SELECT * FROM mi_users where email=$1;", email)
	// if err2 != nil {
	// 	fmt.Printf("Unable to make query: %v\n", err2)
	// 	return false, models.User{}
	// }

	// pgx_user, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	// if err != nil {
	// 	fmt.Printf("CollectRows error: %v", err)
	// 	//return
	// 	return false, models.User{}
	// }

	// // fmt.Print("pgx user::", pgx_user)

	// bytearray_user, err2 := json.Marshal(pgx_user)
	// if err2 != nil {
	// 	fmt.Printf("Json error: %v", err2)
	// }

	// var struct_user models.User

	// if string(bytearray_user) != "null" {
	// 	struct_user = to_struct(bytearray_user)[0]
	// } else {
	// 	return false, models.User{}
	// }

	// fmt.Print(struct_user)

	// AKO JE SVE OKEJ ŠALJE SE MEJL PREKO MEJL SERVERA

	// if already_authenticated {

	// 	// fmt.Printf("\nalready authenticated: Prošlo je!\n")
	// 	struct_user.Hash_lozinka = ""
	// 	return true, struct_user

	// } else {

	// 	err = bcrypt.CompareHashAndPassword([]byte(struct_user.Hash_lozinka), password)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error from bcrypt Dencryption: %s\n", err)
	// 		return false, models.User{}
	// 	} else {
	// 		// fmt.Printf("\nProšlo je!\n")
	// 		struct_user.Hash_lozinka = ""
	// 		return true, struct_user
	// 	}

	// }

	return true

}
