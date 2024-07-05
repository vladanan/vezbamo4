package errorlogres

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/vladanan/vezbamo4/src/models"
)

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

// Modifikovana funkcija koja koja daje f za false, u za user{}
// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) OutputIzmenjen(s string) (bool, models.User, error) {
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
	return false, models.User{}, err
}

// **********************************************************************

func StringIzError(e error) string {
	return fmt.Sprintf("%s", e)
}

/*
Daje:

	za l: Logger.OutputIzmenjen(string) koja daje (false, User{}, error) za AuthenticateUser()
	za i: StringIzError(error) koja daje string za l

Na taj način se rade tri stvari u vrlo malecnom if e != nil{} kodu koji:
  - hendluje error
  - loguje grešku na konzoli tamo gde je i nastala
  - šalje response bool, User{} i error routeru.
*/
func GetELRvars() (func(string) (bool, models.User, error), func(error) string) {
	loger := Logger{Out: os.Stdout, Prefix: "", Flag: log.LstdFlags | log.Lshortfile}
	return loger.OutputIzmenjen, StringIzError
}

// func GetELRvars() (func(error) string, func(string) (bool, models.User, error)) {
// 	loger := Logger{Out: os.Stdout, Prefix: "", Flag: log.LstdFlags | log.Lshortfile}
// 	return func(e error) string {
// 			return fmt.Sprintf("%s", e)
// 		},
// 		loger.OutputModified
// }

// type ShortError struct{}

// func (se ShortError) Error(er error) string {
// 	return fmt.Sprintf("%s", er)
// }

// **********************************************************************

// m := Logger{Out: os.Stdout, Prefix: "", Flag: log.LstdFlags | log.Lshortfile}
// o := m.J

// return F(), U(), m.J(I(e))
// m.O(P(e))
// m.O(s.R(e))
// m.o(1, s.r(e))
// m.o(1, e.Error())
// return f, u, l._(1, e._())
// return f, u, log.Output(1, e.Error())

// func Demo() (bool, models.User, error) {
// 	l := log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
// 	f := false
// 	u := models.User{}
// 	m := Logger{Out: os.Stdout, Prefix: "", Flag: log.LstdFlags | log.Lshortfile}

// 	_, e := strconv.Atoi("v")
// 	if e != nil {
// 		m.J(I(e))
// 		return f, u, m.J(I(e))
// 	}

// 	_, e = strconv.Atoi("v")
// 	if e != nil {
// 		l.Print(e)
// 		return f, u, nil
// 	}

// 	_, e = strconv.Atoi("v")
// 	if e != nil {
// 		return f, u, m.J(I(e))
// 	}

// 	return f, u, m.J(I(e))
// }
