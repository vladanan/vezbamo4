package vezbamo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/vladanan/vezbamo4/src/clr"
	"github.com/vladanan/vezbamo4/src/db"
)

type TestHandler struct {
	db db.DB
}

func NewTestHandler(db db.DB) *TestHandler {
	return &TestHandler{db: db}
}

func (h *TestHandler) HandleGetTests(w http.ResponseWriter, r *http.Request) error {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// io.WriteString(w, string(db.GetQuestions()))
	// curl http://127.0.0.1:7331/api_get_tests

	l := clr.GetELRfunc2()

	vars := mux.Vars(r)
	record := vars["id"]
	var g_id int

	fmt.Println("id:", record, "broj:", g_id, r.Method, r.URL.Path)

	if record != "" {
		var err error
		g_id, err = strconv.Atoi(record)
		if err != nil {
			return l(r, 0, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax"))
		}
	}

	data, err := h.db.GetTests(g_id, r)
	if err != nil {
		return err
	}
	if data != nil {
		return clr.WriteJSON(w, 200, data)
	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	// w.Write(db.GetQuestions())
	// return db.GetQuestions()
}
