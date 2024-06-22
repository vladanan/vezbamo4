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

func getLang(store sessions.Store, r *http.Request) []string {
	session, _ := store.Get(r, "vezbamo.onrender.com-users")
	lang_map := session.Values["language"]
	accept := r.Header.Get("Accept-Language")
	language := ""

	if lang_map != nil {
		language = lang_map.(string)
	}
	if language == "" {
		language = strings.Split(accept, ",")[0]
	}

	var languageNames = map[string]string{
		"ar": "Arabic    - العربية: ar",
		"zh": "Chinese - 中文 (汉语): zh",
		"en": "English  : en",
		"sh": "Ex-yug srpskohrvatski: sh",
		"fr": "French   - français: fr",
		"de": "German - Deutch: de",
		"hi": "Hindi      - हिन्दी: hi",
		"it": "Italian    ;- italiano: it",
		"ru": "Russian  - русский: ru",
		"sr": "Serbian  - српски: sr",
		"es": "Spanish  - español: es",
	}

	// fmt.Println("lang:", language, "accept:", accept, "0:", strings.Split(accept, ",")[0])

	return []string{language, languageNames[language]}
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"relative top-0 left-2 w-6\"><a href=\"/\" class=\"\"><img src=\"static/site/vezbamo_logo6.svg\" height=\"25\" width=\"25\" alt=\"Vezbamo\"></a></div><a href=\"/\" class=\"absolute text-blue-400 top-1 left-10 text-sm\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(Translate(store, r, "Home"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 111, Col: 94}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a><div class=\"absolute top-0 right-2\"><select onchange=\"sendLang(event)\" name=\"lang\" id=\"lang\" class=\" w-[72px] text-sm bg-sky-950 text-blue-400\"><option class=\"text-blue-100 font-bold \" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(getLang(store, r)[0]))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(getLang(store, r)[1])
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 124, Col: 94}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option> <option disabled value=\"ar\">Arabic &nbsp;&nbsp;&nbsp;- العربية: ar</option> <option disabled value=\"zh\">Chinese - 中文 (汉语): zh</option> <option value=\"en\">English&nbsp;&nbsp;: en</option> <option value=\"sh\">Ex-yug srpskohrvatski: sh</option> <option disabled value=\"fr\">French &nbsp;&nbsp;- français: fr</option> <option disabled value=\"de\">German - Deutch: de</option> <option disabled value=\"hi\">Hindi &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- हिन्दी: hi</option> <option disabled value=\"it\">Italian &nbsp;&nbsp;&nbsp;- italiano: it</option> <option disabled value=\"ru\">Russian &nbsp;- русский: ru</option> <option disabled value=\"sr\">Serbian &nbsp;- српски: sr</option> <option disabled value=\"es\">Spanish &nbsp;- español: es</option></select>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if session, err := store.Get(r, "vezbamo.onrender.com-users"); err == nil {
			if auth, _ := session.Values["authenticated"].(bool); !auth {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("   <button class=\"text-sm ml-1 px-1 text-blue-300\" type=\"button\"><a href=\"/sign_in\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(Translate(store, r, "Sign_in"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 150, Col: 55}
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
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"text-sm ml-1 px-1 bg-gradient-to-r from-blue-700 to-blue-400 rounded-sm\" type=\"button\"><a href=\"/sign_out\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var5 string
				templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(Translate(store, r, "Sign_out"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/header.templ`, Line: 154, Col: 57}
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><script>\n\t\tfunction pageReload() {\n\t\t\tlocation.reload()\n\t\t}\n\t\tfunction delayReload2() {\n\t\t\tsetTimeout(pageReload(), 2000);\n\t\t}\n\t</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
