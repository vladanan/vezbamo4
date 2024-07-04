package models

import (
	"time"
)

type User struct {
	U_id                 int       `db:"u_id"`
	Created_at_time      time.Time `db:"created_at_time"`
	Hash_lozinka         string    `db:"hash_lozinka"`
	Email                string    `db:"email"`
	User_name            string    `db:"user_name"`
	Mode                 string    `db:"user_mode"`
	Level                string    `db:"user_level"`
	Basic                bool      `db:"basic"`
	Js                   bool      `db:"js"`
	C                    bool      `db:"c"`
	Payment_date         time.Time `db:"payment_date"`
	Payment_expire       time.Time `db:"payment_expire"`
	Payment_amount       int       `db:"payment_amount"`
	Payment_currency     string    `db:"payment_currency"`
	Verified_email       string    `db:"verified_email"`
	Last_sign_in_time    time.Time `db:"last_sign_in_time"`
	Last_sign_in_headers string    `db:"last_sign_in_headers"`
	Created_at_headers   string    `db:"created_at_headers"`
	Deleted_at           time.Time `db:"deleted_at"`
	Bad_sign_in_attempts int       `db:"bad_sign_in_attempts"`
	Bad_sign_in_time     time.Time `db:"bad_sign_in_time"`
	Updated_at           time.Time `db:"updated_at"`
}

type Settings struct {
	S_id                       int       `db:"s_id"`
	Updated_at                 time.Time `db:"updated_at"`
	Bad_sign_in_attempts_limit string    `db:"bad_sign_in_attempts_limit"`
	Bad_sign_in_time_limit     string    `db:"bad_sign_in_time_limit"`
	Same_ip_sign_up_time_limit string    `db:"same_ip_sign_up_time_limit"`
}
