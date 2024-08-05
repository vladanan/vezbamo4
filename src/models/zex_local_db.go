package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vladanan/vezbamo4/src/controllers/clr"
)

func GetLocal(r *http.Request) {

	l := clr.GetELRfunc2()

	godotenv.Load("../../../.env")

	conn, err := pgx.Connect(context.Background(), os.Getenv("FEDORA_CONNECTION_STRING"))
	if err != nil {
		l(r, 8, err)
	}
	defer conn.Close(context.Background())

	// fmt.Println("fedora conn string:", os.Getenv("FEDORA_CONNECTION_STRING"))

	rows2, err := conn.Query(context.Background(), "SELECT * FROM app_g_user_blog;")
	if err != nil {
		l(r, 8, err)
	}
	for rows2.Next() {
		if val, err := rows2.Values(); err != nil {
			log.Print(err)
		} else {
			fmt.Println("proba local pg:", fmt.Sprint(val))
		}
	}

	// bytearray_tests, err2 := json.Marshal(pgx_tests)
	// if err2 != nil {
	// 	fmt.Printf("Json error: %v", err2)
	// }
	// jsonstring_pitanja := string(bytearray_pitanja)
	// fmt.Println("json string pitanja:", jsonstring_pitanja)

}
