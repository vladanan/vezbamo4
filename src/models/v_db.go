package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vladanan/vezbamo4/src/controllers/clr"
)

func (db DB) PostOne(table string, recordData any, r *http.Request) (string, error) {

	l := clr.GetELRfunc2()

	godotenv.Load("../../../.env")

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		return "", l(r, 8, err)
	}
	defer conn.Close(context.Background())

	switch data := recordData.(type) {

	case Test:
		commandTag, err := conn.Exec(context.Background(), `INSERT INTO `+table+`
			(
				tip,
				obrazovni_profil,
				razred,
				predmet,
				oblast
			)
				VALUES ($1, $2, $3, $4, $5);`,
			data.Tip,
			data.Obrazovni_profil,
			data.Razred,
			data.Predmet,
			data.Oblast,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		if commandTag.String() != "INSERT 0 1" {
			return "", l(r, 0, fmt.Errorf("no records inserted"))
		} else {
			newRecord, err := json.Marshal(data)
			if err != nil {
				return "", l(r, 8, err)
			}
			// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
			// OVO DA SE URADI I U POST I U PUT
			return string(newRecord), nil
		}

	case User:
		newRecord, err := json.Marshal(data)
		if err != nil {
			return "", l(r, 8, err)
		}
		// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
		// OVO DA SE URADI I U POST I U PUT
		return string(newRecord), clr.NewAPIError(418, "još mora da se radi na add user")

	case Note:
		commandTag, err := conn.Exec(context.Background(), `INSERT INTO `+table+`
			(
				ime_tag,
				mejl,
				tema,
				poruka,
				user_id
			)
				VALUES ($1, $2, $3, $4, $5);`,
			data.Ime_tag,
			data.Mejl,
			data.Tema,
			data.Poruka,
			data.User_id,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		if commandTag.String() != "INSERT 0 1" {
			return "", l(r, 0, fmt.Errorf("no records inserted"))
		} else {
			newRecord, err := json.Marshal(data)
			if err != nil {
				return "", l(r, 8, err)
			}
			// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
			// OVO DA SE URADI I U POST I U PUT
			return string(newRecord), nil
		}

	default:
		return "", l(r, 8, fmt.Errorf("post record ne pripada nijednom tipu"))
	}

}

func (db DB) GetOne(table string, field string, record any, r *http.Request) (any, error) {

	l := clr.GetELRfunc2()

	//https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx

	// sistem radi i bez učitavanja .env jer je valjda već učitano u routes.go ali kada se radi unit test onda mora i ovde jer se prilikom testa izgleda ne učitavaju svi fajlovi nego samo ono što je potrebno
	// zato se ovde .env učitava samo sa pathom za unit test jer sa normalan režim ovde nema ni potrebe da se učitava .env
	// zato nije ni potrebno da se reaguje ni kada van test okruženja učitavanje pukne jer je već pravilno učitano u routes
	// godotenv.Load(".env")
	godotenv.Load("../../../.env")

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		return nil, l(r, 8, err)
	}
	defer conn.Close(context.Background())

	var pgxData any

	rows, err := conn.Query(context.Background(), "SELECT * FROM "+table+" WHERE "+field+"=$1;", record)
	if err != nil {
		return nil, l(r, 8, err)
	}

	// for rows.Next() {
	// 	if val, err := rows.Values(); err != nil {
	// 		fmt.Println("rows greška:", err)
	// 		// return nil, l(r, 8, err)
	// 	} else {
	// 		fmt.Println("row:", fmt.Sprint(val))
	// 	}
	// }

	switch table {
	case "g_pitanja_c_testovi":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Test])
		if err != nil {
			return nil, l(r, 8, err)
		}
	case "mi_users":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[User])
		if err != nil {
			return nil, l(r, 8, err)
		}
	case "g_user_blog":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Note])
		if err != nil {
			return nil, l(r, 8, err)
		}
	case "v_settings":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Settings])
		if err != nil {
			return nil, l(r, 8, err)
		}
	default:
		return nil, l(r, 8, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax 3"))
	}

	if fmt.Sprint(pgxData) == "[]" {
		return nil, nil
	}

	return pgxData, nil

}

func (db DB) GetMany(table string, r *http.Request) (any, error) {

	l := clr.GetELRfunc2()

	godotenv.Load("../../../.env")

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		return nil, l(r, 8, err)
	}
	defer conn.Close(context.Background())

	var pgxData any

	rows, err := conn.Query(context.Background(), "SELECT * FROM "+table+";")
	if err != nil {
		return nil, l(r, 8, err)
	}

	switch table {
	case "g_pitanja_c_testovi":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Test])
		if err != nil {
			return nil, l(r, 8, err)
		}
	case "mi_users":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[User])
		if err != nil {
			return nil, l(r, 8, err)
		}
	case "g_user_blog":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Note])
		if err != nil {
			return nil, l(r, 8, err)
		}
	case "v_settings":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Settings])
		if err != nil {
			return nil, l(r, 8, err)
		}
	default:
		return nil, l(r, 8, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax 3"))
	}

	if fmt.Sprint(pgxData) == "[]" {
		return nil, nil
	}

	return pgxData, nil

}

func (db DB) PutOne(table string, field string, record any, recordData any, r *http.Request) (string, error) {

	l := clr.GetELRfunc2()

	godotenv.Load("../../../.env")

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		return "", l(r, 8, err)
	}
	defer conn.Close(context.Background())

	switch data := recordData.(type) {

	case Test:
		// za update napraviti kod koji na osnovu poslatih polja za izmenu i već postojećih napravi skroz novi upis za isti id tako da se izbegnu kompleksni (string contactenation) query i kompleksan kod
		commandTag, err := conn.Exec(context.Background(), `UPDATE `+table+` SET
			tip=$1,
			obrazovni_profil=$2,
			razred=$3,
			predmet=$4,
			oblast=$5
			WHERE `+field+`=$6;`,
			data.Tip,
			data.Obrazovni_profil,
			data.Razred,
			data.Predmet,
			data.Oblast,
			record,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		if commandTag.String() != "UPDATE 1" {
			return "", l(r, 0, fmt.Errorf("no records updated"))
		} else {
			newRecord, err := json.Marshal(data)
			if err != nil {
				return "", l(r, 8, err)
			}
			// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
			// OVO DA SE URADI I U POST I U PUT
			return string(newRecord), nil
		}

	case Note:
		// za update napraviti kod koji na osnovu poslatih polja za izmenu i već postojećih napravi skroz novi upis za isti id tako da se izbegnu kompleksni (string contactenation) query i kompleksan kod
		commandTag, err := conn.Exec(context.Background(), `UPDATE `+table+` SET
			ime_tag=$1,
			mejl=$2,
			tema=$3,
			poruka=$4,
			user_id=$5
			WHERE `+field+`=$6;`,
			data.Ime_tag,
			data.Mejl,
			data.Tema,
			data.Poruka,
			data.User_id,
			record,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		// fmt.Println("put commandTag:", commandTag)
		if commandTag.String() != "UPDATE 1" {
			return "", l(r, 0, fmt.Errorf("no records updated"))
		} else {
			newRecord, err := json.Marshal(data)
			if err != nil {
				return "", l(r, 8, err)
			}
			// OVO VRAĆA KOMPLETAN TIP TJ. STRUKTURU TABELE SA NAZIVIMA POLJA U DB: PROMENITI DA VRATI SAMO POSLATE PODATKE A IMENA POLJA DA SE SAKRIJU U TIPU MODELS TAKOD A NE BUDU ISTA KAO U DB OSIM VELIKIH SLOVA
			// OVO DA SE URADI I U POST I U PUT
			return string(newRecord), nil
		}

	default:
		return "", l(r, 8, fmt.Errorf("put record ne pripada nijednom tipu"))
	}

}

func (db DB) DeleteOne(table string, field string, record any, r *http.Request) error {

	l := clr.GetELRfunc2()

	godotenv.Load("../../../.env")

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		return l(r, 8, err)
	}
	defer conn.Close(context.Background())

	if r.Method == "DELETE" {

		commandTag, err := conn.Exec(context.Background(), "DELETE FROM "+table+" WHERE "+field+"=$1;", record)
		log.Println("delete test:", commandTag, commandTag.String())
		if err != nil {
			return l(r, 8, err)
		}
		if commandTag.String() == "DELETE 0" {
			return l(r, 0, fmt.Errorf("no such record for delete"))
		} else {
			l(r, 0, fmt.Errorf("record deleted"))
			return nil
		}
	} else {
		return l(r, 8, fmt.Errorf("wrong http method api call"))
	}

}
