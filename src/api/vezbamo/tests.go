package vezbamo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vladanan/vezbamo4/src/clr"
	"github.com/vladanan/vezbamo4/src/db"
)

func GetTests(w http.ResponseWriter, r *http.Request) error {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// io.WriteString(w, string(db.GetQuestions()))
	// curl http://127.0.0.1:7331/api_get_tests

	l := clr.GetELRfunc2()

	vars := mux.Vars(r)
	id := vars["id"]
	var g_id int

	switch id {
	case "":
	default:
		var err error
		g_id, err = strconv.Atoi(id)
		if err != nil {
			l(r, 4, err)
		}
	}

	fmt.Println("traži se test:", id, "broj:", g_id, r.Method)

	allTests, err := db.GetTests(r, g_id)
	if err != nil {
		return err
	}
	return clr.WriteJSON(w, 200, allTests)

	// w.Write(db.GetQuestions())
	// return db.GetQuestions()
}
