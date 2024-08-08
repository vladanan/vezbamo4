package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/vladanan/vezbamo4/src/controllers/clr"
	"github.com/vladanan/vezbamo4/src/controllers/vet"
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

var apiToDb = map[string]Tim{
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

func (h *VezbamoHandler) HandlePostOne(w http.ResponseWriter, r *http.Request) error {

	l := clr.GetELRfunc2()

	vars := mux.Vars(r)
	tableApi := vars["table"]
	var tableDb string

	var return_data string

	for a := range apiToDb {

		if a == tableApi {

			tableDb = apiToDb[a].table

			body, err := io.ReadAll(r.Body)
			if err != nil {
				return l(r, 8, err)
			}
			dec := json.NewDecoder(strings.NewReader(string(body)))

			switch tableDb {

			case "g_pitanja_c_testovi":
				var recordData models.Test
				if err := dec.Decode(&recordData); err != nil {
					return l(r, 8, err)
				}
				// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
				// OVO DA SE URADI I U POST I U PUT
				return_data, err = h.db.PostOne(recordData, r)
				if err != nil {
					return l(r, 8, err)
				}

			case "mi_users":
				var recordData models.User
				if err := dec.Decode(&recordData); err != nil {
					return l(r, 4, err)
				}

				if err := vet.ValidateUserData(recordData, r); err != nil {
					return err
				}

				// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
				// OVO DA SE URADI I U POST I U PUT
				//453253fsdf456

				return_data, err = h.db.PostOne(recordData, r)
				if err != nil {
					return l(r, 4, err)
				}
				// log.Println(return_data)
				l(r, 4, fmt.Errorf(return_data))

			case "g_user_blog":
				var recordData models.Note
				if err := dec.Decode(&recordData); err != nil {
					return l(r, 8, err)
				}
				// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
				// OVO DA SE URADI I U POST I U PUT
				return_data, err = h.db.PostOne(recordData, r)
				if err != nil {
					return l(r, 8, err)
				}

			case "v_settings":
				return nil

			default:
				return clr.NewAPIError(
					http.StatusNotAcceptable,
					"request data rejected",
				)
			}

		}
	}

	if tableDb == "" {
		return clr.NewAPIError(
			http.StatusNotAcceptable,
			"request data rejected",
		)
	} else {
		return clr.WriteJSON(w, 200, return_data)
	}

}

func (h *VezbamoHandler) HandleGetOne(w http.ResponseWriter, r *http.Request) error {
	// both work the same (sending json string)
	// but with w.Write there is no extra conversion to string but uses []byte from db
	// io.WriteString(w, string(db.GetQuestions()))

	l := clr.GetELRfunc2()

	vars := mux.Vars(r)
	tableApi := vars["table"]
	fieldApi := vars["field"]
	recordApi := vars["record"]

	var tableDb, fieldDb string
	var recordDb any

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

	for a := range apiToDb {
		// fmt.Println("deo od map:", a)
		if a == tableApi {
			tableDb = apiToDb[a].table
			switch fieldApi {
			case "id":
				fieldDb = apiToDb[a].id
			case "mail":
				fieldDb = apiToDb[a].mail
			}
		}
	}

	if tableDb == "" || fieldDb == "" {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	data, err := h.db.GetOne(tableDb, fieldDb, recordDb, r)
	if err != nil {
		return err
	}
	if data != nil {
		return clr.WriteJSON(w, 200, data)
	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

}

func (h *VezbamoHandler) HandleGetMany(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	tableApi := vars["table"]

	var tableDb string

	for a := range apiToDb {
		if a == tableApi {
			tableDb = apiToDb[a].table
		}
	}

	if tableDb == "" {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	data, err := h.db.GetMany(tableDb, r)
	if err != nil {
		return err
	}
	if data != nil {
		return clr.WriteJSON(w, 200, data)
	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

}

func (h *VezbamoHandler) HandlePutOne(w http.ResponseWriter, r *http.Request) error {
	l := clr.GetELRfunc2()

	vars := mux.Vars(r)

	tableApi := vars["table"]
	fieldApi := vars["field"]
	recordApi := vars["record"]

	var tableDb, fieldDb string
	var recordDb any

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

	var return_data string

	for a := range apiToDb {

		if a == tableApi {

			tableDb = apiToDb[a].table

			switch fieldApi {
			case "id":
				fieldDb = apiToDb[a].id
			case "mail":
				fieldDb = apiToDb[a].mail
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				return l(r, 8, err)
			}
			dec := json.NewDecoder(strings.NewReader(string(body)))

			switch tableDb {

			case "g_pitanja_c_testovi":
				var recordData models.Test
				if err := dec.Decode(&recordData); err != nil {
					return l(r, 8, err)
				}
				return_data, err = h.db.PutOne(tableDb, fieldDb, recordDb, recordData, r)
				if err != nil {
					return l(r, 8, err)
				}

			case "mi_users":
				return nil

			case "g_user_blog":
				var recordData models.Note
				if err := dec.Decode(&recordData); err != nil {
					return l(r, 8, err)
				}
				return_data, err = h.db.PutOne(tableDb, fieldDb, recordDb, recordData, r)
				if err != nil {
					return l(r, 8, err)
				}

			case "v_settings":
				return nil

			default:
				return clr.NewAPIError(
					http.StatusNotAcceptable,
					"no (available) content that conforms to the criteria given",
				)
			}

		}
	}

	if tableDb == "" {
		return clr.NewAPIError(
			http.StatusNotAcceptable,
			"no (available) content that conforms to the criteria given",
		)
	} else {
		return clr.WriteJSON(w, 200, return_data)
	}

}

/*
Briše jedan zapis iz bilo koje tabele
*/
func (h *VezbamoHandler) HandleDeleteOne(w http.ResponseWriter, r *http.Request) error {
	l := clr.GetELRfunc2()

	vars := mux.Vars(r)

	tableApi := vars["table"]
	fieldApi := vars["field"]
	recordApi := vars["record"]

	var tableDb, fieldDb string
	var recordDb any

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

	for a := range apiToDb {
		// fmt.Println("deo od map:", a)
		if a == tableApi {
			tableDb = apiToDb[a].table
			switch fieldApi {
			case "id":
				fieldDb = apiToDb[a].id
			case "mail":
				fieldDb = apiToDb[a].mail
			}
		}
	}

	if tableDb == "" || fieldDb == "" {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	}

	err := h.db.DeleteOne(tableDb, fieldDb, recordDb, r)
	if err != nil {
		return clr.NewAPIError(http.StatusNotAcceptable, "no (available) content that conforms to the criteria given")
	} else {
		return clr.NewAPIError(http.StatusOK, "one record deleted for: "+r.URL.Path)
	}

}
