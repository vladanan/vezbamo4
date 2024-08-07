// Package routes služi da obrađuje zahvete iz main
package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"

	"github.com/vladanan/vezbamo4/src/controllers/api"
	"github.com/vladanan/vezbamo4/src/controllers/clr"

	"github.com/vladanan/vezbamo4/src/models"
	"github.com/vladanan/vezbamo4/src/views"
	"github.com/vladanan/vezbamo4/src/views/assignments"
	"github.com/vladanan/vezbamo4/src/views/dashboard"
	"github.com/vladanan/vezbamo4/src/views/site"
	"github.com/vladanan/vezbamo4/src/views/tests"
)

// var store string = ""

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

// http://127.0.0.1:7331

func RouterSite(r *mux.Router) {
	r.HandleFunc("/", Index)
	r.HandleFunc("/assignments", Assignments)
	r.HandleFunc("/tests", Tests)
	r.HandleFunc("/user_portal", UserPortal)
	r.HandleFunc("/tests_api", TestsAPI)
	r.HandleFunc("/mega_increment", MegaIncrement)
	r.HandleFunc("/custom_apis", CustomAPIs)
	r.HandleFunc("/history", History)
	r.HandleFunc("/privacy", Privacy)
	r.HandleFunc("/terms", Terms)
	r.HandleFunc("/komponents", Komponents)
}

func RouterTests(r *mux.Router) {
	r.HandleFunc("/htmx_get_tests", HtmxGetTests)
}

func RouterAssignments(r *mux.Router) {
	r.HandleFunc("/primary_grade_1", PrimaryGrade1)
	r.HandleFunc("/primary_grade_2", PrimaryGrade2)
	r.HandleFunc("/secondary_grade_1", SecondaryGrade1)
}

func RouterUsers(r *mux.Router) {
	r.HandleFunc("/sign_up", Sign_up)
	r.HandleFunc("/sign_up_post", Sign_up_post)
	r.HandleFunc("/sign_in", Sign_in)
	r.HandleFunc("/sign_in_post", Sign_in_post)
	r.HandleFunc("/auto_login_user", AutoLoginUser)
	r.HandleFunc("/auto_login_admin", AutoLoginAdmin)
	r.HandleFunc("/dashboard", Dashboard)
	r.HandleFunc("/sign_out", Sign_out)
	// samo query koji ima u sebi tačno određene promenljive može da prođe
	r.HandleFunc("/vmk/{key}", CheckLinkFromEmail).Queries("email", "") // , "user", "vladan")
	// isto kao i ono gore:
	// vmk := r.PathPrefix("/vmk").Subrouter()
	// vmk.HandleFunc("/{key}", CheckLinkFromEmail).Queries("email", "")
	r.HandleFunc("/ext/verify_email.html", GetVerifyEmailHtml)
}

func RouterI18n(r *mux.Router) {
	r.HandleFunc("/sh", SetSh)
	r.HandleFunc("/en", SetEn)
	r.HandleFunc("/es", SetEs)
}

// static sa funkcijom koja pravi niz mux PathPrefix handlera jer bi se inače main zagušio sa vazdan njih
// za svaki folder gde se koristi path sa promeljivima kao što je r.HandleFunc("/vmk/{key}"
// https://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkits
func ServeStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"/":   "" + staticDirectory,
		"vmk": "/vmk" + staticDirectory,
		// "qapi": "/questions" + staticDirectory,
	}
	for _, pathValue := range staticPaths {
		// pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathValue).Handler(http.StripPrefix(pathValue, http.FileServer(http.Dir("assets"))))
	}
}

func apiCallGet(table string, field string, record string) (*json.Decoder, error) {

	switch {
	case table != "" && field != "" && record != "":
		table = "/" + table
		field = "/" + field
		record = "/" + record
	case table != "":
		table = "/" + table
	default:
		return nil, fmt.Errorf("no api path elements specified")
	}

	log.Println("api call:", table, field, record)

	var url string
	if os.Getenv("PRODUCTION") == "FALSE" {
		url = "http://127.0.0.1:7331/api/v" + table + field + record
	} else {
		url = "https://vezbamo.onrender.com/api/v"
	}

	var dec *json.Decoder

	resp, err := http.Get(url)
	if err != nil {
		// we will get an error at this stage if the request fails, such as if the
		// requested URL is not found, or if the server is not reachable.
		return nil, err
	} else {
		// if we want to check for a specific status code, we can do so here
		// for example, a successful request should return a 200 OK status
		if resp.StatusCode != http.StatusOK {
			// if the status code is not 200, we should log the status code and the
			// status string, then exit with a fatal error
			return nil, fmt.Errorf("apiCallGet url error: %v, url: %v", resp.Status, url)
		} else {
			// print the response
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			} else {
				dec = json.NewDecoder(strings.NewReader(string(data)))
			}
		}
	}
	defer resp.Body.Close()

	return dec, nil
}

////**** SITE

func Index(w http.ResponseWriter, r *http.Request) {
	// ovo mora da bude tu da bi store i ostalo radili oko os.Getenv("SESSION_KEY")
	if godotevnErr != nil {
		fmt.Printf("Error loading .env file")
	}
	templ.Handler(views.Index(store, r)).Component.Render(context.Background(), w)
}

func GoTo404(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.Page404()).Component.Render(context.Background(), w)
}

func Tests(w http.ResponseWriter, r *http.Request) {
	templ.Handler(tests.Tests(store, r)).Component.Render(context.Background(), w)
}

func UserPortal(w http.ResponseWriter, r *http.Request) {

	l := clr.GetELRfunc2()

	if dec, err := apiCallGet("note", "", ""); err != nil {
		templ.Handler(site.ServerError(clr.CheckErr(l(r, 8, err)))).Component.Render(context.Background(), w)
	} else {
		var notes []models.Note
		if err := dec.Decode(&notes); err != nil {
			templ.Handler(site.ServerError(clr.CheckErr(l(r, 8, err)))).Component.Render(context.Background(), w)
		} else {
			templ.Handler(site.UserPortal(store, r, notes)).Component.Render(context.Background(), w)
		}
	}

}

func TestsAPI(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.TestsAPI(store, r)).Component.Render(context.Background(), w)
}

func MegaIncrement(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.MegaIncrement(store, r)).Component.Render(context.Background(), w)
}

func CustomAPIs(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.CustomAPIs(store, r)).Component.Render(context.Background(), w)
}

func History(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.History(store, r)).Component.Render(context.Background(), w)
}

func Privacy(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.Privacy(store, r)).Component.Render(context.Background(), w)
}

func Terms(w http.ResponseWriter, r *http.Request) {
	templ.Handler(site.Terms(store, r)).Component.Render(context.Background(), w)
}

func Komponents(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(db.GetNotes())
	templ.Handler(views.Komponents()).Component.Render(context.Background(), w)
}

// //**** QUESTIONS

func HtmxGetTests(w http.ResponseWriter, r *http.Request) {

	l := clr.GetELRfunc2()

	if dec, err := apiCallGet("test", "", ""); err != nil {
		templ.Handler(site.ServerError(clr.CheckErr(l(r, 8, err)))).Component.Render(context.Background(), w)
	} else {
		var testsa []models.Test
		if err := dec.Decode(&testsa); err != nil {
			templ.Handler(site.ServerError(clr.CheckErr(l(r, 8, err)))).Component.Render(context.Background(), w)
		} else {
			templ.Handler(tests.List(testsa)).Component.Render(context.Background(), w)
		}
	}

	// https://stackoverflow.com/questions/13765797/the-best-way-to-get-a-string-from-a-writer
	// ch := vezbamo.NewTestHandler(db.DB{})
	// rr := httptest.NewRecorder()
	// err := ch.HandleGetAll(rr, r)
	// err := vezbamo.GetTests(rr, r)
	// if err != nil {
	// 	// log.Println("greška na api")
	// 	templ.Handler(site.ServerError(clr.CheckErr(err))).Component.Render(context.Background(), w)
	// } else {
	// 	list_string := rr.Body.String() // r.Body is a *bytes.Buffer
	// 	dec := json.NewDecoder(strings.NewReader(list_string))
	// 	var all_tests []models.Test
	// 	if err := dec.Decode(&all_tests); err != nil {
	// 		// log.Println("greška json dekodera")
	// 		templ.Handler(site.ServerError(clr.CheckErr(err))).Component.Render(context.Background(), w)
	// 	} else {
	// 		templ.Handler(tests.List(all_tests)).Component.Render(context.Background(), w)
	// 	}
	// }

}

// //**** ASSIGNMENTS
func Assignments(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.Assignments(store, r)).Component.Render(context.Background(), w)
}

func PrimaryGrade1(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.PrimaryGrade1(store, r)).Component.Render(context.Background(), w)
}

func PrimaryGrade2(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.PrimaryGrade2(store, r)).Component.Render(context.Background(), w)
}

func SecondaryGrade1(w http.ResponseWriter, r *http.Request) {
	templ.Handler(assignments.SecondaryGrade1(store, r)).Component.Render(context.Background(), w)
}

////**** USERS

func Sign_up(w http.ResponseWriter, r *http.Request) {
	templ.Handler(dashboard.Sign_up(store, r)).Component.Render(context.Background(), w)
}

func Sign_up_post(w http.ResponseWriter, r *http.Request) {
	email1 := r.FormValue("email1")
	email2 := r.FormValue("email2")
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

	user_map := session.Values["user_email"]
	user_email := ""
	if user_map != nil {
		user_email = user_map.(string)
	}

	if already_authenticated {

		_, data, _ := models.AuthenticateUser(user_email, "", already_authenticated, r)
		if user, ok := data.(models.User); ok {
			templ.Handler(dashboard.Dashboard(store, r, user)).Component.Render(context.Background(), w)
		} else {
			templ.Handler(dashboard.Dashboard(store, r, models.User{})).Component.Render(context.Background(), w)
		}

	} else {

		// validacija za UPIS NOVOG KORISNIKA a-zA-Z09 .,+-*:!?() min char 8 max 32 ISTO URADITI I NA FE UZ ARGUMENTS I JS
		// fmt.Println("SING UP POST form:", r.Form, len(r.Form) == 0)

		var validated bool

		if email1 != email2 || password1 != password2 {
			validated = false
			log.Println("Sign_up_post: validacija za upis korisnika nije prošla ISTI EMAIL/PASS")
		} else if len(r.Form) == 0 {
			validated = false
			log.Println("Sign_up_post: validacija za upis korisnika nije prošla PRAZAN FORM")
		} else {
			// NA DB PROVERITI DA LI VEĆ POSTOJI EMAIL I USER NAME i vratiti odgovarajuće poruke nazad osim bool za validated
			// NA DB PROVERITI da li je sa istog ip-a već bio upis u prethodnih 10min u odnosu na created_at
			validated = api.AddUser(email1, userName, password1, r)
			log.Println("Sign_up_post: validacija IZ DB:", validated)
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

	user_map := session.Values["user_email"]
	user_email := ""
	if user_map != nil {
		user_email = user_map.(string)
	}

	if already_authenticated {
		_, data, _ := models.AuthenticateUser(user_email, "", already_authenticated, r)
		if user, ok := data.(models.User); ok {
			templ.Handler(dashboard.Dashboard(store, r, user)).Component.Render(context.Background(), w)
		} else {
			templ.Handler(dashboard.Dashboard(store, r, models.User{})).Component.Render(context.Background(), w)
		}
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

	// fmt.Println("SING IN POST:", r.FormValue("email"), r.FormValue("password"), r.Body, r.MultipartForm, r.URL, r.PostForm, r.Form, r.FormValue("email"))
	fmt.Println("\nSIGN IN POST form:", r.FormValue("mail"), r.FormValue("password"))

	email := r.FormValue("email")
	password := r.FormValue("password")
	authenticated, data, err := models.AuthenticateUser(email, password, false, r)
	msg_fe := ""
	if err != nil {
		msg_fe = "Email_or_pass_wrong"
	} else {
		msg_fe = "Welcome"
	}

	// Set user as authenticated
	if authenticated {
		// fmt.Println("Sign_in_post: mail i user name", user.Email, user.User_name)
		session.Values["authenticated"] = true

		if user, ok := data.(models.User); ok {
			session.Values["user_email"] = user.Email
			session.Values["user_name"] = user.User_name
		} else {
			log.Println("greška pri dodeli mejla i user_name za sesiju")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
		fmt.Println("Sign_in_post: autentikacija korisnika NIJE prošla, msg_fe: ", msg_fe)
		session.Values["authenticated"] = false
		session.Values["user_email"] = "bbb"
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		templ.Handler(dashboard.UserNotLogedPage(store, r, msg_fe)).Component.Render(context.Background(), w)
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
	authenticated, data, err := models.AuthenticateUser(email, password, false, r)
	msg_fe := ""
	if err != nil {
		msg_fe = "Mail_or_pass_wrong"
	} else {
		msg_fe = "Unwelcome"
	}
	// Set user as authenticated
	if authenticated {
		session.Values["authenticated"] = true

		if user, ok := data.(models.User); ok {
			session.Values["user_email"] = user.Email
			session.Values["user_name"] = user.User_name
		} else {
			log.Println("greška pri dodeli mejla i user_name za sesiju")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
		session.Values["user_email"] = "ccc"
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templ.Handler(dashboard.UserNotLogedPage(store, r, msg_fe)).Component.Render(context.Background(), w)
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
	authenticated, data, err := models.AuthenticateUser(email, password, false, r)
	msg_fe := ""
	if err != nil {
		msg_fe = "Email_or_pass_wrong"
	} else {
		msg_fe = "Unwelcome"
	}
	// Set user as authenticated
	if authenticated {
		session.Values["authenticated"] = true

		if user, ok := data.(models.User); ok {
			session.Values["user_email"] = user.Email
			session.Values["user_name"] = user.User_name
		} else {
			log.Println("greška pri dodeli mejla i user_name za sesiju")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
		session.Values["user_email"] = "ccc"
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templ.Handler(dashboard.UserNotLogedPage(store, r, msg_fe)).Component.Render(context.Background(), w)
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

	user_map := session.Values["user_email"]
	user_email := ""
	if user_map == nil {
		// fmt.Println("nema email:", session.Values["user_email"])
	} else {
		// fmt.Println("ima email:", session.Values["user_email"])
		user_email = user_map.(string)
	}
	// Set user as authenticated
	if already_authenticated {
		_, data, _ := models.AuthenticateUser(user_email, "", already_authenticated, r)
		if user, ok := data.(models.User); ok {
			templ.Handler(dashboard.Dashboard(store, r, user)).Component.Render(context.Background(), w)
		} else {
			templ.Handler(dashboard.Dashboard(store, r, models.User{})).Component.Render(context.Background(), w)
		}
	} else {
		templ.Handler(dashboard.Dashboard(store, r, models.User{})).Component.Render(context.Background(), w)
	}
}

func Sign_out(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["user_email"] = nil
	session.Save(r, w)
	templ.Handler(views.Index(store, r)).Component.Render(context.Background(), w)
}

func CheckLinkFromEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values

	// deo iz query URL.Query i FormValue ne rade na isti način pogotovo ako u r ima body i multipart form
	fmt.Print("CheckLinkFromEmail: url vars and queries:", vars, r.URL.Query()["email"][0], r.FormValue("email"), "\n")

	// delovi patha-a tj. urla
	// title := vars["title"]
	key := vars["key"]
	email := r.URL.Query()["email"][0]

	// fmt.Println("ceo url query", r.URL.Query())

	emailVerified := models.AuthenticateEmail(key, email)

	if emailVerified {
		templ.Handler(dashboard.EmailVerified(store, r)).Component.Render(context.Background(), w)
		// fmt.Fprint(w, "Your mail is registered. You can go back to homepage and sign in")
	} else {
		// fmt.Fprintf(w, "You mail is NOT REGISTERED. Contact user support.")
		// fmt.Fprintf(w, "You want to register this key from mail bre: %s\n", key)
		templ.Handler(dashboard.EmailNotVerified(store, r)).Component.Render(context.Background(), w)
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

	user_email_map := session.Values["user_email"]
	user_email := ""
	if user_email_map == nil {
		// fmt.Println("nema mail:", session.Values["user_email"])
	} else {
		// fmt.Println("ima mail:", session.Values["user_email"])
		user_email = user_email_map.(string)
	}

	user_name_map := session.Values["user_name"]
	user_name := ""
	if user_name_map == nil {
		// fmt.Println("nema mail:", session.Values["user_email"])
	} else {
		// fmt.Println("ima mail:", session.Values["user_email"])
		user_name = user_name_map.(string)
	}

	// uzima se html fajl za mejl za verifikaciju
	dat, err1 := os.ReadFile("src/ext/verify_email.html")
	if err1 != nil {
		fmt.Printf("getVerifyEmailHtml: greška čitanje html fajla: %v\n", err1)
	}
	html := strings.ReplaceAll(string(dat), "+userName+", user_name)
	html = strings.ReplaceAll(html, "+urlDomainForEmail+", "http://127.0.0.1:7331/vmk/$2a$07$IFkFJy1NufwawNGqoef6kuJLuVFKzhqI4v_hYYwK2f_Y6Y3pP2eGu?email=y.emailbox-proba@yahoo.com")
	html = strings.ReplaceAll(html, "+emailForEmail+", user_email)

	w.Write([]byte(html))
}

////*** API

////**** i18n

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
