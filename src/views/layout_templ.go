// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func Translate(store sessions.Store, r *http.Request, item string) string {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("assets/i18n/active.en.toml")
	bundle.MustLoadMessageFile("assets/i18n/active.sh.toml")
	// bundle.MustLoadMessageFile("assets/i18n/active.es.toml")

	// https://pkg.go.dev/github.com/gorilla/sessions@v1.2.2#section-documentation
	// https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-cookie-same-site-00
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
	// https://stackoverflow.com/questions/67821709/this-set-cookie-didnt-specify-a-samesite-attribute-and-was-default-to-samesi

	session, err := store.Get(r, "vezbamo.onrender.com-users")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError) 876876
		// return
		log.Println("layout: Translate: Error on get store:", err)
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
		// SameSite: http.SameSiteDefaultMode,
		// SameSite: http.SameSiteLaxMode,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSite(0),
	}

	mail_map := session.Values["user_email"]

	if mail_map == nil {
		session.Values["user_email"] = "aaa"
	}

	auth_map := session.Values["authenticated"]

	if auth_map == nil {
		session.Values["authenticated"] = false
	}

	lang_map := session.Values["language"]
	sessionLanguage := ""

	if lang_map != nil {
		sessionLanguage = lang_map.(string)
	}

	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")

	if sessionLanguage != "" {
		accept = sessionLanguage
	}

	//fmt.Println("language: ", lang, "header: ", accept) ,jk

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

func Layout(store sessions.Store, r *http.Request) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en-US\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><!-- <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1\" /> --><link rel=\"icon\" href=\"static/site/vezbamo_ico4.svg\"><title>Vezbamo</title><script src=\"static/reload.js\"></script><script src=\"static/htmx.min.js\"></script><link href=\"static/output.css\" rel=\"stylesheet\"><!-- Link Swiper's CSS --><link rel=\"stylesheet\" href=\"https://cdn.jsdelivr.net/npm/swiper@11/swiper-bundle.min.css\"><!-- Demo styles height: 100%; --><style>\n\t\t\t\t.swiper {\n\t\t\t\t\twidth: 100%;\n\t\t\t\t\theight: 550px;\n\n\t\t\t\t}\n\n\t\t\t\t.swiper-slide {\n\t\t\t\t\tfont-size: 18px;\n\t\t\t\t\tcolor: #fff;\n\t\t\t\t\t-webkit-box-sizing: border-box;\n\t\t\t\t\tbox-sizing: border-box;\n\t\t\t\t\tpadding: 40px 60px;\n\t\t\t\t}\n\n\t\t\t\t.parallax-bg {\n\t\t\t\t\tposition: absolute;\n\t\t\t\t\tleft: 0;\n\t\t\t\t\ttop: 0;\n\t\t\t\t\twidth: 130%;\n\t\t\t\t\theight: 100%;\n\t\t\t\t\t-webkit-background-size: cover;\n\t\t\t\t\tbackground-size: cover;\n\t\t\t\t\tbackground-position: center;\n\t\t\t\t}\n\n\t\t\t\t.swiper-slide .title {\n\t\t\t\t\tfont-size: 41px;\n\t\t\t\t\tfont-weight: 300;\n\t\t\t\t}\n\n\t\t\t\t.swiper-slide .subtitle {\n\t\t\t\t\tpadding: 0.5rem;\n\t\t\t\t\tfont-size: 21px;\n\t\t\t\t\tmax-width: 400px;\n\t\t\t\t\tline-height: 1.1;\n\t\t\t\t}\n\n\t\t\t\t.swiper-slide .text {\n\t\t\t\t\tpadding: 0.5rem;\n\t\t\t\t\tmargin-top: 10px;\n\t\t\t\t\t\n\t\t\t\t\tmax-width: 400px;\n\t\t\t\t\tline-height: 1.5;\n\t\t\t\t}\n\t\t\t</style></head><body class=\"relative m-auto max-w-7xl top-2 h-screen bg-gradient-to-br from-slate-950 via-sky-950 to-blue-950\"><div id=\"heading\" class=\"relative h-auto\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Header(store, r).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div id=\"wrapper\" class=\"relative h-auto\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div id=\"footer\" class=\"relative h-auto\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Footer().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><!-- Swiper JS --><script src=\"https://cdn.jsdelivr.net/npm/swiper@11/swiper-bundle.min.js\"></script><!-- Initialize Swiper --><script>\n\t\t\t\tvar swiper = new Swiper(\".mySwiper\", {\n\t\t\t\t\tspeed: 600,\n\t\t\t\t\tparallax: true,\n\t\t\t\t\tpagination: {\n\t\t\t\t\t\tel: \".swiper-pagination\",\n\t\t\t\t\t\tclickable: true,\n\t\t\t\t\t},\n\t\t\t\t\tnavigation: {\n\t\t\t\t\t\tnextEl: \".swiper-button-next\",\n\t\t\t\t\t\tprevEl: \".swiper-button-prev\",\n\t\t\t\t\t},\n\t\t\t\t});\n\t\t\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
