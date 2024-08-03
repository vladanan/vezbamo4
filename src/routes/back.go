// Package routes služi da obrađuje zahvete iz main
package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/vladanan/vezbamo4/src/api/custom/eone"
	"github.com/vladanan/vezbamo4/src/api/vezbamo"
	"github.com/vladanan/vezbamo4/src/clr"

	"github.com/vladanan/vezbamo4/src/models"
)

// var store string = ""

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

func RouterAPI(r *mux.Router) {
	vh := vezbamo.NewVezbamoHandler(models.DB{})
	r.HandleFunc("/api/{table}", clr.CheckFunc(vh.HandleGetMany)).Methods("GET")
	r.HandleFunc("/api/{table}/{field}/{record}", clr.CheckFunc(vh.HandleGetOne)).Methods("GET")

	ch := eone.NewEoneHandler(models.DB{})
	r.HandleFunc("/api/c/billing", clr.CheckFunc(ch.HandleGetBilling))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	r.Handle("/locations", c.Handler(http.HandlerFunc(GetLocationsForAngularFE)))
	r.Handle("/locations/", c.Handler(http.HandlerFunc(GetLocationsForAngularFE)))
}

////*** API

func GetLocationsForAngularFE(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nget locations", r.URL)
	dat, err := os.ReadFile("src/db/locations.json")
	if err != nil {
		log.Println("ne može da se pročita fajl za locations")
	}
	// check(err)

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, string(dat))
	// fmt.Println("\ndat: ", string(dat))
	// w.Write(string(dat)) dfaljfa
}
