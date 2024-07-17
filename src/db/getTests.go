package db

import (
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	elr "github.com/vladanan/vezbamo4/src/errorlogres"
	"github.com/vladanan/vezbamo4/src/models"
)

func GetTests(r *http.Request) ([]models.Test, error) {

	l := elr.GetELRfunc2()

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

	rows, err := conn.Query(context.Background(), "SELECT g_id, tip, oblast FROM g_pitanja_c_testovi;")
	if err != nil {
		return nil, l(r, 8, err)
	}

	pgxTests, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Test])
	if err != nil {
		return nil, l(r, 8, err)
	}

	// rows2, err := conn.Query(context.Background(), "SELECT g_id, tip, oblast, obrazovni_profil FROM g_pitanja_c_testovi;")
	// if err != nil {
	// 	return nil, l(r, 8, err)
	// }
	// for rows2.Next() {
	// 	if val, err := rows2.Values(); err != nil {
	// 		log.Print(err)
	// 	} else {
	// 		println("proba pgx:", fmt.Sprint(val))
	// 	}
	// }

	// bytearray_tests, err2 := json.Marshal(pgx_tests)
	// if err2 != nil {
	// 	fmt.Printf("Json error: %v", err2)
	// }
	// jsonstring_pitanja := string(bytearray_pitanja)
	// fmt.Println("json string pitanja:", jsonstring_pitanja)

	return pgxTests, nil

}
