package models

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vladanan/vezbamo4/src/clr"
)

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

	// var rows pgx.Rows
	// var pgxTests []models.Test

	// fmt.Println("iz api za query:", table, field, record)

	// rows, err := conn.Query(context.Background(), "SELECT * FROM g_pitanja_c_testovi WHERE g_id=$1;", record)
	// if err != nil {
	// 	return nil, l(r, 8, err)
	// }

	// switch {
	// case g_id == 0 && r.Method == "GET":
	// 	rows, err = conn.Query(context.Background(), "SELECT g_id, tip, oblast FROM g_pitanja_c_testovi;")
	// 	if err != nil {
	// 		return nil, l(r, 8, err)
	// 	}

	// case g_id != 0 && r.Method == "GET":
	// 	rows, err = conn.Query(context.Background(), "SELECT g_id, tip, oblast FROM g_pitanja_c_testovi WHERE g_id=$1;", g_id)
	// 	if err != nil {
	// 		return nil, l(r, 8, err)
	// 	}

	// case g_id == 22 && r.Method == "POST":
	// 	_, err := conn.Exec(context.Background(), `INSERT INTO g_pitanja_c_testovi
	// 	(
	// 		tip,
	// 		oblast
	// 	)
	// 		VALUES ($1, $2);`,
	// 		"test",
	// 		"go language",
	// 	)
	// 	// log.Println("new test:", commandTag)
	// 	if err != nil {
	// 		return nil, l(r, 8, err)
	// 	}
	// 	pgxTests = []models.Test{}
	// 	return pgxTests, nil

	// case g_id == 37 && r.Method == "DELETE":
	// 	commandTag, err := conn.Exec(context.Background(), `DELETE FROM g_pitanja_c_testovi WHERE g_id = $1`, g_id)
	// 	log.Println("delete test:", commandTag, commandTag.String())
	// 	if err != nil {
	// 		return nil, l(r, 8, err)
	// 	}
	// 	if commandTag.String() == "DELETE 0" {
	// 		return nil, l(r, 0, fmt.Errorf("no such record for delete"))
	// 	}
	// 	pgxTests = []models.Test{}
	// 	return pgxTests, nil

	// case g_id == 38 && r.Method == "PUT":
	// 	// za update napraviti kod koji na osnovu poslatih polja za izmenu i već postojećih napravi skroz novi upis za isti id tako da se izbegnu kompleksni (string contactenation) query i kompleksan kod
	// 	commandTag, err := conn.Exec(context.Background(), `UPDATE g_pitanja_c_testovi SET obrazovni_profil=$1 WHERE
	// 		g_id=$2`,
	// 		"programeri ccccccccc		",
	// 		g_id)
	// 	log.Println("update test:", commandTag)
	// 	if err != nil {
	// 		return nil, l(r, 8, err)
	// 	}
	// 	pgxTests = []models.Test{}
	// 	return pgxTests, nil

	// default:
	// 	return nil, l(r, 4, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax"))
	// }

	// if rows != nil {
	// 	pgxTests, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.Test])
	// 	if err != nil {
	// 		return nil, l(r, 8, err)
	// 	}
	// }

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

	// fmt.Println("string concat rows:", pgxData)
	// fmt.Println("string concat rows:", pgxTests)

	// bytearray_tests, err2 := json.Marshal(pgx_tests)
	// if err2 != nil {
	// 	fmt.Printf("Json error: %v", err2)
	// }s
	// jsonstring_pitanja := string(bytearray_pitanja)
	// fmt.Println("json string pitanja:", jsonstring_pitanja)

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