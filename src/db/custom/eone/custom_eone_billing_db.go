package eone

import (
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vladanan/vezbamo4/src/clr"
	"github.com/vladanan/vezbamo4/src/models"
)

type DBeone struct{}

func (db DBeone) GetBilling(r *http.Request) (any, error) {

	l := clr.GetELRfunc2()

	godotenv.Load("../../../.env")

	conn, err := pgx.Connect(context.Background(), os.Getenv("FEDORA_CONNECTION_STRING"))
	if err != nil {
		return nil, l(r, 8, err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM billing ORDER BY id ASC;")
	if err != nil {
		return nil, l(r, 8, err)
	}

	pgxBilling, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Billing])
	if err != nil {
		return nil, l(r, 8, err)
	}

	// fmt.Println("string concat rows:", pgxTests)

	// bytearray_tests, err2 := json.Marshal(pgx_tests)
	// if err2 != nil {
	// 	fmt.Printf("Json error: %v", err2)
	// }s
	// jsonstring_pitanja := string(bytearray_pitanja)
	// fmt.Println("json string pitanja:", jsonstring_pitanja)

	return pgxBilling, nil

}
