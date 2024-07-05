package db

import (
	"context"
	// "errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	// "sync/atomic"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	model "github.com/vladanan/vezbamo4/src/models"

	"log"
)

func to_struct(user []byte) []model.User {
	var p []model.User
	err := json.Unmarshal(user, &p)
	if err != nil {
		fmt.Printf("Json error: %v", err)
	}
	return p
}

// **********************************************************************

// type ll struct {
// 	lls string
// }

// func (ll ll) Output2(level int, s string) string {
// 	return "log poruka: " + strconv.Itoa(level) + " " + s
// }

// type ee struct {
// 	ees string
// }

// func (ee ee) Error2() string {
// 	return ee.ees
// }

// **********************************************************************

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix on each line to identify the logger (but see Lmsgprefix)
	flag   int        // properties
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
	// isDiscard atomic.Bool // whether out == io.Discard
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   - l.prefix (if it's not blank and Lmsgprefix is unset),
//   - date and/or time (if corresponding flags are provided),
//   - file and line number (if corresponding flags are provided),
//   - l.prefix (if it's not blank and Lmsgprefix is set).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	if l.flag&Lmsgprefix == 0 {
		*buf = append(*buf, l.prefix...)
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
	if l.flag&Lmsgprefix != 0 {
		*buf = append(*buf, l.prefix...)
	}
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) o(s string) error {
	calldepth := 1
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

// **********************************************************************

type shortError struct {
}

// Adding this Error method makes argError implement the error interface.

func (e shortError) r(er error) string {
	return fmt.Sprintf("%s", er)
}

// **********************************************************************

func AuthenticateUser(email string, password_str string, already_authenticated bool, r *http.Request) (bool, model.User) {
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/

	f := false
	u := model.User{}
	l := log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
	m := Logger{out: os.Stdout, prefix: "", flag: log.LstdFlags | log.Lshortfile}
	s := shortError{}

	//2354213jdjh232

	// llv := ll{lls: "string iz ll"}
	// eev := ee{ees: "string iz ee"}

	if _, e := strconv.Atoi("v"); e != nil {
		m.o(s.r(e))
		// m.o(1, s.r(e))
		// m.o(1, e.Error())
		// return f, u, l._(1, e._())
		// return f, u, log.Output(1, e.Error())
	}

	password := []byte(password_str)

	// ENV, BAZA, UZIMANJE USERA

	e := godotenv.Load(".env")
	if e != nil {
		l.Print(e)
		return f, u
	}

	conn, e := pgx.Connect(context.Background(), os.Getenv("SUPABASE_CONNECTION_STRING"))
	if e != nil {
		l.Print(e)
		return f, u
		// os.Exit(1)
	}
	defer conn.Close(context.Background())
	rows, e := conn.Query(context.Background(), "SELECT * FROM mi_users where email=$1;", email)
	if e != nil {
		l.Print(e)
		return f, u
	}
	pgx_user, e := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if e != nil {
		l.Print(e)
		return f, u
	}
	// fmt.Print("AuthenticateUser: pgx user:", pgx_user)
	bytearray_user, e := json.Marshal(pgx_user)
	if e != nil {
		l.Print(e)
		return f, u
	}

	// fmt.Print("bytearray user: ", bytearray_user)

	var struct_user model.User
	if string(bytearray_user) != "null" { // array nije prazan tj. ima zapisa sa odgovarajućim mejlom
		struct_user = pgx_user[0] //to_struct(bytearray_user)[0]
	}

	// fmt.Print("AuthenticateUser: ",struct_user)

	// UZIMANJE PROMENLJIVIH IZ ENV I DB ZA BAD ATTEMPT LIMITE

	var bad_sign_in_attempts_limit int64 = 2
	var bad_sign_in_time_limit int64 = 8

	BAD_SIGN_IN_ATTEMPTS_LIMIT, err := strconv.ParseInt(os.Getenv("BAD_SIGN_IN_ATTEMPTS_LIMIT"), 0, 8)
	if err != nil {
		BAD_SIGN_IN_ATTEMPTS_LIMIT = 0
	}
	BAD_SIGN_IN_TIME_LIMIT, err := strconv.ParseInt(os.Getenv("BAD_SIGN_IN_TIME_LIMIT"), 0, 8)
	if err != nil {
		BAD_SIGN_IN_TIME_LIMIT = 0
	}

	rows2, e := conn.Query(context.Background(), "SELECT * FROM v_settings where s_id=1;")
	if e != nil {
		l.Print(e)
		return f, u
	}
	pgx_settings, e := pgx.CollectRows(rows2, pgx.RowToStructByName[model.Settings])
	if e != nil {
		l.Print(e)
		return f, u
	}

	db_bad_sign_in_attempts_limit, err := strconv.ParseInt(pgx_settings[0].Bad_sign_in_attempts_limit, 0, 8)
	if err != nil {
		db_bad_sign_in_attempts_limit = 0
	}
	db_bad_sign_in_time_limit, err := strconv.ParseInt(pgx_settings[0].Bad_sign_in_time_limit, 0, 8)
	if err != nil {
		db_bad_sign_in_time_limit = 0
	}

	if BAD_SIGN_IN_ATTEMPTS_LIMIT != 0 {
		bad_sign_in_attempts_limit = BAD_SIGN_IN_ATTEMPTS_LIMIT
	} else if db_bad_sign_in_attempts_limit != 0 {
		bad_sign_in_attempts_limit = db_bad_sign_in_attempts_limit
	}
	if BAD_SIGN_IN_TIME_LIMIT != 0 {
		bad_sign_in_time_limit = BAD_SIGN_IN_TIME_LIMIT
	} else if db_bad_sign_in_time_limit != 0 {
		bad_sign_in_time_limit = db_bad_sign_in_time_limit
	}

	// fmt.Print("AuthenticateUser: bad env: ", BAD_SIGN_IN_ATTEMPTS_LIMIT, BAD_SIGN_IN_TIME_LIMIT, "\n")
	// fmt.Print("AuthenticateUser: bad db: ", db_bad_sign_in_attempts_limit, db_bad_sign_in_time_limit, "\n")
	// fmt.Print("AuthenticateUser: bad real: ", bad_sign_in_attempts_limit, bad_sign_in_time_limit, "\n")

	if already_authenticated {

		fmt.Print("AuthenticateUser: Already authenticated\n")
		struct_user.Hash_lozinka = ""
		return true, struct_user

	} else if int64(struct_user.Bad_sign_in_attempts) < bad_sign_in_attempts_limit { // ako je broj neuselih pokušaja manji od limita ide se na dalje proverese broj neuspelih pokušaja

		fmt.Print("AuthenticateUser: add 1 to bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")

		_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
			struct_user.Bad_sign_in_attempts+1,
			email,
		)
		if e != nil {
			l.Print(e)
			return f, u
		}

		// fmt.Print("AuthenticateUser: prošlo minuta od poslednjeg lošeg sign in-a: ", time.Since(struct_user.Bad_sign_in_time).Minutes(), "\n")
		fmt.Print("AuthenticateUser: set last bad sign time\n")
		_, e = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_time=$1 where email=$2;`,
			time.Now(),
			email,
		)
		if e != nil {
			l.Print(e)
			return f, u
		}

		if string(bytearray_user) != "null" { // array nije prazan tj. ima zapisa sa odgovarajućim mejlom

			e = bcrypt.CompareHashAndPassword([]byte(struct_user.Hash_lozinka), password) // provera lozinke

			if e != nil {
				// fmt.Fprintf(os.Stderr, "AuthenticateUser: Loša lozinka: %s\n", err)
				l.Print(e)
				return f, u

			} else if struct_user.Verified_email == "verified" { // ako je lozinka dobra onda se proverava da li je mejl verifikovan

				_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_time=$1 where email=$2;`,
					time.Now(),
					email,
				)
				if e != nil {
					l.Print(e)
					return f, u
				}

				bytearray_headers, e := json.Marshal(r.Header)
				if e != nil {
					l.Print(e)
					return f, u
				}
				_, e = conn.Exec(context.Background(), `UPDATE mi_users SET last_sign_in_headers=$1 where email=$2;`,
					string(bytearray_headers),
					email,
				)
				if e != nil {
					l.Print(e)
					return f, u
				}

				fmt.Print("\nAuthenticateUser: zero bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")

				_, e = conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
					0,
					email,
				)
				if e != nil {
					l.Print(e)
					return f, u
				}

				// fmt.Print("\nAuthenticateUser: Prošlo je!\n")
				struct_user.Hash_lozinka = ""
				return true, struct_user

			} else {

				l.Print("Mejl nije verifikovan!\n")
				return f, u

			}

		} else {

			l.Print("Nema korisnika sa tim mejlom i lozinkom", email, password_str, "\n")
			return f, u

		}

	} else { // ako je broj neuselih pokušaja veći od limita gleda se da li je prošlo više vremena od limita

		if time.Since(struct_user.Bad_sign_in_time).Minutes() < float64(bad_sign_in_time_limit) {

			l.Print("AuthenticateUser: previše pokušaja za sign in\n")
			l.Print("AuthenticateUser: pokušati za minuta:", float64(bad_sign_in_time_limit)-time.Since(struct_user.Bad_sign_in_time).Minutes(), "\n")
			return f, u

		} else {

			fmt.Print("AuthenticateUser: zeroing bad sign in attempts:", struct_user.Bad_sign_in_attempts, "\n")
			_, e := conn.Exec(context.Background(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
				0,
				email,
			)
			if e != nil {
				l.Print(e)
				return f, u
			}

			l.Print("Moguće je ponovo probati sign in\n")
			return f, u

		}

	}

}
