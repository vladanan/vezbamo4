package main

import (
	// "context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	// "github.com/a-h/templ"
	"github.com/rs/cors"
	"github.com/vladanan/vezbamo4/src/routes"

	// site "github.com/vladanan/vezbamo4/src/views/site"

	"github.com/gorilla/mux"
)

func main() {

	var dir string

	flag.StringVar(&dir, "dir", "assets", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	// R O U T E R  M U X
	r := mux.NewRouter()

	//***  P A G E S

	// http://localhost:10000
	// http://127.0.0.1:7331
	r.HandleFunc("/", routes.GoToIndex)

	//http.Handle("/questions", templ.Handler(questions.Questions()))
	r.HandleFunc("/questions", routes.GoToQuestions)
	r.HandleFunc("/questions_api", routes.GoToQuestionsAPI)

	r.HandleFunc("/assignments", routes.GoToAssignments)
	r.HandleFunc("/primary_grade_1", routes.GoToPrimaryGrade1)
	r.HandleFunc("/primary_grade_2", routes.GoToPrimaryGrade2)
	r.HandleFunc("/secondary_grade_1", routes.GoToSecondaryGrade1)

	r.HandleFunc("/user_portal", routes.GoToUserPortal)

	r.HandleFunc("/mega_increment", routes.GoToMegaIncrement)

	r.HandleFunc("/custom_apis", routes.GoToCustomAPIs)

	r.HandleFunc("/history", routes.GoToHistory)
	r.HandleFunc("/privacy", routes.GoToPrivacy)
	r.HandleFunc("/terms", routes.GoToTerms)
	// http.HandleFunc("/da", routes.GoToDa)

	// http.Handle("/404", http.NotFoundHandler())
	// r.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
	// 	templ.Handler(site.Page404()).Component.Render(context.Background(), w)
	// })
	r.NotFoundHandler = http.HandlerFunc(routes.GoTo404)

	r.HandleFunc("/auto_login", routes.AutoLogin)
	r.HandleFunc("/sign_in", routes.Sign_in)
	r.HandleFunc("/sign_in_post", routes.Sign_in_post)
	r.HandleFunc("/dashboard", routes.Dashboard)
	r.HandleFunc("/sign_out", routes.Sign_out)
	r.HandleFunc("/sign_up", routes.Sign_up)
	r.HandleFunc("/sign_up_post", routes.Sign_up_post)

	r.HandleFunc("/komponents", routes.GoToKomponents)

	//***  I N T E R N A L S

	// fs := http.FileServer(http.Dir("assets/"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.PathPrefix("/vmk/static/").Handler(http.StripPrefix("/vmk/static/", http.FileServer(http.Dir(dir))))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	r.HandleFunc("/htmx_get_questions", routes.HtmxGetQuestions)

	r.HandleFunc("/en", routes.SetEn)
	r.HandleFunc("/es", routes.SetEs)
	r.HandleFunc("/sh", routes.SetSh)
	// http.HandleFunc("/browser", routes.SetBrowserLang)

	// r.Path("/vmk/{key}").Queries("mail", "a").HandlerFunc(routes.CheckLinkFromEmail)
	// r.NewRoute().Path("/vmk/{key}").HandlerFunc(routes.CheckLinkFromEmail).Queries("mail", "a")
	// r.NewRoute().Path("/vmk/{key}").Handler(http.HandlerFunc(routes.CheckLinkFromEmail)).Queries("mail", "a")
	// https://chromium.googlesource.com/external/github.com/gorilla/mux/+/refs/tags/v1.7.4/mux_test.go
	// https://stackoverflow.com/questions/64279203/gorilla-mux-router-get-url-with-multiple-query-param-only-matches-if-called-by-a
	// https://pkg.go.dev/github.com/gorilla/mux@v1.8.1#Route.Queries
	// Queries(key, value) pairs: mail "" može bilo šta mail=bla@bla.com, "user" "vladan" može samo user=vladan
	// takođe "mail" "{mail}" isto može bilo šta (nije mi jasna razlika između "" i "{mail}")
	// tako samo query koji ima u sebi tačno određene promenljive sa određenim vrednostima (regex: "id", "{id:[0-9]+}") može da prođe
	// to otežava provaljivanje i povećava bezbednost, može da ima i više promenljivih u query nego što je definisano ali mora da ima one koje su definisane
	r.HandleFunc("/vmk/{key}", routes.CheckLinkFromEmail).Queries("mail", "") //, "user", "vladan")

	r.HandleFunc("/html/verify_email.html", routes.GetVerifyEmailHtml)

	//***  A P I

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	r.Handle("/locations", c.Handler(http.HandlerFunc(routes.GetLocationsForAngularFE)))
	r.Handle("/locations/", c.Handler(http.HandlerFunc(routes.GetLocationsForAngularFE)))

	r.HandleFunc("/api_questions", routes.APIgetQuestions)

	//***  S E R V E R

	fmt.Println("Main done", time.Now().Second())

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:10000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// log.Fatal(srv.ListenAndServe())

	// var err = srv.ListenAndServe("0.0.0.0:10000", r)
	var err = srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
