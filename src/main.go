package main

import (
	"context"
	"errors"
	"time"

	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"

	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/vladanan/vezbamo4/db"
	views "github.com/vladanan/vezbamo4/views"
	pitanja "github.com/vladanan/vezbamo4/views/pitanja"

	"github.com/a-h/templ"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var globalLanguage string = ""

//const keyServerAddr = "serverAddr"
//curl -X POST -d 'This is the body' 'http://localhost:3333?first=1&second='

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// templ: https://templ.guide/

//https://tailwindcss.com/docs/installation/play-cdn

// func get404(res http.ResponseWriter, req *http.Request) {
// 	dat, err := os.ReadFile("views/404.html")
// 	check(err)
// 	fmt.Println(req.URL.Path)
// 	io.WriteString(res, string(dat))
// }

func httpPOSTfromHTMX(w http.ResponseWriter, r *http.Request) {
	fmt.Println("db.Db()")
	fmt.Println(db.Db())
	io.WriteString(w, db.Db())
}

// func getProbaJS(res http.ResponseWriter, req *http.Request) {
// 	fmt.Println("get js proba")
// 	dat, err := os.ReadFile("views/proba.js")
// 	check(err)
// 	io.WriteString(res, string(dat))
// }

// func getHTMXlibrary(res http.ResponseWriter, req *http.Request) {
// 	fmt.Println("get htmx library")
// 	dat, err := os.ReadFile("assets/htmx.min.js")
// 	check(err)
// 	io.WriteString(res, string(dat))
// }

// func getTailwindCSS(res http.ResponseWriter, req *http.Request) {
// 	fmt.Println("get css")
// 	dat, err := os.ReadFile("views/output.css")
// 	check(err)
// 	io.WriteString(res, string(dat))
// }

func api() {

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

	fmt.Print(blogs[0])
	// for _, b := range blogs {
	// 	fmt.Printf("%v, %s: $%s\n", b.B_id, b.Poruka, b.Tema)
	// }

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
		fmt.Printf("\nProšlo je!\n")
	}

}

func getLocationsForAngularFE(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nget locations", r.URL)
	dat, err := os.ReadFile("api/locations.json")
	check(err)

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, string(dat))
	//fmt.Println("\ndat: ", string(dat))
	//w.Write(string(dat)) dfaljfa
}

func setEn(w http.ResponseWriter, r *http.Request) {
	globalLanguage = "en-US"
}
func setEs(w http.ResponseWriter, r *http.Request) {
	globalLanguage = "es"
}
func setSr(w http.ResponseWriter, r *http.Request) {
	globalLanguage = "sr"
}
func setBrowserLang(w http.ResponseWriter, r *http.Request) {
	globalLanguage = ""
}

func goToPitanja(w http.ResponseWriter, r *http.Request) {
	templ.Handler(pitanja.Pitanja(globalLanguage, r)).Component.Render(context.Background(), w)
}

func goToIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
	} else {
		templ.Handler(views.Page404()).Component.Render(context.Background(), w)
	}
}

func getPitanja(w http.ResponseWriter, r *http.Request) {
	fmt.Println("main get pitanja")
	spisak := db.GetPitanja()
	templ.Handler(pitanja.Spisak([]byte(spisak))).Component.Render(context.Background(), w)
}

func main() {

	// api()

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", goToIndex)

	//http.Handle("/pitanja", templ.Handler(pitanja.Pitanja()))
	http.HandleFunc("/pitanja", goToPitanja)

	// http.Handle("/404", http.NotFoundHandler())
	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Page404()).Component.Render(context.Background(), w)
	})

	http.HandleFunc("/zgrabi-iz-db", httpPOSTfromHTMX)
	http.HandleFunc("/get_pitanja", getPitanja)
	// http.HandleFunc("/proba.js", getProbaJS)
	// http.HandleFunc("/assets/htmx.min.js", getHTMXlibrary)
	// http.HandleFunc("/output.css", getTailwindCSS)

	http.HandleFunc("/en", setEn)
	http.HandleFunc("/es", setEs)
	http.HandleFunc("/sr", setSr)
	http.HandleFunc("/browser", setBrowserLang)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})

	http.Handle("/locations", c.Handler(http.HandlerFunc(getLocationsForAngularFE)))
	http.Handle("/locations/", c.Handler(http.HandlerFunc(getLocationsForAngularFE)))

	fmt.Println("Done", time.Now().Second())

	var err = http.ListenAndServe("0.0.0.0:10000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
