// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
	// "fmt"
)

// func showLanguage (store sessions.Store, r *http.Request, item string) bool {

// 	session, err := store.Get(r, "vezbamo.onrender.com-users")
// 	if err != nil {
// 		// http.Error(w, err.Error(), http.StatusInternalServerError)
// 		// return
// 		fmt.Println("Error on get store:", err)
// 	}

// 	lang_map := session.Values["language"]
// 	sessionLanguage := ""

// 	if lang_map != nil {
// 		sessionLanguage = lang_map.(string)
// 	}

// 	accept := r.Header.Get("Accept-Language")

// 	// fmt.Println("\nshadow:", sessionLanguage, item, strings.Split(accept, ",")[0])

// 	if sessionLanguage == "" && item == "browser" {
// 		return true
// 	} else
// 	if (
// 		(sessionLanguage == "en-US" || sessionLanguage == "") &&
// 		(sessionLanguage == "en-US" || (strings.Split(accept, ",")[0] == "en-US" || strings.Split(accept, ",")[0] == "en")) &&
// 		item == "en") {
// 		return true
// 	} else
// 	if (
// 		(sessionLanguage == "es" || sessionLanguage == "") &&
// 		(sessionLanguage == "es" || strings.Split(accept, ",")[0] == "es") &&
// 		item == "es") {
// 		return true
// 	} else
// 	if (
// 		(sessionLanguage == "sr" || sessionLanguage == "") &&
// 		(sessionLanguage == "sr" || strings.Split(accept, ",")[0] == "sr") &&
// 		item == "sr") {
// 		return true
// 	} else {
// 		return false
// 	}

// }

// var cssheader_lang_button string = `text-sm ml-1 px-1 text-red-500`

// templ header_lang_button(link string, class string, text string, store sessions.Store, r *http.Request) {
// 	<button
// 		hx-post={"/"+link}
// 		hx-on:click="delayReload()"
// 		class={class, templ.KV("text-blue-300", !showLanguage(store, r, link)), templ.KV("text-blue-100", showLanguage(store, r, link))}
// 	>
// 		{text}
// 	</button>
// }

func getLang(store sessions.Store, r *http.Request) string {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")
	lang_map := session.Values["language"]
	accept := r.Header.Get("Accept-Language")
	language := ""
	// fmt.Println("lang, accept, 0:", language, accept, strings.Split(accept, ",")[0])

	// fdsafSDGF

	if lang_map != nil {
		language = lang_map.(string)
	}

	if language == "" {
		language = strings.Split(accept, ",")[0]
	}

	return language
}

func Header(store sessions.Store, r *http.Request) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"relative top-0 left-2 w-6\"><a href=\"/\" class=\"\"><img src=\"static/vezbamo_logo6.svg\" height=\"25\" width=\"25\" alt=\"Vezbamo\"></a></div><a href=\"/\" class=\"absolute text-blue-400 top-1 left-10 text-sm\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(Translate(store, r, "Home"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 99, Col: 94}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a><div class=\"absolute top-0 right-2\"><select onchange=\"sendLang(event)\" name=\"lang\" id=\"lang\" class=\"text-sm bg-sky-950 text-blue-400\"><option class=\"text-blue-100 font-bold \" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(getLang(store, r)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(getLang(store, r))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 112, Col: 88}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option> <option value=\"eng\">eng</option> <option value=\"srh\">srh</option> <option value=\"esp\">esp</option> <option value=\"browser\">auto</option></select>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if session, err := store.Get(r, "vezbamo.onrender.com-users"); err == nil {
			if auth, _ := session.Values["authenticated"].(bool); !auth {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("  <button class=\"text-sm ml-1 px-1 text-blue-300\" type=\"button\"><a href=\"/login\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(Translate(store, r, "Login"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 129, Col: 51}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"text-sm ml-1 px-1 bg-gradient-to-r from-blue-700 to-blue-400 rounded-sm\" type=\"button\"><a href=\"/logout\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var5 string
				templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(Translate(store, r, "Logout"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 133, Col: 53}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
