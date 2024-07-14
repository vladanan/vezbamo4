package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	elr "github.com/vladanan/vezbamo4/src/errorlogres"
	"github.com/vladanan/vezbamo4/src/models"
)

func GetTests() ([]models.Test, error) {

	l := elr.GetELRfunc2()

	//https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		return nil, l(err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT g_id, tip, oblast FROM g_pitanja_c_testovi;")
	if err != nil {
		return nil, l(err)
		// fmt.Printf("Unable to make query: %v\n", err)
	}

	pgxTests, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Test])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		//return 56
	}

	// bytearray_tests, err2 := json.Marshal(pgx_tests)
	// if err2 != nil {
	// 	fmt.Printf("Json error: %v", err2)
	// }

	// jsonstring_pitanja := string(bytearray_pitanja)
	// fmt.Println("json string pitanja:", jsonstring_pitanja)

	log.Println("db okej")

	return pgxTests, nil

}
