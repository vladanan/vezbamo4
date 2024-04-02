package db

import (
	"context"

	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"encoding/json"
)

type Pitanje struct {
	G_id   int8   `db:"g_id"`
	Tip    string `db:"tip"`
	Oblast string `db:"oblast"`
}

func GetPitanja() []byte {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	//https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, _ := conn.Query(context.Background(), "SELECT g_id, tip, oblast FROM g_pitanja_c_testovi;")
	if err != nil {
		fmt.Printf("Unable to make query: %v\n", err)
	}

	pgx_pitanja, err := pgx.CollectRows(rows, pgx.RowToStructByName[Pitanje])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		//return
	}

	bytearray_pitanja, err2 := json.Marshal(pgx_pitanja)
	if err2 != nil {
		fmt.Printf("Json error: %v", err2)
	}

	// jsonstring_pitanja := string(bytearray_pitanja)
	// fmt.Println("json string pitanja:", jsonstring_pitanja)

	return bytearray_pitanja

}
