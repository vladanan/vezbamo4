package testovi

import (
	"net/http"

	"github.com/vladanan/vezbamo4/src/db"
	elr "github.com/vladanan/vezbamo4/src/errorlogres"
)

func GetTests(w http.ResponseWriter, r *http.Request) error {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// io.WriteString(w, string(db.GetQuestions()))
	// curl http://127.0.0.1:7331/api_get_tests

	allTests, err := db.GetTests(r)
	if err != nil {
		return err
	}
	return elr.WriteJSON(w, 200, allTests)

	// w.Write(db.GetQuestions())
	// return db.GetQuestions()
}
