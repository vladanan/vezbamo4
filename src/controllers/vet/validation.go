package vet

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/vladanan/vezbamo4/src/controllers/clr"
	"github.com/vladanan/vezbamo4/src/models"
)

func ValidateUserData(recordData models.User, r *http.Request) error {
	l := clr.GetELRfunc2()

	// validacija za UPIS NOVOG KORISNIKA a-zA-Z09 .,+-*:!?() min char 8 max 32 ISTO URADITI I NA FE UZ ARGUMENTS I JS

	// https://gobyexample.com/regular-expressions
	// https://pkg.go.dev/regexp
	// https://regex101.com/

	// matched, err := regexp.MatchString(`[^a-zA-Z\d]`, recordData.Email)
	re, err := regexp.Compile(`[^a-zA-Z\d@.-_]g`)
	if err != nil {
		return l(r, 4, err)
	}

	// fmt.Println("regex match:", re_email.MatchString(recordData.Email))

	if len(recordData.Email) < 8 ||
		len(recordData.Email) > 32 ||
		!strings.ContainsAny(recordData.Email, "@.") ||
		re.MatchString(recordData.Email) {
		return l(r, 4, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax mail"))
	}

	if len(recordData.Hash_lozinka) < 8 ||
		len(recordData.Hash_lozinka) > 32 ||
		re.MatchString(recordData.Hash_lozinka) {
		return l(r, 4, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax pass"))
	}

	if len(recordData.User_name) < 8 ||
		len(recordData.User_name) > 32 ||
		re.MatchString(recordData.User_name) {
		return l(r, 4, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax name"))
	}

	return nil
}
