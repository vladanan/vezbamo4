package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/vladanan/vezbamo4/src/db"
	views "github.com/vladanan/vezbamo4/src/views"
	questions "github.com/vladanan/vezbamo4/src/views/questions"
	site "github.com/vladanan/vezbamo4/src/views/site"
)

var globalLanguage string = ""

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var godotevn_err = godotenv.Load(".env")

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	// key   = []byte("super-secret-key")
	key   = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)
)

//***  P A G E S

// http://127.0.0.1:7331

func GoToIndex(w http.ResponseWriter, r *http.Request) {

	// if godotevn_err != nil {
	// 	fmt.Printf("Error loading .env file")
	// }

	// session, err1 := store.Get(r, "vezbamo.onrender.com-lang")
	// if err1 != nil {
	// 	fmt.Println("index greška get sessio")
	// 	// http.Error(w, err1.Error(), http.StatusInternalServerError)
	// 	// return
	// }

	// // Set some session values. ghghsdhg
	// // session.Values["language"] = "srpski"

	// session.Options = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   86400 * 7,
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteStrictMode,
	// 	// SameSite: http.SameSite(0),
	// }

	// err2 := session.Save(r, w)
	// if err2 != nil {
	// 	fmt.Println("index greška save sessio")
	// 	// http.Error(w, err2.Error(), http.StatusInternalServerError)
	// 	// return
	// }

	if r.URL.Path == "/" {
		templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
	} else {
		templ.Handler(site.Page404()).Component.Render(context.Background(), w)
	}
}

func GoToQuestions(w http.ResponseWriter, r *http.Request) {
	templ.Handler(questions.Questions(globalLanguage, r)).Component.Render(context.Background(), w)
}

func GoToNotes(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(db.GetNotes())
	templ.Handler(site.Notes(globalLanguage, r, db.GetNotes())).Component.Render(context.Background(), w)
}

func Login(w http.ResponseWriter, r *http.Request) {

	// https://pkg.go.dev/github.com/gorilla/sessions@v1.2.2#section-documentation
	// https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-cookie-same-site-00

	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	// session.Values["language"] = "sr"

	// Authentication goes here
	// ...
	email := "vladan.andjelkovic@gmail.com"
	password := "vezbamo.2015"
	authenticated := db.AuthenticateUser(email, password)

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSite(0),
	}

	// Set user as authenticated
	if authenticated {
		session.Values["authenticated"] = true
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Admin(w, r)
	} else {
		session.Values["authenticated"] = false
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
	}
}

func Admin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")

	// fmt.Println("jezik:", session.Values["language"])

	// Check if user is authenticated
	// auth2, ok2 := session.Values["authenticated"].(bool)
	// fmt.Println("pristup admin sajtu: auth:", auth2, "ok:", ok2)
	// ulogovan: pristup admin sajtu: auth: true ok: true
	// neulogovan: pristup admin sajtu: auth: false ok: true
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		templ.Handler(site.UserNotLogedPage(globalLanguage, r)).Component.Render(context.Background(), w)
		return
	}

	// Print secret message
	// fmt.Println("You are logged in!")
	templ.Handler(site.Admin(globalLanguage, r)).Component.Render(context.Background(), w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
}

//*** I N T E R N A L

func HtmxGetQuestions(w http.ResponseWriter, r *http.Request) {
	list := db.GetQuestions()
	templ.Handler(questions.List(list)).Component.Render(context.Background(), w)
}

func SetEn(w http.ResponseWriter, r *http.Request) {
	globalLanguage = "en-US"
}
func SetEs(w http.ResponseWriter, r *http.Request) {
	// session, err1 := store.Get(r, "vezbamo.onrender.com-lang")
	// if err1 != nil {
	// 	fmt.Println("špnski greška get sessio")
	// 	// http.Error(w, err1.Error(), http.StatusInternalServerError)
	// 	// return
	// }

	// // Set some session values. sfdsalkčk
	// // session.Values["language"] = "spanski sinjor"

	// err2 := session.Save(r, w)
	// if err2 != nil {
	// 	fmt.Println("šplanski greška save sessio")
	// 	// http.Error(w, err2.Error(), http.StatusInternalServerError)
	// 	// return
	// }
	globalLanguage = "es"
}
func SetSr(w http.ResponseWriter, r *http.Request) {
	globalLanguage = "sr"
}
func SetBrowserLang(w http.ResponseWriter, r *http.Request) {
	globalLanguage = ""
}

//*** A P I

func GetLocationsForAngularFE(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nget locations", r.URL)
	dat, err := os.ReadFile("src/db/locations.json")
	check(err)

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, string(dat))
	//fmt.Println("\ndat: ", string(dat))
	//w.Write(string(dat)) dfaljfa
}

func APIgetQuestions(w http.ResponseWriter, r *http.Request) {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// curl http://127.0.0.1:7331/api_questions
	// io.WriteString(w, string(db.GetQuestions()))
	w.Write(db.GetQuestions())
}
