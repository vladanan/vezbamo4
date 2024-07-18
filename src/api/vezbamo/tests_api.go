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

	if id != "" {
		var err error
		g_id, err = strconv.Atoi(id)
		if err != nil {
			return l(r, 0, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax"))
		}
	}

	fmt.Println("id:", id, "broj:", g_id, r.Method)

	data, err := db.GetTests(r, g_id)
	if err != nil {
		return err
	}
	if data != nil {
		return clr.WriteJSON(w, 200, data)
	} else {
		return clr.NewAPIError(406, "no (available) test with requested id")
	}

	// w.Write(db.GetQuestions())
	// return db.GetQuestions()
}
