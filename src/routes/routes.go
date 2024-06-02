package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/vladanan/vezbamo4/src/db"
	views "github.com/vladanan/vezbamo4/src/views"
	assignments "github.com/vladanan/vezbamo4/src/views/assignments"
	da "github.com/vladanan/vezbamo4/src/views/da"
	questions "github.com/vladanan/vezbamo4/src/views/questions"
	site "github.com/vladanan/vezbamo4/src/views/site"
)

// var store string = ""

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
	if godotevn_err != nil {
		fmt.Printf("Error loading .env file")
	}

	if r.URL.Path == "/" {
		templ.Handler(views.Index(store, r)).Component.Render(context.Background(), w)
	} else {
		templ.Handler(site.Page404()).Component.Render(context.Background(), w)
	}
}

func GoToQuestions(w http.ResponseWriter, r *http.Request) {
	templ.Handler(questions.Questions(store, r)).Component.Render(context.Background(), w)
}

func GoToQuestionsAPI(w http.ResponseWriter, r *http.Request) {
	templ.Handler(questions.QuestionsAPI(store, r)).Component.Render(context.Background(), w)
}

func GoToAssignments(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.Assignments(store, r)).Component.Render(context.Background(), w)
}
func GoToPrimaryGrade1(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.PrimaryGrade1(store, r)).Component.Render(context.Background(), w)
}
func GoToPrimaryGrade2(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.PrimaryGrade2(store, r)).Component.Render(context.Background(), w)
}
func GoToSecondaryGrade1(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.SecondaryGrade1(store, r)).Component.Render(context.Background(), w)
}

func GoToUserPortal(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(db.GetNotes())
	templ.Handler(site.UserPortal(store, r, db.GetNotes())).Component.Render(context.Background(), w)
}

func GoToMegaIncrement(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.MegaIncrement(store, r)).Component.Render(context.Background(), w)
}

func GoToCustomAPIs(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.CustomAPIs(store, r)).Component.Render(context.Background(), w)
}

func GoToHistory(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.History(store, r)).Component.Render(context.Background(), w)
}

func GoToPrivacy(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.Privacy(store, r)).Component.Render(context.Background(), w)
}

func GoToTerms(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.Terms(store, r)).Component.Render(context.Background(), w)
}

func GoToDa(w http.ResponseWriter, r *http.Request) {
	templ.Handler(da.Da(store, r)).Component.Render(context.Background(), w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Authentication goes here
	email := "vladan.andjelkovic@gmail.com"
	password := "vezbamo.2015"
	authenticated := db.AuthenticateUser(email, password)
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
		templ.Handler(views.Index(store, r)).Component.Render(context.Background(), w)
	}
}

func Admin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")
	// Check if user is authenticated
	// auth2, ok2 := session.Values["authenticated"].(bool)
	// fmt.Println("pristup admin sajtu: auth:", auth2, "ok:", ok2)
	// ulogovan: pristup admin sajtu: auth: true ok: true
	// neulogovan: pristup admin sajtu: auth: false ok: true
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		templ.Handler(site.UserNotLogedPage(store, r)).Component.Render(context.Background(), w)
		return
	}
	templ.Handler(site.Admin(store, r)).Component.Render(context.Background(), w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	templ.Handler(views.Index(store, r)).Component.Render(context.Background(), w)
}

func GoToKomponents(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(db.GetNotes())
	templ.Handler(views.Komponents()).Component.Render(context.Background(), w)
}

//*** I N T E R N A L

func HtmxGetQuestions(w http.ResponseWriter, r *http.Request) {
	list := db.GetQuestions()
	templ.Handler(questions.List(list)).Component.Render(context.Background(), w)
}

func SetEn(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	// fmt.Println("engleski podešavanje")
	if err != nil {
		// fmt.Println("engleski greška get sessio")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["language"] = "en"
	err2 := session.Save(r, w)
	if err2 != nil {
		// fmt.Println("engleski greška save sessio")
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}

func SetEs(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("španski podešavanje")
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		// fmt.Println("špnski greška get sessio")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["language"] = "es"
	err2 := session.Save(r, w)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}

func SetSh(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		// fmt.Println("srpski greška get sessio")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["language"] = "sh"
	err2 := session.Save(r, w)
	if err2 != nil {
		// fmt.Println("srpski greška save sessio")
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}

// func SetBrowserLang(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "vezbamo.onrender.com-users")
// 	if err != nil {
// 		// fmt.Println("browser greška get sessio")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	session.Values["language"] = ""
// 	err2 := session.Save(r, w)
// 	if err2 != nil {
// 		// fmt.Println("brower greška save sessio")
// 		http.Error(w, err2.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

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
