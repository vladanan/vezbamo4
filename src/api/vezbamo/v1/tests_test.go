package testsAPI

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vladanan/vezbamo4/src/errorlogres"
	"github.com/vladanan/vezbamo4/src/models"
)

func TestGetTests(t *testing.T) {
	// https://go.dev/doc/code
	// https://www.cloudbees.com/blog/testing-http-handlers-go

	r, err := http.NewRequest(http.MethodGet, "/api_get_tests", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(errorlogres.CheckFunc(GetTests))
	handler.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("api handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	list_string := rr.Body.String() // r.Body is a *bytes.Buffer
	dec := json.NewDecoder(strings.NewReader(list_string))
	var all_tests []models.Test
	if err := dec.Decode(&all_tests); err != nil {
		t.Errorf("json decode greška: %v", err)
	}

	if len(all_tests) < 1 {
		t.Errorf("api handler returned unexpected body: got %v, want %v", len(all_tests), ">1")

	}

	// cases := []struct {
	// 	in (http.Request, http.ResponseWriter)
	// 	want []models.Test
	// }{
	// }

}
