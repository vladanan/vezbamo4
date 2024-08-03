package vezbamo

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/vladanan/vezbamo4/src/clr"
	"github.com/vladanan/vezbamo4/src/models"
)

type VezbamoHandler struct {
	db models.DB
}

func NewVezbamoHandler(db models.DB) *VezbamoHandler {
	return &VezbamoHandler{db: db}
}

type Tim struct {
	table string
	id    string
	mail  string
}

func (h *VezbamoHandler) HandleGetOne(w http.ResponseWriter, r *http.Request) error {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// io.WriteString(w, string(db.GetQuestions()))
	// curl http://127.0.0.1:7331/api_get_tests

	l := clr.GetELRfunc2()

	vars := mux.Vars(r)

	// OVO NE RADI KADA SE IDE SA SAJTA NEGO SE VIDI DA JE PATH htmx_get_tests
	// zato i treba sve da ide preko http.Get() da bi sve išlo stvarno preko http api poziva
	// fmt.Println("url path ceo:", r.URL.Path, vars)
	tableApi := vars["table"]
	fieldApi := vars["field"]
	recordApi := vars["record"]

	// tableApi := strings.Split(r.URL.Path, "/")[2]
	// fieldApi := strings.Split(r.URL.Path, "/")[3]

	// fieldApi := fmt.Sprint(vars)
	// fieldApi = strings.ReplaceAll(fieldApi, "map[", "")
	// fieldApi = strings.Split(fieldApi, ":")[0]

	// http://127.0.0.1:7331/api/note/mail/n@n.com

	var tableDb, fieldDb string
	var recordDb any

	apiToDb2 := map[string]Tim{
		"test": {
			table: "g_pitanja_c_testovi",
			id:    "g_id",
			mail:  "user_id",
		},
		"user": {
			table: "mi_users",
			id:    "u_id",
			mail:  "email",
		},
		"note": {
			table: "g_user_blog",
			id:    "b_id",
			mail:  "user_mail",
		},
		"setting": {
			table: "v_settings",
			id:    "s_id",
			mail:  "",
		},
	}

	for a := range apiToDb2 {
		// fmt.Println("deo od map:", a)
		if a == tableApi {
			tableDb = apiToDb2[a].table
			switch fieldApi {
			case "id":
				fieldDb = apiToDb2[a].id
				// recordApi = vars[fieldApi]
			case "mail":
				fieldDb = apiToDb2[a].mail
				// recordApi = vars[fieldApi]
			}
		}
	}
	if tableDb == "" || fieldDb == "" {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	// fmt.Println("iz apija:", tableApi, fieldApi)

	//**********************************************

	// db.GetLocal(r)

	//**********************************************

	// // url := "http://127.0.0.1:1401/send"
	// // url := "http://127.0.0.1:8080/secure/balance"
	// url := "http://127.0.0.1:8080/secure/send"

	// resp, err := http.Get(url)
	// if err != nil {
	// 	// we will get an error at this stage if the request fails, such as if the
	// 	// requested URL is not found, or if the server is not reachable.
	// 	// log.Println("spoljni api greška", err)

	// 	return l(r, 8, err)
	// } else {
	// 	// if we want to check for a specific status code, we can do so here
	// 	// for example, a successful request should return a 200 OK status
	// 	if resp.StatusCode != http.StatusOK {
	// 		// if the status code is not 200, we should log the status code and the
	// 		// status string, then exit with a fatal error
	// 		l(r, 8, fmt.Errorf("http jasmin res: %v, url: %v", resp.Status, url))
	// 		// log.Println("http jasmin api status code:", resp.StatusCode, resp.Status) //
	// 	}

	// 	// print the response
	// 	data1, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Println("io greška:", err)
	// 	}
	// 	log.Println("tetka jasmin kaže:", string(data1))

	// }
	// defer resp.Body.Close()

	//**********************************************

	// fmt.Println("id:", record, "broj:", g_id, r.Method, r.URL.Path)

	if recordApi != "" && fieldApi == "id" {
		var err error
		recordDb, err = strconv.Atoi(recordApi)
		if err != nil {
			return l(r, 0, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax 0"))
		}
	}

	if recordApi != "" && fieldApi == "mail" {
		if m := strings.ContainsAny(recordApi, "@."); !m {
			return l(r, 0, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax 1"))
		}
		// napraviti funkciju za validaciju i sanitaciju za mejl itd.
		if m := strings.ContainsAny(recordApi, ",:;()[]<>{}/\\"); m {
			return l(r, 0, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax 2"))
		}
		recordDb = recordApi
	}

	// fmt.Println("za db podaci:", tableDb, fieldDb, recordApi, recordDb)

	data, err := h.db.GetOne(tableDb, fieldDb, recordDb, r)
	if err != nil {
		return err
	}
	// fmt.Println("api data:", data)
	if data != nil {
		return clr.WriteJSON(w, 200, data)
	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	// w.Write(db.GetQuestions())
	// return db.GetQuestions()
}

func (h *VezbamoHandler) HandleGetMany(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)

	// OVO NE RADI KADA SE IDE SA SAJTA NEGO SE VIDI DA JE PATH htmx_get_tests
	// zato i treba sve da ide preko http.Get() da bi sve išlo stvarno preko http api poziva
	// fmt.Println("url path ceo:", r.URL.Path, vars)
	tableApi := vars["table"]

	var tableDb string

	apiToDb2 := map[string]Tim{
		"test": {
			table: "g_pitanja_c_testovi",
			id:    "g_id",
			mail:  "user_id",
		},
		"user": {
			table: "mi_users",
			id:    "u_id",
			mail:  "email",
		},
		"note": {
			table: "g_user_blog",
			id:    "b_id",
			mail:  "user_mail",
		},
		"setting": {
			table: "v_settings",
			id:    "s_id",
			mail:  "",
		},
	}

	for a := range apiToDb2 {
		// fmt.Println("deo od map:", a)
		if a == tableApi {
			tableDb = apiToDb2[a].table
		}
	}
	if tableDb == "" {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	data, err := h.db.GetMany(tableDb, r)
	if err != nil {
		return err
	}
	// fmt.Println("api data:", data)
	if data != nil {
		return clr.WriteJSON(w, 200, data)
	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

}
