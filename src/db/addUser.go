package db

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	// "crypto/tls"
	"fmt"

	// "log"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/vladanan/vezbamo4/src/models"

	"os"

	// "time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	// "encoding/json"
	// "github.com/vladanan/vezbamo4/src/models"
)

func to_map(headers []byte) map[string][]string {
	var h map[string][]string
	err := json.Unmarshal(headers, &h)
	if err != nil {
		fmt.Printf("AddUser: to_map:%v\n", err)
		// return nil, err
	}
	return h
}

func AddUser(email_str, user_name, password_str string, r *http.Request) bool {

	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/
	// GenerateFromPassword does not accept passwords longer than 72 bytes, which is the longest password bcrypt will operate on.
	// nije dobro da se key za proveru mejla pravi na osnovu lozinke
	// da se ne bi desilo da neko proba da rekonstruiše password iz poslatog linka
	// najbolje samo iz mejla jer je mejl svakako već poznat onome ko ima link a novi link svakako ne može sam da generiše

	password := []byte(password_str)
	email := []byte(email_str)

	ciphertext_sign_in, err := bcrypt.GenerateFromPassword(password, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddUser: greška bcrypt ciphertext_sign_in: %s\n", err)
		return false
	}
	ciphertext_verify_mail, err := bcrypt.GenerateFromPassword([]byte(email), 7)
	if err != nil {
		fmt.Fprintf(os.Stderr, "AddUser: greška bcrypt ciphertext_verify_mail: %s\n", err)
		return false
	}

	// kada se key koristi bez zamena / i . onda ne može da se koristi kao url jer / dovodi do toga da je url pogrešan
	// možda se to ne dešava sa . ali sam zamenio za svaki slučaj da se ne brka domen ili tako nešto
	ciphertext_verify_mail_string1 := strings.ReplaceAll(string(ciphertext_verify_mail), "/", "-")
	ciphertext_verify_mail_string2 := strings.ReplaceAll(string(ciphertext_verify_mail_string1), ".", "_")

	// fmt.Println("Ciphertexts: ", string(ciphertext_sign_in))

	// UZIMAMO HEADER ZA KASNIJE POREĐENJE X-FORWARDED-FOR I ZA UBACIVANJE U DB POLJE RADI EVIDENCIJE I POREĐENJA
	bytearray_headers, err := json.Marshal(r.Header)
	if err != nil {
		fmt.Printf("AddUser: json 1: %v", err)
	}

	// UZIMA SE HTML FAJL ZA MEJL ZA VERIFIKACIJU
	// https://gobyexample.com/reading-files
	dat, err1 := os.ReadFile("src/html/verify_email.html")
	if err1 != nil {
		fmt.Printf("AddUser: greška čitanje html fajla: %v\n", err1)
		return false
	}
	// fmt.Print("AddUser: html fajl:", string(dat))
	html := string(dat)

	// DB
	err2 := godotenv.Load(".env")
	if err2 != nil {
		fmt.Printf("AddUser: greška učitavanja za .env")
		return false // uraditi da ne bude samo return false nego neku fukciju koja radi ERROR i LOG na srpskom a i return false uz neku poruku za fe stranu tako da radi sa i18n
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("AddUser: povezivanje sa db: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// PROVERA DA LI NEMA VEĆ USER-A SA ISTIM MEJLOM I USER_NAME

	rows, err2 := conn.Query(context.Background(), "SELECT * FROM mi_users where email=$1 OR user_name=$2;", email_str, user_name)
	if err2 != nil {
		fmt.Printf("AddUser: Query 1: %v\n", err2)
		return false
	}
	pgx_user, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		fmt.Printf("AddUser: CollectRows 1: %v\n", err)
		return false
	}
	// fmt.Print("AddUser: pgx user:", pgx_user)
	bytearray_user, err2 := json.Marshal(pgx_user)
	if err2 != nil {
		fmt.Printf("AddUser: json 2: %v\n", err2)
		return false
	}
	// fmt.Print("bytearray user: ", bytearray_user)

	if string(bytearray_user) != "null" { // array nije prazan tj. ima zapisa sa istim mejlom ili userom

		fmt.Print("AddUser: već ima korisnik sa takvim mejlom/user name\n")
		return false

	} else {

		// PROVERA DA LI ima istog takvog hedera u ostalim zapisima deset minuta unazad tj. koji imaju created_at deset minuta stariji od sadašnjeg trenutka

		// fmt.Print("AddUser: r header: ", to_map(bytearray_headers)["X-Forwarded-For"][0], "\n")

		// https://stackoverflow.com/questions/23320945/postgresql-select-if-string-contains
		// https://stackoverflow.com/questions/45849494/how-do-i-search-for-a-specific-string-in-a-json-postgres-data-type-column
		// razni sql pokušaji i varijante sa i bez promenljivih i jsonb polja:
		// SELECT id FROM TAG_TABLE WHERE position(tag_name in 'aaaaaaaaaaa')>0;
		// rows2, err2 := conn.Query(context.Background(), "SELECT * FROM mi_users where position($1 IN created_at_headers) > 0;", to_map(bytearray_headers)["X-Forwarded-For"][0])
		// rows2, err2 := conn.Query(context.Background(), `with  vars as (select '127.0.0.1' as var1) select * from  mi_users,  vars where jsonb_path_exists(created_at_headers,'$.X-Forwarded-For ? (@ == var1)');`, to_map(bytearray_headers)["X-Forwarded-For"][0])

		// ne mogu da se koriste prepared statements niti sql promenljive u upitima za sadržaj X-Forwarded-For heder u jsonb polju (može sa string concatenation ali to je opasno) tako da mora prvo da se pokupe upisi u poslednjih 10 min i zatim da se kod svih uporedi X-Forwarded-For sa aktuelnim

		// UZIMANJE PROMENLJIVIH IZ ENV I DB ZA ATTEMPT TIME LIMIT

		var same_ip_sign_up_time_limit = "2m"

		SAME_IP_SIGN_UP_TIME_LIMIT := os.Getenv("SAME_IP_SIGN_UP_TIME_LIMIT")
		if SAME_IP_SIGN_UP_TIME_LIMIT == "" || SAME_IP_SIGN_UP_TIME_LIMIT == "0" {
			SAME_IP_SIGN_UP_TIME_LIMIT = "0m"
		}

		rows2, err := conn.Query(context.Background(), "SELECT * FROM v_settings where s_id=1;")
		if err != nil {
			fmt.Printf("AddUser: Unable to make query: %v\n", err)
			return false
		}
		pgx_settings, err := pgx.CollectRows(rows2, pgx.RowToStructByName[models.Settings])
		if err != nil {
			fmt.Printf("AddUser: CollectRows error: %v\n", err)
			return false
		}

		db_same_ip_sign_up_time_limit := pgx_settings[0].Same_ip_sign_up_time_limit
		if db_same_ip_sign_up_time_limit == "" || db_same_ip_sign_up_time_limit == "0" {
			db_same_ip_sign_up_time_limit = "0m"
		}

		if SAME_IP_SIGN_UP_TIME_LIMIT != "0m" {
			same_ip_sign_up_time_limit = "'" + SAME_IP_SIGN_UP_TIME_LIMIT + "m'"
		} else if db_same_ip_sign_up_time_limit != "0m" {
			same_ip_sign_up_time_limit = "'" + db_same_ip_sign_up_time_limit + "m'"
		}

		fmt.Print("AddUser: bad env: ", SAME_IP_SIGN_UP_TIME_LIMIT, db_same_ip_sign_up_time_limit, same_ip_sign_up_time_limit, "\n")

		rows2, err2 := conn.Query(
			context.Background(),
			`select * from  mi_users where (now() :: timestamp - created_at_time) < interval `+same_ip_sign_up_time_limit)
		// rows2, err2 := conn.Query(context.Background(), `select * from  mi_users where (now() :: timestamp - created_at_time) < interval '3m'`)
		if err2 != nil {
			fmt.Printf("AddUser: query 2: %v\n", err2)
			return false
		}
		pgx_user2, err := pgx.CollectRows(rows2, pgx.RowToStructByName[models.User])
		if err != nil {
			fmt.Printf("AddUser: CollectRows 2: %v\n", err)
			return false
		}

		// fmt.Print("AddUser: pgx:", pgx_user2, "\n")

		if pgx_user2 == nil {

			fmt.Print("AddUser: nema upisa od pre " + same_ip_sign_up_time_limit + ":\n")

		} else {

			for _, item := range pgx_user2 {

				if to_map([]byte(item.Created_at_headers))["X-Forwarded-For"][0] == to_map(bytearray_headers)["X-Forwarded-For"][0] {

					fmt.Print("AddUser: ima upis u posledjih "+same_ip_sign_up_time_limit+" i JESTE isti ip:", to_map([]byte(item.Created_at_headers))["X-Forwarded-For"][0], "\n")
					return false

				} else {

					fmt.Print("AddUser: ima upis u posledjih "+same_ip_sign_up_time_limit+" ali NIJE isti ip:", to_map([]byte(item.Created_at_headers))["X-Forwarded-For"][0], "\n")

				}

			}

		}

	}

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
			verified_email,
			created_at_headers
		)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		ciphertext_sign_in,
		email_str,
		user_name,
		"user",
		0,
		true, true, false,
		ciphertext_verify_mail_string2,
		string(bytearray_headers),
		// time.Now(),
		// time.DateTime,
	)
	if err != nil {

		fmt.Printf("AddUser: ne može da se uradi insert u db: %v\n", err)
		return false

	} else {

		fmt.Printf("AddUser: insert rezultat: %v\n", commandTag)

		// AKO JE SVE OKEJ ŠALJE SE MEJL PREKO MEJL KLIJENTA

		// Set up authentication information.
		auth := sasl.NewPlainClient("", os.Getenv("SMTP_MAIL"), os.Getenv("SMTP_APP_PASSWORD_VEZBAMO"))

		// Connect to the server, authenticate, set the sender and recipient,and send the email all in one step.

		var mail_for_mail string
		var url_domain_for_mail string

		if os.Getenv("PRODUCTION") == "FALSE" {
			mail_for_mail = email_str //"vladan_zasve@yahoo.com"
			url_domain_for_mail = "http://127.0.0.1:7331/vmk/" + string(ciphertext_verify_mail_string2) + "?mail=" + mail_for_mail
		} else {
			mail_for_mail = email_str
			url_domain_for_mail = "https://vezbamo.onrender.com/vmk/" + string(ciphertext_verify_mail_string2) + "?mail=" + mail_for_mail
		}

		to := []string{mail_for_mail}

		html = strings.ReplaceAll(html, "+user_name+", user_name)
		html = strings.ReplaceAll(html, "+url_domain_for_mail+", url_domain_for_mail)
		html = strings.ReplaceAll(html, "+mail_for_mail+", mail_for_mail)
		// fmt.Print("AddUser: html za mejl:", html)

		msg := strings.NewReader(
			`Content-Transfer-Encoding: quoted-printable` + "\r\n" +
				`Content-Type: text/html; charset="UTF-8"` + "\r\n" +
				`To: ` + mail_for_mail + "\r\n" +
				`Subject: Dobrodošli na portal Vežbamo!` + "\r\n" +
				"\r\n" +

				html +

				// `<html lang="en-US"><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"/></head><body>` +
				// `<h1 style="color: darkblue">Zdravo ` + user_name + `!<h1>` +
				// `<h2>Za verifikaciju svog naloga prekopiraj ovaj link u svoj browser:<h2>` +
				// `<p style="font-size: small">` + url_domain_for_mail + `<p>` +
				// `<p>ili klikni na ovo dugme za verifikaciju mejla:<p>` +
				// `<a href="` + url_domain_for_mail + `"><button type="button" style="color: darkgreen; font-size: large; font-weight: bold">Verifikacija mejla</button></a><br></br>` +
				// `<p style="font-size: small">Ako "` + mail_for_mail + `" nije tvoj email ili ne želiš da se registruješ na portal Vežbamo javi nam se na email adresu: y.emailbox-vezbamo@yahoo.com .<p>` +
				// `</body></html>` + "\r\n" +

				"\r\n")

		err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("SMTP_MAIL"), to, msg)
		if err != nil {
			fmt.Printf("AddUser: korisnik je upisan ali nije poslat email za verifikaciju:%v\n", err)
			return false
			// log.Fatal(err)
		}

		return true

	}

}

/*
<html dir=3D"ltr" lang=3D"en">=0A=0A  <head>=0A    <meta http-equiv=3D"Cont=
ent-Type" content=3D"text/html; charset=3Dutf-8" />

*/
