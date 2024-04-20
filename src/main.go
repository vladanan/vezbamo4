package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/rs/cors"
	"github.com/vladanan/vezbamo4/src/routes"
	site "github.com/vladanan/vezbamo4/src/views/site"
)

func main() {

	//***  P A G E S

	// http://localhost:10000
	// http://127.0.0.1:7331
	http.HandleFunc("/", routes.GoToIndex)

	//http.Handle("/questions", templ.Handler(questions.Questions()))
	http.HandleFunc("/questions", routes.GoToQuestions)
	http.HandleFunc("/questions_api", routes.GoToQuestionsAPI)
	http.HandleFunc("/assignments", routes.GoToAssignments)

	http.HandleFunc("/user_portal", routes.GoToUserPortal)

	http.HandleFunc("/mega_increment", routes.GoToMegaIncrement)

	http.HandleFunc("/custom_apis", routes.GoToCustomAPIs)

	http.HandleFunc("/history", routes.GoToHistory)
	http.HandleFunc("/privacy", routes.GoToPrivacy)
	http.HandleFunc("/terms", routes.GoToTerms)

	// http.Handle("/404", http.NotFoundHandler())
	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(site.Page404()).Component.Render(context.Background(), w)
	})

	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/admin", routes.Admin)
	http.HandleFunc("/logout", routes.Logout)

	http.HandleFunc("/komponents", routes.GoToKomponents)

	//***  I N T E R N A L S

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/htmx_get_questions", routes.HtmxGetQuestions)

	http.HandleFunc("/en", routes.SetEn)
	http.HandleFunc("/es", routes.SetEs)
	http.HandleFunc("/sh", routes.SetSh)
	// http.HandleFunc("/browser", routes.SetBrowserLang)

	//***  A P I

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	http.Handle("/locations", c.Handler(http.HandlerFunc(routes.GetLocationsForAngularFE)))
	http.Handle("/locations/", c.Handler(http.HandlerFunc(routes.GetLocationsForAngularFE)))

	http.HandleFunc("/api_questions", routes.APIgetQuestions)

	//***  S E R V E R

	fmt.Println("Main done", time.Now().Second())

	var err = http.ListenAndServe("0.0.0.0:10000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
