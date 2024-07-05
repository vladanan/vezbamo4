// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package site

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/vladanan/vezbamo4/src/views"
	"net/http"
)

type Note struct {
	B_id        int    `db:"b_id"`
	Ime_tag     string `db:"ime_tag"`
	Mejl        string `db:"mejl"`
	Tema        string `db:"tema"`
	Poruka      string `db:"poruka"`
	User_id     string `db:"user_id"`
	User_mail   string `db:"user_mail"`
	From_url    string `db:"from_url"`
	Datum_upisa any    `db:"datum_upisa"`
}

func to_struct(notes []byte) []Note {
	var p []Note
	err := json.Unmarshal(notes, &p)
	if err != nil {
		fmt.Printf("Json error: %v", err)
	}
	return p
}

func reverser(notes []Note) []Note {
	var structLength int
	var newNotes []Note
	for i := range notes {
		structLength = i
	}
	for j := structLength; j > -1; j-- {
		newNotes = append(newNotes, notes[j])
	}
	return newNotes
}

func UserPortal(store sessions.Store, r *http.Request, notes []byte) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"m-5 text-blue-300 font-bold text-3xl\">Notes, contact, FAQ</p><ul class=\"\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, item := range reverser(to_struct(notes)) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"m-5 px-2 border text-xl  rounded-m text-blue-300  \">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(item.B_id))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/site/user_portal.templ`, Line: 49, Col: 94}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(": ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(item.Ime_tag)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/site/user_portal.templ`, Line: 49, Col: 112}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(", ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var5 string
				templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(item.Tema)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/site/user_portal.templ`, Line: 49, Col: 127}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(", ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(item.Poruka)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/site/user_portal.templ`, Line: 49, Col: 144}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul><div class=\"m-5 text-blue-300\"><p>Najčešće postavljana pitanja Zašto na sajtu nema odgovora na sva pitanja? Zato što, po Godingovoj teoremi nekompletnosti, u ovom svetu ili skup čenjenica nikada ne može biti kompletan ili sve činjenice nikada ne mogu biti dokazane. To je nešto slično onome što bi rekao Vitgenštajn da se ovaj svet može razumeti samo ako se neko popne iznad njega. Druga mogućnost je da korisnik koji je unosio pitanja nije uneo odgovore za neka od njih. To nije moguće potpuno iskontrolisati a svakako ne škodi da se odgovor potraži u knjizi ili svesci ;)<p></p>Da li svako može da upisuje svoja pitanja i testove na sajt? Da. Korisnici mogu da upišu svoje komplete i pitanja a možemo i mi to da uradimo umesto njih po povoljnim cenama koje su izložene ovde. Pre upisa svakako je poželjno pregledati ono što je već upisano. Pogledati i ostale napomene na pravilima sajta.</p><p>Mislio sam da ću ovde naći gotova pitanja i testove za oblast koja me zanima. Zašto ih nema? Ovaj sajt je namenjen da bude pomoć u učenju kroz tri koraka: pisanje pitanja i testova, prolazak kroz nasumično ispitivanje, rešavanje testa. To znači da je pisanje pitanja i testova (koje drugi već nisu upisali na sajt) sastavni deo učenja i mnogo bolje pomaže u učenju određene materije nego korišćenje već gotovih testova. To je zato što sam proces pregledavanja (tekstualnih, zvučnih i video) materijala i pisanja testova predstavlja značajan deo u procesu učenja. Korišćenje prećica u procesu učenja nije uvek mudro.</p><p></p></div><br><br><br><br><br><br>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.Layout(store, r).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
