package main

import (
	"context"
	"errors"
	"time"

	"fmt"
	"io"

	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/vladanan/vezbamo4/db"

	views "github.com/vladanan/vezbamo4/views"
	questions "github.com/vladanan/vezbamo4/views/questions"
	site "github.com/vladanan/vezbamo4/views/site"

	"github.com/a-h/templ"

	"github.com/gorilla/sessions"
)

var globalLanguage string = ""

//const keyServerAddr = "serverAddr"
//curl -X POST -d 'This is the body' 'http://localhost:3333?first=1&second='

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// templ: https://templ.guide/ fsfas

//https://tailwindcss.com/docs/installation/play-cdn

// func get404(res http.ResponseWriter, req *http.Request) {
// 	dat, err := os.ReadFile("views/404.html")
// 	check(err)
// 	fmt.Println(req.URL.Path)
// 	io.WriteString(res, string(dat))
// }

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

func goToIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
	} else {
		templ.Handler(views.Page404()).Component.Render(context.Background(), w)
	}
}

func goToQuestions(w http.ResponseWriter, r *http.Request) {
	templ.Handler(questions.Questions(globalLanguage, r)).Component.Render(context.Background(), w)
}

func htmxGetQuestions(w http.ResponseWriter, r *http.Request) {
	list := db.GetQuestions()
	templ.Handler(questions.List(list)).Component.Render(context.Background(), w)
}

func APIgetQuestions(w http.ResponseWriter, r *http.Request) {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// curl http://127.0.0.1:7331/api_questions
	// io.WriteString(w, string(db.GetQuestions()))
	w.Write(db.GetQuestions())
}

func goToNotes(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(db.GetNotes())
	templ.Handler(site.Notes(globalLanguage, r, db.GetNotes())).Component.Render(context.Background(), w)
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func user(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		templ.Handler(site.UserNotLogedPage(globalLanguage, r)).Component.Render(context.Background(), w)
		return
	}

	// Print secret message
	// fmt.Println("You are logged in!")
	templ.Handler(site.UserPage(globalLanguage, r)).Component.Render(context.Background(), w)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	// Authentication goes here
	// ...
	email := "vladan.andjelkovic@gmail.com"
	password := "vezbamo.2015"
	authenticated := db.AuthenticateUser(email, password)
	// Set user as authenticated
	if authenticated {
		session.Values["authenticated"] = true
		session.Save(r, w)
		user(w, r)
	} else {
		session.Values["authenticated"] = false
		session.Save(r, w)
		templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
}

func main() {

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", goToIndex)

	//http.Handle("/questions", templ.Handler(questions.Questions()))
	http.HandleFunc("/questions", goToQuestions)
	http.HandleFunc("/htmx_get_questions", htmxGetQuestions)
	http.HandleFunc("/api_questions", APIgetQuestions)

	http.HandleFunc("/notes", goToNotes)

	// http.Handle("/404", http.NotFoundHandler())
	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Page404()).Component.Render(context.Background(), w)
	})

	// http.HandleFunc("/proba.js", getProbaJS)
	// http.HandleFunc("/assets/htmx.min.js", getHTMXlibrary)
	// http.HandleFunc("/output.css", getTailwindCSS)

	http.HandleFunc("/en", setEn) //dasfa
	http.HandleFunc("/es", setEs)
	http.HandleFunc("/sr", setSr)
	http.HandleFunc("/browser", setBrowserLang)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})

	http.Handle("/locations", c.Handler(http.HandlerFunc(getLocationsForAngularFE)))
	http.Handle("/locations/", c.Handler(http.HandlerFunc(getLocationsForAngularFE)))

	http.HandleFunc("/user", user)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	//http://127.0.0.1:7331
	//curl -s http://127.0.0.1:7331/secret
	//curl -s -I http://127.0.0.1:7331/login
	// cookie-name=MTcxMjE2ODg3MnxEdi1CQkFFQ180SUFBUkFCRUFBQUpmLUNBQUVHYzNSeWFXNW5EQThBRFdGMWRHaGxiblJwWTJGMFpXUUVZbTl2YkFJQ0FBRT18PIMIqmKy6k41-1TZIkA9j7QXEQ79mZmcVIJPsKONzQQ=; Path=/; Expires=Fri, 03 May 2024 18:27:52 GMT; Max-Age=2592000

	// curl -s --cookie "cookie-name=MTcxMjE2ODg3MnxEdi1CQkFFQ180SUFBUkFCRUFBQUpmLUNBQUVHYzNSeWFXNW5EQThBRFdGMWRHaGxiblJwWTJGMFpXUUVZbTl2YkFJQ0FBRT18PIMIqmKy6k41-1TZIkA9j7QXEQ79mZmcVIJPsKONzQQ=; Path=/; Expires=Fri, 03 May 2024 18:27:52 GMT; Max-Age=2592000" http://127.0.0.1:7331/secret

	fmt.Println("Main done", time.Now().Second())

	var err = http.ListenAndServe("0.0.0.0:10000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
