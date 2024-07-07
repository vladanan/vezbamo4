package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/vladanan/vezbamo4/src/errorlogres"
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

func AuthenticateUser(email string, password_str string, already_authenticated bool, r *http.Request) (bool, models.User, string) {
	l := errorlogres.GetELRfunc()
	log.SetFlags(log.Ltime | log.Lshortfile)
	password := []byte(password_str)

	// if _, e := strconv.Atoi("v"); e != nil {
	// 	l("neka greškaaaaa bez veze")
	// 	return l(e)
	// }

	// ENV, BAZA, UZIMANJE USERA 323sfds

	e := godotenv.Load(".env")
	if e != nil {
		return l(e)
	}
	conn, e := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if e != nil {
		return l(e)
		// os.Exit(1)
	}
	defer conn.Close(context.Background())
	rows, e := conn.Query(context.Background(), "SELECT * FROM mi_users where email=$1;", email)
	if e != nil {
		return l(e)
	}
	pgx_user, e := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if e != nil {
		return l(e)
	}
	bytearray_user, e := json.Marshal(pgx_user)
	if e != nil {
		return l(e)
	}
	var struct_user models.User
	if string(bytearray_user) != "null" { // array nije prazan tj. ima zapisa sa odgovarajućim mejlom
		struct_user = pgx_user[0] // to_struct(bytearray_user)[0]
	}

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
		return l(e)
	}
	pgx_settings, e := pgx.CollectRows(rows2, pgx.RowToStructByName[models.Settings])
	if e != nil {
		return l(e)
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

	if already_authenticated {

		log.Println("Already authenticated")
		struct_user.Hash_lozinka = ""
		return true, struct_user, ""

	} else if string(bytearray_user) == "null" { // array je prazan tj. nema korisnika sa takvim mejlom ali se to ne odaje nego se piše i lozinka

		return l(fmt.Sprintln("Nema korisnika sa tim mejlom ili lozinkom!")) //, email, password_str

	} else if struct_user.Verified_email != "verified" { // ako ima mejla proverava se verifikacija

		return l("Mejl nije verifikovan!\n")

		//
	} else if int64(struct_user.Bad_sign_in_attempts) < bad_sign_in_attempts_limit { // mejl je verifikovan i ide se na proveru broja neuspelih pokušaja: ako je broj neuspelih pokušaja manji od limita upisuje se pokušaj i ide se na proveru lozinke

		_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
			struct_user.Bad_sign_in_attempts+1,
			email,
		)
		if e != nil {
			return l(e)
		} else {

			log.Println("Add 1 to bad sign in attempts:", struct_user.Bad_sign_in_attempts)
		}

		_, e = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_time=$1 where email=$2;`,
			time.Now(),
			email)
		if e != nil {
			return l(e)
		} else {
			log.Println("Set last bad sign time")
		}

		// https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
		// https://gowebexamples.com/password-hashing/
		e = bcrypt.CompareHashAndPassword([]byte(struct_user.Hash_lozinka), password) // provera lozinke

		if e != nil { // LOŠA LOZINKA

			return l(errors.New(fmt.Sprintln("Pogrešna lozinka za:", email)))

		} else { // SVE JE OKEJ (mejl, ver, pokušaji, pass) UPISUJE SE U BAZU SVE ŠTO TREBA I PODACI ŠALJU RUTERU

			_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_time=$1 where email=$2;`, time.Now(), email)
			if e != nil {
				return l(e)
			}
			bytearray_headers, e := json.Marshal(r.Header)
			if e != nil {
				return l(e)
			}
			_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_headers=$1 where email=$2;`,
				string(bytearray_headers),
				email,
			)
			if e != nil {
				return l(e)
			}
			_, e = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`, 0, email)
			if e != nil {
				return l(e)
			} else {
				log.Println("OK Zeroing bad sign in attempts:", struct_user.Bad_sign_in_attempts)
			}

			struct_user.Hash_lozinka = ""
			log.Println("AuthenticateUser: Prošlo je!")
			return true, struct_user, ""
		}

	} else { // ako je broj neuselih pokušaja veći od limita gleda se da li je prošlo više vremena od limita

		if time.Since(struct_user.Bad_sign_in_time).Minutes() < float64(bad_sign_in_time_limit) { //

			return l(fmt.Sprint("Previše pokušaja za sign in. Pokušati za minuta: ",
				float64(bad_sign_in_time_limit)-time.Since(struct_user.Bad_sign_in_time).Minutes()))

			//
		} else { // kada je prošlo dovoljno vremena resetuje se broj neuspelih pokušaja

			_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`, 0, email)
			if e != nil {
				return l(e)
			} else {
				log.Println("Time up zeroing bad sign in attempts:", struct_user.Bad_sign_in_attempts)
			}

			return l("Moguće je ponovo probati sign in")

		}
	}

}
