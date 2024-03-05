package main

import (
	"context"

	// "crypto/rsa"
	// "crypto/sha512"
	// "crypto/md5"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"

	// "math/big"
	"net/http"
	"os"

	"github.com/vladanan/vezbamo4/db"
	"github.com/vladanan/vezbamo4/views"

	"github.com/a-h/templ"
	// "github.com/BurntSushi/toml"
	// "github.com/nicksnyder/go-i18n/v2/i18n"
	// "golang.org/x/text/language"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var global string = "en"

//const keyServerAddr = "serverAddr"
//curl -X POST -d 'This is the body' 'http://localhost:3333?first=1&second='

func check(e error) {
	if e != nil {
		panic(e)
	}

}

// templ: https://templ.guide/

//https://tailwindcss.com/docs/installation/play-cdn

func getRoot(res http.ResponseWriter, req *http.Request) {
	//ctx := req.Context()

	dat, err := os.ReadFile("views/index.html")
	check(err)

	hasFirst := req.URL.Query().Has("first")
	first := req.URL.Query().Get("first")
	hasSecond := req.URL.Query().Has("second")
	second := req.URL.Query().Get("second")

	body, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Println(req.URL.Path)

	fmt.Printf("got / request. first(%t)=%s, second(%t)=%s, body:\n%s\n",
		//ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second,
		body)

	io.WriteString(res, string(dat))
}

// curl -X POST -F 'myName=Sammy' 'http://localhost:3333/hello'
func getHello(res http.ResponseWriter, req *http.Request) {
	myName := req.PostFormValue("myName")

	if myName == "" {
		res.Header().Set("x-missing-field", "myName")
		res.WriteHeader(http.StatusBadRequest)
		return
		//myName = "HTTP"
	}

	io.WriteString(res, fmt.Sprintf("Hello, %s!\n", myName))
}

func getClicked(res http.ResponseWriter, req *http.Request) {
	fmt.Println("db.Db()")
	fmt.Println(db.Db())

	io.WriteString(res, db.Db())
}

func proba(res http.ResponseWriter, req *http.Request) {
	fmt.Println("get proba")
	dat, err := os.ReadFile("assets/proba.js")
	check(err)

	io.WriteString(res, string(dat))
}

func getHTMX(res http.ResponseWriter, req *http.Request) {
	fmt.Println("get htmx")
	dat, err := os.ReadFile("assets/htmx.min.js")
	check(err)

	io.WriteString(res, string(dat))
}

func getCSS(res http.ResponseWriter, req *http.Request) {
	fmt.Println("get css")
	dat, err := os.ReadFile("src/output.css")
	check(err)

	io.WriteString(res, string(dat))
}

func setEn(res http.ResponseWriter, req *http.Request) {
	fmt.Println("set en")
	global = "en"
}
func setEs(res http.ResponseWriter, req *http.Request) {
	fmt.Println("set es")
	global = "es"
}

func kripto() {

	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	var email = "vladan.andjelkovic@gmail.com"
	password := []byte("vezbamo.2015")

	ciphertext, err := bcrypt.GenerateFromPassword(password, 5) //df
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from bcrypt encryption: %s\n", err)
		return
	}

	fmt.Println("Ciphertext: ", string(ciphertext))

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	//https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx

	type Blog struct {
		B_id   int8   `db:"b_id"`
		Tema   string `db:"tema"`
		Poruka string `db:"poruka"`
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	rows, _ := conn.Query(context.Background(), "SELECT b_id, tema, poruka FROM g_user_blog;")
	if err != nil {
		fmt.Printf("Unable to make query: %v\n", err)
	}

	blogs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Blog])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		//return
	}

	for _, b := range blogs {
		fmt.Printf("%v, %s: $%s\n", b.B_id, b.Poruka, b.Tema)
	}

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

	defer conn.Close(context.Background())

	err = bcrypt.CompareHashAndPassword([]byte(hloz), password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from bcrypt Dencryption: %s\n", err)
		return
	} else {
		fmt.Printf("\nProšlo je!")
	}

}

func main() {

	kripto()

	component := views.Hello("John")
	// component2 = views.Htmx("en")

	//component.Render(context.Background(), os.Stdout)

	// fmt.Println("main db: ")
	// fmt.Println(db.Db())

	http.Handle("/cao", templ.Handler(component))
	// http.Handle("/htmx", templ.Handler(component2))
	http.Handle("/htmx", templ.Handler(views.Htmx(global)))

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/clicked", getClicked)

	http.HandleFunc("/assets/proba.js", proba)
	http.HandleFunc("/assets/htmx.min.js", getHTMX)
	http.HandleFunc("/src/output.css", getCSS)

	http.HandleFunc("/en", setEn)
	http.HandleFunc("/es", setEs)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
