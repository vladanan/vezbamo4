package db

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	"github.com/vladanan/vezbamo4/src/models"
)

func to_struct(user []byte) []models.User {
	var p []models.User
	err := json.Unmarshal(user, &p)
	if err != nil {
		fmt.Printf("Json error: %v", err)
	}
	return p
}

func AuthenticateUser(email string, password_str string, already_authenticated bool, r *http.Request) (bool, models.User) {
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	password := []byte(password_str)

	// ENV, BAZA, UZIMANJE USERA

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("AuthenticateUser: Error loading .env file\n")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("AuthenticateUser: Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err2 := conn.Query(context.Background(), "SELECT * FROM mi_users where email=$1;", email)
	if err2 != nil {
		fmt.Printf("AuthenticateUser: Unable to make query: %v\n", err2)
		return false, models.User{}
	}

	pgx_user, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		fmt.Printf("AuthenticateUser: CollectRows error: %v\n", err)
		return false, models.User{}
	}

	// fmt.Print("AuthenticateUser: pgx user:", pgx_user)

	bytearray_user, err2 := json.Marshal(pgx_user)
	if err2 != nil {
		fmt.Printf("AuthenticateUser: user JSON error: %v\n", err2)
	}

	// fmt.Print("bytearray user: ", bytearray_user)

	var struct_user models.User
	if string(bytearray_user) != "null" { // array nije prazan tj. ima zapisa sa odgovarajućim mejlom
		struct_user = pgx_user[0] //to_struct(bytearray_user)[0]
	}

	// fmt.Print("AuthenticateUser: ",struct_user)

	// UZIMANJE PROMENLJIVIH IZ ENV I DB ZA BAD ATTEMPT LIMITE

	var bad_sign_in_attempts_limit int64 = 2
	var bad_sign_in_time_limit int64 = 8

	BAD_SIGN_IN_ATTEMPTS_LIMIT, err := strconv.ParseInt(os.Getenv("BAD_SIGN_IN_ATTEMPTS_LIMIT"), 0, 8)
	if err != nil {
		BAD_SIGN_IN_ATTEMPTS_LIMIT = 0
	}
	BAD_SIGN_IN_TIME_LIMIT, err := strconv.ParseInt(os.Getenv("BAD_SIGN_IN_TIME_LIMIT"), 0, 8)
	if err != nil {
		BAD_SIGN_IN_TIME_LIMIT = 0
	}

	rows2, err2 := conn.Query(context.Background(), "SELECT * FROM v_settings where s_id=1;")
	if err2 != nil {
		fmt.Printf("AuthenticateUser: Unable to make query: %v\n", err2)
		return false, models.User{}
	}
	pgx_settings, err := pgx.CollectRows(rows2, pgx.RowToStructByName[models.Settings])
	if err != nil {
		fmt.Printf("AuthenticateUser: CollectRows error: %v\n", err)
		return false, models.User{}
	}

	db_bad_sign_in_attempts_limit, err := strconv.ParseInt(pgx_settings[0].Bad_sign_in_attempts_limit, 0, 8)
	if err != nil {
		db_bad_sign_in_attempts_limit = 0
	}
	db_bad_sign_in_time_limit, err := strconv.ParseInt(pgx_settings[0].Bad_sign_in_time_limit, 0, 8)
	if err != nil {
		db_bad_sign_in_time_limit = 0
	}

	if BAD_SIGN_IN_ATTEMPTS_LIMIT != 0 {
		bad_sign_in_attempts_limit = BAD_SIGN_IN_ATTEMPTS_LIMIT
	} else if db_bad_sign_in_attempts_limit != 0 {
		bad_sign_in_attempts_limit = db_bad_sign_in_attempts_limit
	}
	if BAD_SIGN_IN_TIME_LIMIT != 0 {
		bad_sign_in_time_limit = BAD_SIGN_IN_TIME_LIMIT
	} else if db_bad_sign_in_time_limit != 0 {
		bad_sign_in_time_limit = db_bad_sign_in_time_limit
	}

	// fmt.Print("AuthenticateUser: bad env: ", BAD_SIGN_IN_ATTEMPTS_LIMIT, BAD_SIGN_IN_TIME_LIMIT, "\n")
	// fmt.Print("AuthenticateUser: bad db: ", db_bad_sign_in_attempts_limit, db_bad_sign_in_time_limit, "\n")
	// fmt.Print("AuthenticateUser: bad real: ", bad_sign_in_attempts_limit, bad_sign_in_time_limit, "\n")

	if already_authenticated {

		fmt.Print("AuthenticateUser: Already authenticated\n")
		struct_user.Hash_lozinka = ""
		return true, struct_user

	} else if int64(struct_user.Bad_sign_in_attempts) < bad_sign_in_attempts_limit { // ako je broj neuselih pokušaja manji od limita ide se na dalje proverese broj neuspelih pokušaja

		fmt.Print("AuthenticateUser: add 1 to bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")
		_, err := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
			struct_user.Bad_sign_in_attempts+1,
			email,
		)
		if err != nil {
			fmt.Printf("AuthenticateUser: Unable to connect to database to write bad sign in attempts:%v\n", err)
			return false, models.User{}
		}

		// fmt.Print("AuthenticateUser: prošlo minuta od poslednjeg lošeg sign in-a: ", time.Since(struct_user.Bad_sign_in_time).Minutes(), "\n")
		fmt.Print("AuthenticateUser: set last bad sign time\n")
		_, err = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_time=$1 where email=$2;`,
			time.Now(),
			email,
		)
		if err != nil {
			fmt.Printf("AuthenticateUser: Unable to connect to database to write bad sign in attempts:%v\n", err)
			return false, models.User{}
		}

		if string(bytearray_user) != "null" { // array nije prazan tj. ima zapisa sa odgovarajućim mejlom

			err = bcrypt.CompareHashAndPassword([]byte(struct_user.Hash_lozinka), password) // provera lozinke

			if err != nil {

				fmt.Fprintf(os.Stderr, "AuthenticateUser: Loša lozinka: %s\n", err)
				return false, models.User{}

			} else if struct_user.Verified_email == "verified" { // ako je lozinka dobra onda se proverava da li je mejl verifikovan

				_, err = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_time=$1 where email=$2;`,
					time.Now(),
					email,
				)
				if err != nil {
					fmt.Print("AuthenticateUser: Unable to connect to database to write last sign in field:", err, "\n")
					return false, models.User{}
				}

				bytearray_headers, err := json.Marshal(r.Header)
				if err != nil {
					fmt.Printf("AuthenticateUser: headers JSON error: %v", err)
				}
				_, err = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_headers=$1 where email=$2;`,
					string(bytearray_headers),
					email,
				)
				if err != nil {
					fmt.Print("AuthenticateUser: Unable to connect to database to write last sign in headers field:", err, "\n")
					return false, models.User{}
				}

				fmt.Print("\nAuthenticateUser: zero bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")
				_, err = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
					0,
					email,
				)
				if err != nil {
					fmt.Print("AuthenticateUser: Unable to connect to database to write bad attempts field:", err, "\n")
					return false, models.User{}
				}

				// fmt.Print("\nAuthenticateUser: Prošlo je!\n")
				struct_user.Hash_lozinka = ""
				return true, struct_user

			} else {

				fmt.Print("AuthenticateUser: Mejl nije verifikovan!\n")
				return false, models.User{}

			}

		} else {

			fmt.Print("AuthenticateUser: nema korisnika sa tim mejlom i lozinkom", email, password_str, "\n")
			return false, models.User{}

		}

	} else { // ako je broj neuselih pokušaja veći od limita gleda se da li je prošlo više vremena od limita

		if time.Since(struct_user.Bad_sign_in_time).Minutes() < float64(bad_sign_in_time_limit) {

			fmt.Print("AuthenticateUser: previše pokušaja za sign in\n")
			fmt.Print("AuthenticateUser: pokušati za minuta:", float64(bad_sign_in_time_limit)-time.Since(struct_user.Bad_sign_in_time).Minutes(), "\n")
			return false, models.User{}

		} else {

			fmt.Print("AuthenticateUser: zeroing bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")
			_, err := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
				0,
				email,
			)
			if err != nil {
				fmt.Printf("AuthenticateUser: Unable to connect to database to write bad sign in attempts:%v\n", err)
				return false, models.User{}
			}

			fmt.Print("AuthenticateUser: moguće je ponovo probati sign in\n")
			return false, models.User{}

		}

	}

}
