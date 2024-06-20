package db

import (
	"context"
	"strings"

	// "crypto/tls"
	"fmt"

	// "log"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"

	"os"

	// "time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	// "encoding/json"
	// "github.com/vladanan/vezbamo4/src/models"
)

func AddUser(email_str, user_name, password_str string) bool {

	// PROVERA DA LI NEMA VEĆ USER-A SA ISTIM MEJLOM I USER_NAME
	// PROVERA DA LI ima već upisa sa istog ili sičnog ip u poslednjih 10min u odnosu na created_at

	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	// GenerateFromPassword does not accept passwords longer than 72 bytes, which is the longest password bcrypt will operate on.
	// praviti da se key za proveru mejla pravi na osnovu mejla i lozinke jer su skupa 32+32=64 ispod broja koji prihvata bcrypt

	password := []byte(password_str)

	ciphertext_sign_in, err := bcrypt.GenerateFromPassword(password, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddUser: Error from bcrypt encryption ciphertext_sign_in: %s\n", err)
		return false
	}
	ciphertext_verify_mail, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddUser: Error from bcrypt encryption ciphertext_verify_mail: %s\n", err)
		return false
	}

	// kada se key koristi bez zamena / i . onda ne može da se koristi kao url jer / dovodi do toga da je url pogrešan
	// možda se to ne dešava sa . ali sam zamenio za svaki slučaj da se ne brka domen ili tako nešto
	ciphertext_verify_mail_string1 := strings.ReplaceAll(string(ciphertext_verify_mail), "/", "-")
	ciphertext_verify_mail_string2 := strings.ReplaceAll(string(ciphertext_verify_mail_string1), ".", "=")

	// fmt.Println("Ciphertexts: ", string(ciphertext_sign_in))

	err2 := godotenv.Load(".env")
	if err2 != nil {
		fmt.Printf("AddUser: Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("AddUser: Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	commandTag, err := conn.Exec(context.Background(), `INSERT INTO mi_users
		(
			hash_lozinka,
			email,
			user_name,
			user_mode,
			user_level,
			basic,
			js,
			c,
			verified_email
		)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		ciphertext_sign_in,
		email_str,
		user_name,
		"user",
		0,
		true, true, false,
		ciphertext_verify_mail_string2,
		// time.Now(),
		// time.DateTime,
	)
	if err != nil {

		fmt.Printf("AddUser: Unable to execute insert: %v\n", err)
		return false

	} else {

		fmt.Printf("AddUser: insert result: %v\n", commandTag)

		// AKO JE SVE OKEJ ŠALJE SE MEJL PREKO MEJL KLIJENTA

		// Set up authentication information.
		auth := sasl.NewPlainClient("", os.Getenv("SMTP_MAIL"), os.Getenv("SMTP_APP_PASSWORD_VEZBAMO"))

		// Connect to the server, authenticate, set the sender and recipient,
		// and send the email all in one step.

		var mail_for_mail string
		var url_domain_for_mail string

		if os.Getenv("PRODUCTION") == "FALSE" {
			mail_for_mail = "vladan_zasve@yahoo.com"
			url_domain_for_mail = "http://127.0.0.1:7331/vmk/" + string(ciphertext_verify_mail_string2)
		} else {
			mail_for_mail = email_str
			url_domain_for_mail = "https://vezbamo.onrender.com/vmk/" + string(ciphertext_verify_mail_string2)
		}

		to := []string{mail_for_mail}

		msg := strings.NewReader(
			"To: " + mail_for_mail + "\r\n" +
				"Subject: Dobrodošli na portal Vežbamo!\r\n" +
				"\r\n" +
				"Da bi verifikovao svoj nalog prekopiraj ovaj link u svoj browser: " + url_domain_for_mail + "\r\n" +
				"\r\n")

		err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("SMTP_MAIL"), to, msg)
		if err != nil {
			fmt.Printf("AddUser: Unable to send mail:%v\n", err)
			return false
			// log.Fatal(err)
		}

		return true

	}

}
