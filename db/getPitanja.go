package db

import (
	"context"

	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"encoding/json"
	// pitanja "github.com/vladanan/vezbamo4/views/pitanja"
	// "github.com/a-h/templ"
)

type Pitanje struct {
	G_id   int8   `db:"g_id"`
	Tip    string `db:"tip"`
	Oblast string `db:"oblast"`
}

func GetPitanja() string {

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

	rows, _ := conn.Query(context.Background(), "SELECT g_id, tip, oblast FROM g_pitanja_c_testovi;")
	if err != nil {
		fmt.Printf("Unable to make query: %v\n", err)
	}

	sva_pitanja, err := pgx.CollectRows(rows, pgx.RowToStructByName[Pitanje])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		//return
	}

	jp, err2 := json.Marshal(sva_pitanja)
	if err2 != nil {
		fmt.Printf("Json error: %v", err2)
	}

	// fmt.Print(string(jp))

	return string(jp)

}
