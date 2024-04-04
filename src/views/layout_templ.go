// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"net/http"
)

func Translate(globalLanguage string, r *http.Request, item string) string {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("assets/i18n/active.en.toml")
	bundle.MustLoadMessageFile("assets/i18n/active.es.toml")
	bundle.MustLoadMessageFile("assets/i18n/active.sr.toml")

	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")

	if globalLanguage != "" {
		accept = globalLanguage
	}

	//fmt.Println("language: ", lang, "header: ", accept)

	localizer := i18n.NewLocalizer(bundle, lang, accept)

	helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: item,
		},
	})

	// myUnreadEmails := localizer.MustLocalize(&i18n.LocalizeConfig{
	// 	DefaultMessage: &i18n.Message{
	// 		ID:          "MyUnreadEmails",
	// 		Description: "The number of unread emails I have",
	// 		One:         "I have {{.PluralCount}} unread email.",
	// 		Other:       "I have {{.PluralCount}} unread emails.",
	// 	},

	// })

	return helloPerson
}

func Layout(globalLanguage string, r *http.Request) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en-US\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><link href=\"static/output.css\" rel=\"stylesheet\"><script src=\"static/htmx.min.js\"></script><script src=\"static/reload.js\"></script><link rel=\"icon\" href=\"static/vezbamo_ico4.svg\"><title>Vezbamo</title></head><body class=\"relative m-auto sm:w-auto max-w-md h-max bg-gradient-to-br from-sky-400 via-emerald-200 to-amber-200\"><div id=\"heading\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Header(globalLanguage, r).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div id=\"wrapper\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div id=\"footer\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Footer().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
