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

	model "github.com/vladanan/vezbamo4/src/models"

	"log"
)

func to_struct(user []byte) []model.User {
	var p []model.User
	err := json.Unmarshal(user, &p)
	if err != nil {
		fmt.Printf("Json error: %v", err)
	}
	return p
}

func AuthenticateUser(email string, password_str string, already_authenticated bool, r *http.Request) (bool, model.User) {
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	l := log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
	f := false
	u := model.User{}

	password := []byte(password_str)

	// ENV, BAZA, UZIMANJE USERA

	e := godotenv.Load(".env")
	if e != nil {
		l.Print(e)
		return f, u
	}

	conn, e := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if e != nil {
		l.Print(e)
		return f, u
		// os.Exit(1)
	}
	defer conn.Close(context.Background())
	rows, e := conn.Query(context.Background(), "SELECT * FROM mi_users777 where email=$1;", email)
	if e != nil {
		l.Print(e)
		return f, u
	}
	pgx_user, e := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if e != nil {
		l.Print(e)
		return f, u
	}
	// fmt.Print("AuthenticateUser: pgx user:", pgx_user)
	bytearray_user, e := json.Marshal(pgx_user)
	if e != nil {
		l.Print(e)
		return f, u
	}

	// fmt.Print("bytearray user: ", bytearray_user)

	var struct_user model.User
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

	rows2, e := conn.Query(context.Background(), "SELECT * FROM v_settings where s_id=1;")
	if e != nil {
		l.Print(e)
		return f, u
	}
	pgx_settings, e := pgx.CollectRows(rows2, pgx.RowToStructByName[model.Settings])
	if e != nil {
		l.Print(e)
		return f, u
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

		_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
			struct_user.Bad_sign_in_attempts+1,
			email,
		)
		if e != nil {
			l.Print(e)
			return f, u
		}

		// fmt.Print("AuthenticateUser: prošlo minuta od poslednjeg lošeg sign in-a: ", time.Since(struct_user.Bad_sign_in_time).Minutes(), "\n")
		fmt.Print("AuthenticateUser: set last bad sign time\n")
		_, e = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_time=$1 where email=$2;`,
			time.Now(),
			email,
		)
		if e != nil {
			l.Print(e)
			return f, u
		}

		if string(bytearray_user) != "null" { // array nije prazan tj. ima zapisa sa odgovarajućim mejlom

			e = bcrypt.CompareHashAndPassword([]byte(struct_user.Hash_lozinka), password) // provera lozinke

			if e != nil {
				// fmt.Fprintf(os.Stderr, "AuthenticateUser: Loša lozinka: %s\n", err)
				l.Print(e)
				return f, u

			} else if struct_user.Verified_email == "verified" { // ako je lozinka dobra onda se proverava da li je mejl verifikovan

				_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_time=$1 where email=$2;`,
					time.Now(),
					email,
				)
				if e != nil {
					l.Print(e)
					return f, u
				}

				bytearray_headers, e := json.Marshal(r.Header)
				if e != nil {
					l.Print(e)
					return f, u
				}
				_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_headers=$1 where email=$2;`,
					string(bytearray_headers),
					email,
				)
				if e != nil {
					l.Print(e)
					return f, u
				}

				fmt.Print("\nAuthenticateUser: zero bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")

				_, e = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
					0,
					email,
				)
				if e != nil {
					l.Print(e)
					return f, u
				}

				// fmt.Print("\nAuthenticateUser: Prošlo je!\n")
				struct_user.Hash_lozinka = ""
				return true, struct_user

			} else {

				l.Print("Mejl nije verifikovan!\n")
				return f, u

			}

		} else {

			l.Print("Nema korisnika sa tim mejlom i lozinkom", email, password_str, "\n")
			return f, u

		}

	} else { // ako je broj neuselih pokušaja veći od limita gleda se da li je prošlo više vremena od limita

		if time.Since(struct_user.Bad_sign_in_time).Minutes() < float64(bad_sign_in_time_limit) {

			l.Print("AuthenticateUser: previše pokušaja za sign in\n")
			l.Print("AuthenticateUser: pokušati za minuta:", float64(bad_sign_in_time_limit)-time.Since(struct_user.Bad_sign_in_time).Minutes(), "\n")
			return f, u

		} else {

			fmt.Print("AuthenticateUser: zeroing bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")
			_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
				0,
				email,
			)
			if e != nil {
				l.Print(e)
				return f, u
			}

			l.Print("Moguće je ponovo probati sign in\n")
			return f, u

		}

	}

}
