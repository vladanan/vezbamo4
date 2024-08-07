// package for custom log and error
package clr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/encoding/json"
)

// **********************************************************************
// //**** UNUSED FOR GO COMPILER STUPIDITY

// kompajler javlja grešku za neiskorišćene promenljive i funkcije (za importe radi ekstenzija)
// tako da to ne može da se isključi linterom jer je hardcoded u kompajleru
// može da se zaobiđe tako da se promenljiva upotrebi sa _ kao: a:= 5; _ = a
// može i da se stavi van funkcija tj. kao globalna promenljivi i onda će da se pokazuje kao upozorenje ali će da prolazi kompajler
// još jedan workaround je funkcija koja upotrebljava sve takve promenljive i funkcijem

func X(x ...any) {}

// **********************************************************************

////**** CUSTOM LOGER

// https://pkg.go.dev/golang.org/x/exp/slog#Level
// const (
// 	LevelDebug = -4
// 	LevelInfo  = 0
// 	LevelWarn  = 4
// 	LevelError = 8
// )

// https://opentelemetry.io/docs/specs/otel/logs/data-model/
// Severity	range	Range name	Meaning
// 1-4						TRACE				A fine-grained debugging event. Typically disabled in default configurations.
// 5-8						DEBUG				A debugging event.
// 9-12						INFO				An informational event. Indicates that an event happened.
// 13-16					WARN				A warning event. Not an error but is likely more important than an informational event.
// 17-20					ERROR				An error event. Something went wrong.
// 21-24					FATAL				A fatal error such as application or system crash.

const ( // console escape characters for colors
	Reset          = "\033[0m"
	Black          = "\033[30m"
	Red            = "\033[31m"
	Green          = "\033[32m"
	Yellow         = "\033[33m"
	Blue           = "\033[34m"
	Magenta        = "\033[35m"
	Cyan           = "\033[36m"
	LightGray      = "\033[37m"
	Gray           = "\033[90m"
	LightRed       = "\033[91m"
	LightGreen     = "\033[92m"
	LightYellow    = "\033[93m"
	LightBlue      = "\033[94m"
	LightMagenta   = "\033[95m"
	LightCyan      = "\033[96m"
	White          = "\033[97m"
	BgBlack        = "\033[40m"
	BgRed          = "\033[41m"
	BgGreen        = "\033[42m"
	BgYellow       = "\033[43m"
	BgBlue         = "\033[44m"
	BgMagenta      = "\033[45m"
	BgCyan         = "\033[46m"
	BgLightGray    = "\033[47m"
	BgGray         = "\033[100m"
	BgLightRed     = "\033[101m"
	BgLightGreen   = "\033[102m"
	BgLightYellow  = "\033[103m"
	BgLightBlue    = "\033[104m"
	BgLightMagenta = "\033[105m"
	BgLightCyan    = "\033[106m"
	BgWhite        = "\033[107m"
)

var consoleColors = map[string]string{ // console escape characters for colors
	"Reset":          "\033[0m",
	"Black":          "\033[30m",
	"Red":            "\033[31m",
	"Green":          "\033[32m",
	"Yellow":         "\033[33m",
	"Blue":           "\033[34m",
	"Magenta":        "\033[35m",
	"Cyan":           "\033[36m",
	"LightGray":      "\033[37m",
	"Gray":           "\033[90m",
	"LightRed":       "\033[91m",
	"LightGreen":     "\033[92m",
	"LightYellow":    "\033[93m",
	"LightBlue":      "\033[94m",
	"LightMagenta":   "\033[95m",
	"LightCyan":      "\033[96m",
	"White":          "\033[97m",
	"BgBlack":        "\033[40m",
	"BgRed":          "\033[41m",
	"BgGreen":        "\033[42m",
	"BgYellow":       "\033[43m",
	"BgBlue":         "\033[44m",
	"BgMagenta":      "\033[45m",
	"BgCyan":         "\033[46m",
	"BgLightGray":    "\033[47m",
	"BgGray":         "\033[100m",
	"BgLightRed":     "\033[101m",
	"BgLightGreen":   "\033[102m",
	"BgLightYellow":  "\033[103m",
	"BgLightBlue":    "\033[104m",
	"BgLightMagenta": "\033[105m",
	"BgLightCyan":    "\033[106m",
	"BgWhite":        "\033[107m",
}

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
	Mu     sync.Mutex // ensures atomic writes; protects the following fields
	Prefix string     // prefix on each line to identify the logger (but see Lmsgprefix)
	Flag   int        // properties
	Out    io.Writer  // destination for output
	Buf    []byte     // for accumulating text to write
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
	if l.Flag&Lmsgprefix == 0 {
		*buf = append(*buf, l.Prefix...)
	}
	if l.Flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.Flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.Flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.Flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.Flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.Flag&(Lshortfile|Llongfile) != 0 {
		if l.Flag&Lshortfile != 0 {
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
	if l.Flag&Lmsgprefix != 0 {
		*buf = append(*buf, l.Prefix...)
	}
}

//
//

type FileLog struct {
	Date  string
	Time  string
	File  string
	Error string
	Path  string
}

// ubacuje log zapis na početak log fajla umesto kao što je default za write metode da rade append na kraju fajla
func prependLogToFile(file string, buf []byte) bool {

	// čišćenje teksta od oznaka za boje za log
	// jer je to nepotrebno u fajlovima i pravi probleme sa json parsingom
	bstring := string(buf)

	for color := range consoleColors {
		// log.Println("color from map", consoleColors[color])
		bstring = strings.ReplaceAll(bstring, consoleColors[color], "")
	}

	buf = []byte(bstring)

	// provera za test režim uz pomoć .env fajla jer se fajlovi u tekst režimu ne pokreću iz root nego iz mesta test fajla i onda ostali path ne valjaju
	if _, err := os.ReadFile(".env"); err != nil {
		file = "../../../" + file // ../../../
	}

	dat, err := os.ReadFile(file)
	if err != nil {
		log.Println(err)
		return false
	}

	// If the file doesn't exist, create it, or append to the file
	// sys, err := os.OpenFile("sys.log", os.O_CREATE|os.O_WRONLY, 0644)
	sys, err := os.OpenFile(file, os.O_WRONLY, 0644)
	if err != nil {
		// log.Fatal(err)
		log.Println(err)
		return false
	}
	// WriteAt will *OVERWRITE* the contents from the given offset, so your expected result "12A345" is incorrect. It is not possible to insert characters in the middle of the file with the WriteAt or Write methods.
	written, err := sys.WriteAt(buf, 0)
	if err != nil {
		sys.Close() // ignore error; Write error takes precedence
		// log.Fatal(err)
		log.Println(err)
		return false
	}
	if _, err := sys.WriteAt(dat, int64(written)); err != nil {
		sys.Close() // ignore error; Write error takes precedence
		// log.Fatal(err)
		log.Println(err)
		return false
	}
	sys.Sync()
	defer sys.Close()

	// parse .log file to .json log file
	dat, err = os.ReadFile(file)
	if err != nil {
		log.Println(err)
		return false
	}

	var jsonLog string

	lines := strings.Split(string(dat), "\n")

	for index, line := range lines {
		// log se sastoji iz datuma, vremena, fajla, greške i path
		// prva tri su odvojena od druga dva sa po dva razmaka
		// fajl ne sme da ima na kraju prazan novi red
		// niti sme da ima manje od dva razmaka između prva tri i greške i path
		dtfErrPath := strings.Split(line, "  ")
		dtf := strings.Split(dtfErrPath[0], " ")
		fileLog := FileLog{
			Date:  dtf[0],
			Time:  dtf[1],
			File:  dtf[2],
			Error: dtfErrPath[1],
			Path:  dtfErrPath[2],
		}
		bufJson := new(bytes.Buffer)
		if err := json.NewEncoder(bufJson).Encode(fileLog); err != nil {
			log.Println(err)
			return false
		}
		line = bufJson.String() // returns a string of what was written to it
		if index < len(lines)-1 {
			line = strings.ReplaceAll(line, "\n", ",")
		} else {
			line = strings.ReplaceAll(line, "\n", "")
		}
		jsonLog = jsonLog + line
	}

	// pravljenje json niza
	jsonLog = "[" + jsonLog + "]"

	// beautify json fajla
	jsonLog = strings.ReplaceAll(jsonLog, "{", "\n\t{\n\t\t")
	// ako se koristi `` da bi se našao i ubacio \n onda se on ubacuje kao takav i ne može da se escapuje
	// zato prvo ubacujem sedam , pa umesto njih \n i ostalo šta treba
	jsonLog = strings.ReplaceAll(jsonLog, `","Error":"`, `",,,,,,,"Error":"`)
	jsonLog = strings.ReplaceAll(jsonLog, `","Path":"`, `",,,,,,,"Path":"`)
	jsonLog = strings.ReplaceAll(jsonLog, ",,,,,,,", ",\n\t\t")
	jsonLog = strings.ReplaceAll(jsonLog, "},", "\n\t},")
	jsonLog = strings.ReplaceAll(jsonLog, "}]", "\n\t}\n]")

	if err := os.WriteFile(file+".json", []byte(jsonLog), os.ModePerm); err != nil {
		log.Println(err)
		return false
	}

	return true
}

// Modifikovana funkcija koja koja daje f za false, u za user{}
//
// log.Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) OutputIzmenjen(a any) (bool, any, error) {

	var msg_fe string
	for_sys_log := true
	for_usr_log := false
	var s string
	var e error
	if err, ok := a.(error); ok {
		e = err
	} else {
		e = errors.New(fmt.Sprint(a))
	}

	switch a.(type) {
	case string:
		msg_fe = fmt.Sprint(a)
		// s = BgLightBlue + " " + fmt.Sprint(a) + Reset
		s = Reset + LightYellow + " " + fmt.Sprint(a) + Reset
	case error:
		// s = Reset + LightMagenta + " " + fmt.Sprint(a) + Reset 989876kjhuj
		s = BgRed + " " + fmt.Sprint(a) + Reset
		if strings.Contains(s, "Pogrešna lozinka za:") {
			s = BgRed + " " + fmt.Sprint(a) + "  " + "tmp/r.URL.Path" + Reset
			// s = Reset + LightYellow + " " + e.Error() + "  " + "tmp/r.URL.Path" + Reset
			for_usr_log = true
			for_sys_log = false
			msg_fe = "Email_or_pass_wrong"
		}
	default:
		s = "Funkcija nije dobila ni string ni error!"
		a = errors.New("Funkcija nije dobila ni string ni error!")
	}

	calldepth := 1
	now := time.Now() // get this early.
	var file string
	var line int
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if l.Flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.Mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.Mu.Lock()
	}
	l.Buf = l.Buf[:0]
	l.formatHeader(&l.Buf, now, file, line)
	l.Buf = append(l.Buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.Buf = append(l.Buf, '\n')
	}
	_, err := l.Out.Write(l.Buf)
	if err != nil {
		log.Println(fmt.Sprint(err))
		l.Buf = append(l.Buf, fmt.Sprint(err)...)
	}

	switch a.(type) {
	case error:

		file := "sys.log"

		if for_sys_log {
			file = "sys.log"
		} else if for_usr_log {
			file = "usr.log"
		}

		ok := prependLogToFile(file, l.Buf)
		if !ok {
			log.Println("Nije uspelo dodavanje loga na fajl!")
		}

	default:

	}

	X(msg_fe)

	return false, nil, e
}

//
//

func (l *Logger) OutputIzmenjen2(r *http.Request, level int, e error) error {

	var s string
	var logfile string

	// level 0 info ne ide u fajlove
	// 4 warn je user error ide u user log
	// 8 erro server ide u sys log
	switch level {
	case 0:
		s = Reset + Blue + " " + e.Error() + "  " + r.URL.Path + Reset
	case 4:
		s = Reset + LightYellow + " " + e.Error() + "  " + r.URL.Path + Reset
		logfile = "tmp/usr.log"
	case 8:
		s = BgRed + " " + e.Error() + "  " + r.URL.Path + Reset
		logfile = "tmp/sys.log"
	default:
		s = BgMagenta + " " + e.Error() + "  " + r.URL.Path + Reset
		logfile = "tmp/sys.log"
	}

	if strings.Contains(s, "Pogrešna lozinka za:") {
		logfile = "tmp/usr.log"
	}

	calldepth := 1
	now := time.Now() // get this early.
	var file string
	var line int
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if l.Flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.Mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.Mu.Lock()
	}
	l.Buf = l.Buf[:0]
	l.formatHeader(&l.Buf, now, file, line)
	l.Buf = append(l.Buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.Buf = append(l.Buf, '\n')
	}
	_, err := l.Out.Write(l.Buf)
	if err != nil {
		log.Println(fmt.Sprint(err))
		l.Buf = append(l.Buf, fmt.Sprint(err)...)
	}

	if logfile != "" {
		if ok := prependLogToFile(logfile, l.Buf); !ok {
			log.Println("Nije uspelo dodavanje loga na fajl!")
		}
	}

	return e
}

/*
Daje:

	za l: Logger.OutputIzmenjen(string) koja daje (false, User{}, string) za AuthenticateUser()

Na taj način se rade tri stvari u vrlo malecnom if e != nil{} kodu koji:
  - hendluje error
  - loguje grešku na konzoli tamo gde je i nastala
  - šalje response false, models.User{} i mag_fe ruteru.
*/
func GetELRfunc() func(any) (bool, any, error) {
	loger := Logger{Out: os.Stdout, Prefix: BgBlue, Flag: log.LstdFlags | log.Lshortfile}
	return loger.OutputIzmenjen
}

func GetELRfunc2() func(r *http.Request, level int, e error) error {
	loger := Logger{Out: os.Stdout, Prefix: BgGreen, Flag: log.LstdFlags | log.Lshortfile}
	return loger.OutputIzmenjen2
}

// func GetELRvars_ex() (func(string) (bool, models.User, error), func(error) string) {
// 	loger := Logger{Out: os.Stdout, Prefix: "", Flag: log.LstdFlags | log.Lshortfile}
// 	return loger.OutputIzmenjen1, StringIzError
// }

// func GetELRvars() (func(error) string, func(string) (bool, models.User, error)) {
// 	loger := Logger{Out: os.Stdout, Prefix: "", Flag: log.LstdFlags | log.Lshortfile}
// 	return func(e error) string {
// 			return fmt.Sprintf("%s", e)
// 		},
// 		loger.OutputModified
// }

// **********************************************************************

////**** CUSTOM ERROR

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API erorr: %d, %v", e.StatusCode, e.Msg)
}

func NewAPIError(status int, msg any) APIError {
	return APIError{
		StatusCode: status,
		Msg:        msg,
	}
}

type APIfunc func(w http.ResponseWriter, r *http.Request) error

func CheckFunc(h APIfunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				WriteJSON(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error breee",
				}
				WriteJSON(w, http.StatusInternalServerError, errResp)
			}
			// slog.Error("http api error", "err", err.Error(), "path", r.URL.Path)
			// log.Println("http api error:", err.Error(), "path:", r.URL.Path)
			// slog.Error("on http api:", "path", r.URL.Path)
			// log.Print(Yellow + "error on http api path: " + err.Error() + r.URL.Path + Reset)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func CheckErr(err error) APIError {

	// log.Print(Magenta + "error on internal api call" + Reset)
	if apiErr, ok := err.(APIError); ok {
		return apiErr
	} else {
		return APIError{
			StatusCode: http.StatusInternalServerError,
			Msg:        "internal server error, check again later or contact support",
		}
	}

}

// type ServerError struct {
// 	StatusCode int `json:"statusCode"`
// 	Msg        any `json:"msg"`
// }

// func (e ServerError) Error() string {
// 	return fmt.Sprintf("API erorr %d", e.StatusCode)
// }

// func NewServerError(status int, msg any) ServerError {
// 	return ServerError{
// 		StatusCode: status,
// 		Msg:        msg,
// 	}
// }

// type UserError struct {
// 	StatusCode int `json:"statusCode"`
// 	Msg        any `json:"msg"`
// }

// func (e UserError) Error() string {
// 	return fmt.Sprintf("API erorr %d", e.StatusCode)
// }

// func NewUserError(status int, msg any) UserError {
// 	return UserError{
// 		StatusCode: status,
// 		Msg:        msg,
// 	}
// }

// **********************************************************************
