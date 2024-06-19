package models

import "time"

type User struct {
	U_id           int       `db:"u_id"`
	Created_at     time.Time `db:"created_at"`
	Hash_lozinka   string    `db:"hash_lozinka"`
	Email          string    `db:"email"`
	User_name      string    `db:"user_name"`
	Mode           string    `db:"user_mode"`
	Level          string    `db:"user_level"`
	Basic          bool      `db:"basic"`
	Js             bool      `db:"js"`
	C              bool      `db:"c"`
	Payment_date   time.Time `db:"payment_date"`
	Payment_expire time.Time `db:"payment_expire"`
	Verified_email string    `db:"verified_email"`
}
