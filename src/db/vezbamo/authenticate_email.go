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

// type User struct {
// 	Hash_lozinka string `db:"hash_lozinka"`
// 	Email        string `db:"email"`
// 	Test         string `db:"test"`
// }

func AuthenticateEmail(key, mail string) bool {
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	// PROVERA ZA key (urađena)
	// PROVERA ZA i za mail da se ne desi da ima neki zaostali isti key sa kreiranja nekog drugog neverifikovanog naloga

	// password := []byte(password_str)

	// ciphertext, err := bcrypt.GenerateFromPassword(password, 5) //df
	// _, err := bcrypt.GenerateFromPassword(password, 5) //df
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error from bcrypt encryption: %s\n", err)
	// 	return false, models.User{}
	// }

	// fmt.Println("Ciphertext: ", string(ciphertext))

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		// os.Exit(1)
		return false
	}
	defer conn.Close(context.Background())

	// commandTag, err := conn.Exec(context.Background(), "INSERT INTO mi_users (email, basic, js, c, hash_lozinka) VALUES ($1, $2, $3, $4, $5);", email, true, true, true, ciphertext)
	// if err != nil {
	// 	fmt.Printf("Unable to connect to database: %v\n", err)
	// }
	// fmt.Printf("insert result: %v\n", commandTag)

	// err = bcrypt.CompareHashAndPassword([]byte("$2a$05$HYej4fyvWYnC5LvrSGEOD.bztJzcYn45t62etOTrN8d59BkoD7fhy"), password)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error from bcrypt Dencryption: %s\n", err)
	// 	return
	// } else {
	// 	fmt.Printf("\nProšlo je!")
	// }

	rows, err2 := conn.Query(context.Background(), "SELECT * FROM mi_users where verified_email=$1 and email=$2;", key, mail)
	if err2 != nil {
		fmt.Printf("Unable to make query: %v\n", err2)
		return false
	}

	pgxKey, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		//return
		return false
	}

	// fmt.Print("pgx user::", pgx_user)

	bytearrayKey, err2 := json.Marshal(pgxKey)
	if err2 != nil {
		fmt.Printf("Json error: %v", err2)
	}

	var structUser models.User

	if string(bytearrayKey) != "null" {
		structUser = toStruct(bytearrayKey)[0]
	} else {
		return false // ovo se dešava ako je ključ netačan ili u bazi nema uopšte tog ključa jer je mejl već verifikovan
	}

	// fmt.Print(struct_user)

	if len(structUser.Verified_email) > 8 {

		fmt.Print("Key for db write", structUser.Verified_email)
		//     update mytab set c=3, d=4, e=5 where a = 0;
		commandTag, err := conn.Exec(context.Background(), `UPDATE mi_users SET verified_email=$1 where verified_email=$2;`,
			"verified",
			structUser.Verified_email,
		)
		if err != nil {
			fmt.Printf("Unable to connect to database to write verified field: %v\n", err)
			return false
		} else {
			fmt.Printf("insert result: %v\n", commandTag)

			// AKO JE SVE OKEJ
			fmt.Printf("\nmail verified: Prošlo je!\n")
			return true
		}

	} else {
		return false
	}

}
