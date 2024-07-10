package testsAPI

import (
	"net/http"

	"github.com/vladanan/vezbamo4/src/db"
	"github.com/vladanan/vezbamo4/src/errorlogres"
)

// func to_struct(questions []byte) []models.Test {
// 	var p []models.Test
// 	err := json.Unmarshal(questions, &p)
// 	if err != nil {
// 		fmt.Printf("Json error: %v", err)
// 	}
// 	return p
// }

func GetTests(w http.ResponseWriter, r *http.Request) error {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// io.WriteString(w, string(db.GetQuestions()))
	// curl http://127.0.0.1:7331/api_get_tests

	allTests := db.GetTests()
	// l2 := to_struct(list)
	return errorlogres.WriteJSON(w, 200, allTests)

	// w.Write(db.GetQuestions())

	// return db.GetQuestions()
}
