package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	"github.com/vladanan/vezbamo4/src/models"
)

// type User struct {
// 	Hash_lozinka string `db:"hash_lozinka"`
// 	Email        string `db:"email"`
// 	Test         string `db:"test"`
// }

func to_struct(user []byte) []models.User {
	var p []models.User
	err := json.Unmarshal(user, &p)
	if err != nil {
		fmt.Printf("Json error: %v", err)
	}
	return p
}

func AuthenticateUser(email string, password_str string, already_authenticated bool) (bool, models.User) {
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	password := []byte(password_str)

	// ciphertext, err := bcrypt.GenerateFromPassword(password, 5) //df
	// _, err := bcrypt.GenerateFromPassword(password, 5) //df
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error from bcrypt encryption: %s\n", err)
	// 	return false, models.User{}
	// }

	// fmt.Println("Ciphertext: ", string(ciphertext))

	err := godotenv.Load(".env")
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

	rows, err2 := conn.Query(context.Background(), "SELECT * FROM mi_users where email=$1;", email)
	if err2 != nil {
		fmt.Printf("Unable to make query: %v\n", err2)
		return false, models.User{}
	}

	pgx_user, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		//return
		return false, models.User{}
	}

	// fmt.Print("pgx user::", pgx_user)

	bytearray_user, err2 := json.Marshal(pgx_user)
	if err2 != nil {
		fmt.Printf("Json error: %v", err2)
	}

	// fmt.Print("bytearray user:", bytearray_user)

	// for _, item := range to_struct(bytearray_user) {
	// 	fmt.Print("\n")
	// 	fmt.Print(item.Email)
	// 	fmt.Print("\n")
	// 	fmt.Print(item.Test)
	// 	fmt.Print("\n")
	// 	fmt.Print(item.Hash_lozinka)
	// }

	// var time_stamp time.Time
	// err = conn.QueryRow(context.Background(), "select created_at from mi_users where email=$1", email).Scan(&time_stamp)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	// os.Exit(1)
	// 	return false
	// }

	// fmt.Print("\n")
	// fmt.Print(time_stamp)

	var struct_user models.User

	if string(bytearray_user) != "null" {
		struct_user = to_struct(bytearray_user)[0]
	} else {
		return false, models.User{}
	}

	// fmt.Print(struct_user)

	if already_authenticated {

		// fmt.Printf("\nProšlo je!\n")
		struct_user.Hash_lozinka = ""
		return true, struct_user

	} else {

		err = bcrypt.CompareHashAndPassword([]byte(struct_user.Hash_lozinka), password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error from bcrypt Dencryption: %s\n", err)
			return false, models.User{}
		} else {
			// fmt.Printf("\nProšlo je!\n")
			struct_user.Hash_lozinka = ""
			return true, struct_user
		}

	}

}
