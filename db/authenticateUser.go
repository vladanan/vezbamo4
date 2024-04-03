package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(email string, password_str string) bool {
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	password := []byte(password_str)

	// ciphertext, err := bcrypt.GenerateFromPassword(password, 5) //df
	_, err := bcrypt.GenerateFromPassword(password, 5) //df
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from bcrypt encryption: %s\n", err)
		return false
	}

	// fmt.Println("Ciphertext: ", string(ciphertext))

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// commandTag, err := conn.Exec(context.Background(), "INSERT INTO mi_users (email, basic, js, c, hash_lozinka) VALUES ($1, $2, $3, $4, $5);", email, true, true, true, ciphertext)
	// if err != nil {
	// 	fmt.Printf("Unable to connect to database: %v\n", err)
	// }
	// fmt.Printf("insert result: %v\n", commandTag)

	// err = bcrypt.CompareHashAndPassword([]byte("$2a$05$HYej4fyvWYnC5LvrSGEOD.bztJzcYn45t62etOTrN8d59BkoD7fhy"), password)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error from bcrypt Dencryption: %s\n", err)
	// 	return
	// } else {
	// 	fmt.Printf("\nProšlo je!")
	// }

	var hloz string
	err = conn.QueryRow(context.Background(), "select hash_lozinka from mi_users where email=$1", email).Scan(&hloz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hloz), password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from bcrypt Dencryption: %s\n", err)
		return false
	} else {
		// fmt.Printf("\nProšlo je!\n")
		return true
	}

}
