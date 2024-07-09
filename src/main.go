package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vladanan/vezbamo4/src/errorlogres"
	"github.com/vladanan/vezbamo4/src/routes"
)

func main() {

	////**** ROUTER MUX
	r := mux.NewRouter()

	////**** SITE

	r.HandleFunc("/", routes.Index)
	r.NotFoundHandler = http.HandlerFunc(routes.GoTo404)
	r.HandleFunc("/user_portal", routes.UserPortal)
	r.HandleFunc("/mega_increment", routes.MegaIncrement)
	r.HandleFunc("/custom_apis", routes.CustomAPIs)
	r.HandleFunc("/history", routes.History)
	r.HandleFunc("/privacy", routes.Privacy)
	r.HandleFunc("/terms", routes.Terms)
	r.HandleFunc("/komponents", routes.Komponents)

	////**** QUESTIONS
	r.HandleFunc("/questions", routes.Questions)
	r.HandleFunc("/questions_api", routes.QuestionsAPI)

	////**** ASSIGNMENTS
	r.HandleFunc("/assignments", routes.Assignments)
	r.HandleFunc("/primary_grade_1", routes.PrimaryGrade1)
	r.HandleFunc("/primary_grade_2", routes.PrimaryGrade2)
	r.HandleFunc("/secondary_grade_1", routes.SecondaryGrade1)

	////**** USERS
	r.HandleFunc("/sign_up", routes.Sign_up)
	r.HandleFunc("/sign_up_post", routes.Sign_up_post)
	r.HandleFunc("/sign_in", routes.Sign_in)
	r.HandleFunc("/sign_in_post", routes.Sign_in_post)
	r.HandleFunc("/auto_login_user", routes.AutoLoginUser)
	r.HandleFunc("/auto_login_admin", routes.AutoLoginAdmin)
	r.HandleFunc("/dashboard", routes.Dashboard)
	r.HandleFunc("/sign_out", routes.Sign_out)
	// samo query koji ima u sebi tačno određene promenljive može da prođe
	r.HandleFunc("/vmk/{key}", routes.CheckLinkFromEmail).Queries("mail", "") //, "user", "vladan")
	r.HandleFunc("/html/verify_email.html", routes.GetVerifyEmailHtml)

	////**** API

	r.HandleFunc("/api_get_questions", errorlogres.Check(routes.APIgetQuestions))
	r.HandleFunc("/htmx_get_questions", routes.HtmxGetQuestions)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	r.Handle("/locations", c.Handler(http.HandlerFunc(routes.GetLocationsForAngularFE)))
	r.Handle("/locations/", c.Handler(http.HandlerFunc(routes.GetLocationsForAngularFE)))

	////**** SERVER

	r.HandleFunc("/sh", routes.SetSh)
	r.HandleFunc("/en", routes.SetEn)
	r.HandleFunc("/es", routes.SetEs)

	// static sa funkcijom koja pravi niz mux PathPrefix handlera jer bi se inače main zagušio sa vazdan njih
	// za svaki folder gde se koristi path sa promeljivima kao što je r.HandleFunc("/vmk/{key}"
	routes.ServeStatic(r, "/static/")

	fmt.Println("Main done", time.Now().Local().Format(time.TimeOnly))

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:10000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// var err = srv.ListenAndServe("0.0.0.0:10000", r)
	var err = srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
