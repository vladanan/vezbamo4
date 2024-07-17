package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/vladanan/vezbamo4/src/clr"
	"github.com/vladanan/vezbamo4/src/models"
)

func toStruct(user []byte) []models.User {
	var p []models.User
	err := json.Unmarshal(user, &p)
	if err != nil {
		fmt.Printf("Json error: %v", err)
	}
	return p
}

func AuthenticateUser(email string, passwordStr string, alreadyAuthenticated bool, r *http.Request) (bool, models.User, error) {
	l := clr.GetELRfunc()
	// var Red = "\033[31m"
	var Reset = "\033[0m"
	log.SetFlags(log.Ltime | log.Lshortfile)
	// log.SetPrefix(Red)
	defer func() { log.SetFlags(log.LstdFlags); log.SetPrefix(Reset) }()
	// defer log.SetFlags(log.LstdFlags)

	password := []byte(passwordStr)

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
	pgxUser, e := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if e != nil {
		return l(e)
	}
	bytearrayUser, e := json.Marshal(pgxUser)
	if e != nil {
		return l(e)
	}

	// pgxUser = []models.User{}

	var structUser models.User

	// log.Println("pred dodelu iz niza", pgxUser)

	if len(pgxUser) != 0 {
		structUser = pgxUser[0]
	} else {
		log.Println("prazan user")
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
	pgxSettings, e := pgx.CollectRows(rows2, pgx.RowToStructByName[models.Settings])
	if e != nil {
		return l(e)
	}

	db_bad_sign_in_attempts_limit, err := strconv.ParseInt(pgxSettings[0].Bad_sign_in_attempts_limit, 0, 8)
	if err != nil {
		db_bad_sign_in_attempts_limit = 0
	}
	db_bad_sign_in_time_limit, err := strconv.ParseInt(pgxSettings[0].Bad_sign_in_time_limit, 0, 8)
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

	if alreadyAuthenticated {

		log.Println("Already authenticated")
		structUser.Hash_lozinka = ""
		return true, structUser, nil

	} else if string(bytearrayUser) == "null" { // array je prazan tj. nema korisnika sa takvim mejlom ali se to ne odaje nego se piše i lozinka

		return l(fmt.Sprintln("Nema korisnika sa tim email-om ili lozinkom!")) //, email, password_str

	} else if structUser.Verified_email != "verified" { // ako ima mejla proverava se verifikacija

		return l("Email nije verifikovan!")

		//
	} else if int64(structUser.Bad_sign_in_attempts) < bad_sign_in_attempts_limit { // mejl je verifikovan i ide se na proveru broja neuspelih pokušaja: ako je broj neuspelih pokušaja manji od limita upisuje se pokušaj i ide se na proveru lozinke

		_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
			structUser.Bad_sign_in_attempts+1,
			email,
		)
		if e != nil {
			return l(e)
		} else {

			log.Println("Add 1 to bad sign in attempts:", structUser.Bad_sign_in_attempts)
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
		e = bcrypt.CompareHashAndPassword([]byte(structUser.Hash_lozinka), password) // provera lozinke

		if e != nil { // LOŠA LOZINKA

			return l(fmt.Errorf("Pogrešna lozinka za: %v", email))

		} else { // SVE JE OKEJ (mejl, ver, pokušaji, pass) UPISUJE SE U BAZU SVE ŠTO TREBA I PODACI ŠALJU RUTERU

			_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_time=$1 where email=$2;`, time.Now(), email)
			if e != nil {
				return l(e)
			}
			bytearrayHeaders, e := json.Marshal(r.Header)
			if e != nil {
				return l(e)
			}
			_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_headers=$1 where email=$2;`,
				string(bytearrayHeaders),
				email,
			)
			if e != nil {
				return l(e)
			}
			_, e = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`, 0, email)
			if e != nil {
				return l(e)
			} else {
				log.Println("OK Zeroing bad sign in attempts:", structUser.Bad_sign_in_attempts)
			}

			structUser.Hash_lozinka = ""
			log.Println("AuthenticateUser: Prošlo je!")
			return true, structUser, nil
		}

	} else { // ako je broj neuselih pokušaja veći od limita gleda se da li je prošlo više vremena od limita

		if time.Since(structUser.Bad_sign_in_time).Minutes() < float64(bad_sign_in_time_limit) { //

			return l(fmt.Sprint("Previše pokušaja za sign in. Pokušati za minuta: ",
				float64(bad_sign_in_time_limit)-time.Since(structUser.Bad_sign_in_time).Minutes()))

			//
		} else { // kada je prošlo dovoljno vremena resetuje se broj neuspelih pokušaja

			_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`, 0, email)
			if e != nil {
				return l(e)
			} else {
				log.Println("Time up zeroing bad sign in attempts:", structUser.Bad_sign_in_attempts)
			}

			return l("Moguće je ponovo probati sign in")

		}
	}

}
