package dbvezbamo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/vladanan/vezbamo4/src/models"
)

func GetNotes() []byte {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// err = conn.Ping(context.Background())
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unsuccessful ping test to database: %v\n", err)
	// 	os.Exit(1)
	// }

	rows, _ := conn.Query(context.Background(), "SELECT * FROM g_user_blog;")
	if err != nil {
		fmt.Printf("Unable to make query: %v\n", err)
	}

	pgxNotes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Note])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		//return
	}

	// fmt.Print(pgx_notes[0])
	// for _, b := range pgx_notes {
	// 	fmt.Printf("%v, %s: $%s\n", b.B_id, b.Poruka, b.Tema)
	// }

	bytearrayNotes, err2 := json.Marshal(pgxNotes)
	if err2 != nil {
		fmt.Printf("Json error: %v", err2)
	}

	/*

		pgx also implements QueryRow in the same style as database/sql.

		var name string
		var weight int64
		err := conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
		if err != nil {
				return err
		}
		Use Exec to execute a query that does not return a result set.

		commandTag, err := conn.Exec(context.Background(), "delete from widgets where id=$1", 42)
		if err != nil {
				return err
		}
		if commandTag.RowsAffected() != 1 {
				return errors.New("No row found to delete")
		}

	*/

	// commandTag, err := conn.Exec(context.Background(), "INSERT INTO g_user_blog (ime_tag, mejl, tema, poruka, user_id, user_mail, from_url) VALUES ($1, $2, $3, $4, $5, $6, $7)", "vladan", "v@vezbamo.onrender.com", "pg test", "go pg insert test", "1234", "n@n.com", "www.vezbamo.itd")
	// if err != nil {
	// 	// return err
	// 	fmt.Printf("Insert postgresql error: %v\n", err)
	// }
	// if commandTag.RowsAffected() != 1 {
	// 	// return errors.New("No row found to delete")
	// 	fmt.Printf("No row found: %v\n", err)
	// }

	// fmt.Println(string(bytearray_notes))

	return bytearrayNotes
}
