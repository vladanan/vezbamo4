package routes

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/joho/godotenv"
	"github.com/vladanan/vezbamo4/src/db"
	"github.com/vladanan/vezbamo4/src/models"
	views "github.com/vladanan/vezbamo4/src/views"
	assignments "github.com/vladanan/vezbamo4/src/views/assignments"
	dashboard "github.com/vladanan/vezbamo4/src/views/dashboard"
	questions "github.com/vladanan/vezbamo4/src/views/questions"
	site "github.com/vladanan/vezbamo4/src/views/site"
	// "encoding/json"
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
	// ovo mora da bude tu da bi store i ostalo radili oko os.Getenv("SESSION_KEY")
	if godotevn_err != nil {
		fmt.Printf("Error loading .env file")
	}
	templ.Handler(views.Index(store, r)).Component.Render(context.Background(), w)
}

func GoTo404(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.Page404()).Component.Render(context.Background(), w)
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

// func GoToDa(w http.ResponseWriter, r *http.Request) {
// 	templ.Handler(da.Da(store, r)).Component.Render(context.Background(), w)
// }

func Sign_up(w http.ResponseWriter, r *http.Request) {
	templ.Handler(dashboard.Sign_up(store, r)).Component.Render(context.Background(), w)
}

func Sign_up_post(w http.ResponseWriter, r *http.Request) {
	email1 := r.FormValue("mail1")
	email2 := r.FormValue("mail2")
	userName := r.FormValue("user_name")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	// PROVERA DA LI JE KORISNIK VEĆ PRIJAVLJEN:

	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var already_authenticated bool

	auth_map := session.Values["authenticated"]
	if auth_map != nil {
		already_authenticated = auth_map.(bool)
	}

	user_map := session.Values["user_mail"]
	user_mail := ""
	if user_map != nil {
		user_mail = user_map.(string)
	}

	if already_authenticated {

		_, data, _ := db.AuthenticateUser(user_mail, "", already_authenticated, r)
		templ.Handler(dashboard.Dashboard(store, r, data)).Component.Render(context.Background(), w)

	} else {

		// validacija za UPIS NOVOG KORISNIKA a-zA-Z09 .,+-*:!?() min char 8 max 32 ISTO URADITI I NA FE UZ ARGUMENTS I JS
		// fmt.Println("SING UP POST form:", r.Form, len(r.Form) == 0)

		var validated bool

		if email1 != email2 || password1 != password2 {
			validated = false
			fmt.Print("Sign_up_post: validacija za upis korisnika nije prošla ISTI MEJL/PASS\n")
		} else if len(r.Form) == 0 {
			validated = false
			fmt.Print("Sign_up_post: validacija za upis korisnika nije prošla PRAZAN FORM\n")
		} else {
			// NA DB PROVERITI DA LI VEĆ POSTOJI MAIL I USER NAME i vratiti odgovarajuće poruke nazad osim bool za validated
			// NA DB PROVERITI da li je sa istog ip-a već bio upis u prethodnih 10min u odnosu na created_at
			validated = db.AddUser(email1, userName, password1, r)
			fmt.Print("Sign_up_post: validacija IZ DB:", validated, "\n")
		}

		if validated {
			templ.Handler(dashboard.UserRegistered(store, r)).Component.Render(context.Background(), w)
		} else {
			templ.Handler(dashboard.UserNotRegistered(store, r)).Component.Render(context.Background(), w)
		}

	}
}

func Sign_in(w http.ResponseWriter, r *http.Request) {
	// bytearray_headers, err2 := json.Marshal(r.Header)
	// if err2 != nil {
	// 	fmt.Printf("Sign_in: JSON error: %v", err2)
	// }

	// fmt.Print("\nSign_in: header:", string(bytearray_headers), "\n")
	// for item, index := range r.Header {
	// 	fmt.Print("\nSign_in: header:", item, index, "\n")
	// }

	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var already_authenticated bool

	auth_map := session.Values["authenticated"]
	if auth_map != nil {
		already_authenticated = auth_map.(bool)
	}

	user_map := session.Values["user_mail"]
	user_mail := ""
	if user_map != nil {
		user_mail = user_map.(string)
	}

	if already_authenticated {
		_, data, _ := db.AuthenticateUser(user_mail, "", already_authenticated, r)
		templ.Handler(dashboard.Dashboard(store, r, data)).Component.Render(context.Background(), w)
	} else {
		templ.Handler(dashboard.Sign_in(store, r)).Component.Render(context.Background(), w)
	}
}

func Sign_in_post(w http.ResponseWriter, r *http.Request) {
	// log.SetFlags(log.Ltime | log.Lshortfile)

	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Authentication goes here

	// PROVERA ZA:
	// - user_name: samo a-zA-Z09 .,+-*:!?() min char 8
	// - password: isto, min char 8
	// ISTO URADITI I NA FE UZ ARGUMENTS I JS

	// fmt.Println("SING IN POST:", r.FormValue("mail"), r.FormValue("password"), r.Body, r.MultipartForm, r.URL, r.PostForm, r.Form, r.FormValue("mail"))
	fmt.Println("\nSIGN IN POST form:", r.FormValue("mail"), r.FormValue("password"))

	email := r.FormValue("mail")
	password := r.FormValue("password")
	authenticated, user, msg_fe := db.AuthenticateUser(email, password, false, r)

	// Set user as authenticated
	if authenticated {
		// fmt.Println("Sign_in_post: mail i user name", user.Email, user.User_name)
		session.Values["authenticated"] = true
		session.Values["user_mail"] = user.Email
		session.Values["user_name"] = user.User_name
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("Sign_in_post: autentikacija JE PROŠLA")
		// templ.Handler(dashboard.Dashboard(store, r, data)).Component.Render(context.Background(), w)
		Dashboard(w, r)
	} else {
		fmt.Println("Sign_in_post: autentikacija korisnika NIJE prošla, msg: ", msg_fe)
		session.Values["authenticated"] = false
		session.Values["user_mail"] = "bbb"
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		templ.Handler(dashboard.UserNotLogedPage(store, r)).Component.Render(context.Background(), w)
	}
}

func AutoLoginUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Authentication goes here
	email := "vladan_zasve@yahoo.com"
	password := "b"
	// email := "vladan.andjelkovic@gmail.com"
	// password := "vezbamo.2015"
	authenticated, user, _ := db.AuthenticateUser(email, password, false, r)
	// Set user as authenticated
	if authenticated {
		session.Values["authenticated"] = true
		session.Values["user_mail"] = user.Email
		session.Values["user_name"] = user.User_name
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Dashboard(w, r)
	} else {
		fmt.Println("AutoLogin: autentikacija admina nije prošla")
		session.Values["authenticated"] = false
		session.Values["user_mail"] = "ccc"
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templ.Handler(dashboard.UserNotLogedPage(store, r)).Component.Render(context.Background(), w)
	}
}

func AutoLoginAdmin(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Authentication goes here
	// email := "vladan_zasve@yahoo.com"
	// password := "b"
	email := "vladan.andjelkovic@gmail.com"
	password := "vezbamo.2015"
	authenticated, user, _ := db.AuthenticateUser(email, password, false, r)
	// Set user as authenticated
	if authenticated {
		session.Values["authenticated"] = true
		session.Values["user_mail"] = user.Email
		session.Values["user_name"] = user.User_name
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Dashboard(w, r)
	} else {
		fmt.Println("AutoLogin: autentikacija admina nije prošla")
		session.Values["authenticated"] = false
		session.Values["user_mail"] = "ccc"
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templ.Handler(dashboard.UserNotLogedPage(store, r)).Component.Render(context.Background(), w)
	}
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var already_authenticated bool
	auth_map := session.Values["authenticated"]

	if auth_map == nil {
		// session.Values["authenticated"] = false
		// fmt.Println("nema auth:", session.Values["authenticated"])
	} else {
		// already_authenticated = true
		// fmt.Println("ima auth sesion:", session.Values["authenticated"])
		// fmt.Println("ima auth map:", auth_map)
		already_authenticated = auth_map.(bool)
		// fmt.Println("ima auth map2:", already_authenticated)
	}

	user_map := session.Values["user_mail"]
	user_mail := ""
	if user_map == nil {
		// fmt.Println("nema mail:", session.Values["user_mail"])
	} else {
		// fmt.Println("ima mail:", session.Values["user_mail"])
		user_mail = user_map.(string)
	}
	// Set user as authenticated
	if already_authenticated {
		_, data, _ := db.AuthenticateUser(user_mail, "", already_authenticated, r)
		templ.Handler(dashboard.Dashboard(store, r, data)).Component.Render(context.Background(), w)
	} else {
		templ.Handler(dashboard.Dashboard(store, r, models.User{})).Component.Render(context.Background(), w)
	}
}

func Sign_out(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["user_mail"] = nil
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

func CheckLinkFromEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values

	// deo iz query URL.Query i FormValue ne rade na isti način pogotovo ako u r ima body i multipart form
	fmt.Print("CheckLinkFromEmail: url vars and queries:", vars, r.URL.Query()["mail"][0], r.FormValue("mail"), "\n")

	// delovi patha-a tj. urla
	// title := vars["title"]
	key := vars["key"]
	email := r.URL.Query()["mail"][0]

	mailVerified := db.AuthenticateMail(key, email)

	if mailVerified {
		templ.Handler(dashboard.MailVerified(store, r)).Component.Render(context.Background(), w)
		// fmt.Fprint(w, "Your mail is registered. You can go back to homepage and sign in")
	} else {
		// fmt.Fprintf(w, "You mail is NOT REGISTERED. Contact user support.")
		// fmt.Fprintf(w, "You want to register this key from mail bre: %s\n", key)
		templ.Handler(dashboard.MailNotVerified(store, r)).Component.Render(context.Background(), w)
		// GoToNV(w, r)
		// GoToTerms(w, r)
	}
	// fmt.Print("vmk prošao")
	// GoToTerms(w, r)
}

func GetVerifyEmailHtml(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user_mail_map := session.Values["user_mail"]
	user_mail := ""
	if user_mail_map == nil {
		// fmt.Println("nema mail:", session.Values["user_mail"])
	} else {
		// fmt.Println("ima mail:", session.Values["user_mail"])
		user_mail = user_mail_map.(string)
	}

	user_name_map := session.Values["user_name"]
	user_name := ""
	if user_name_map == nil {
		// fmt.Println("nema mail:", session.Values["user_mail"])
	} else {
		// fmt.Println("ima mail:", session.Values["user_mail"])
		user_name = user_name_map.(string)
	}

	// uzima se html fajl za mejl za verifikaciju
	dat, err1 := os.ReadFile("src/html/verify_email.html")
	if err1 != nil {
		fmt.Printf("getVerifyEmailHtml: greška čitanje html fajla: %v\n", err1)
	}
	html := strings.ReplaceAll(string(dat), "+user_name+", user_name)
	html = strings.ReplaceAll(html, "+url_domain_for_mail+", "http://127.0.0.1:7331/vmk/$2a$07$IFkFJy1NufwawNGqoef6kuJLuVFKzhqI4v_hYYwK2f_Y6Y3pP2eGu?mail=y.emailbox-proba@yahoo.com")
	html = strings.ReplaceAll(html, "+mail_for_mail+", user_mail)

	w.Write([]byte(html))
}

//*** A P I

func GetLocationsForAngularFE(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nget locations", r.URL)
	dat, err := os.ReadFile("src/db/locations.json")
	check(err)

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, string(dat))
	// fmt.Println("\ndat: ", string(dat))
	// w.Write(string(dat)) dfaljfa
}

func APIgetQuestions(w http.ResponseWriter, r *http.Request) {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// curl http://127.0.0.1:7331/api_questions
	// io.WriteString(w, string(db.GetQuestions()))
	w.Write(db.GetQuestions())
}
