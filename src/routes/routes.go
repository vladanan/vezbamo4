package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

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

//***  P A G E S

func GoToIndex(w http.ResponseWriter, r *http.Request) {
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

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func User(w http.ResponseWriter, r *http.Request) {
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

func Login(w http.ResponseWriter, r *http.Request) {
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
		User(w, r)
	} else {
		session.Values["authenticated"] = false
		session.Save(r, w)
		templ.Handler(views.Index(globalLanguage, r)).Component.Render(context.Background(), w)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
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
