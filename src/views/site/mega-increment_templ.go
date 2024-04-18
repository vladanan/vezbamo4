// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package site

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/gorilla/sessions"
	"github.com/vladanan/vezbamo4/src/views"
	"net/http"
)

func MegaIncrement(store sessions.Store, r *http.Request) templ.Component {
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
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"relative top-4\"><a href=\"https://marketplace.visualstudio.com/items?itemName=vladan-andjelkovic.mega-increment\" target=\"_blank\" rel=\"noopener noreferrer\" class=\"swiper-slide cursor-alias\"><div class=\"title\" data-swiper-parallax=\"-300\">MI</div><div class=\"subtitle\" data-swiper-parallax=\"-200\">Mega Increment</div><div class=\"text text-xs lg:text-lg text-sky-200 rounded-md bg-gradient-to-br from-transparent via-neutral-950 to-transparent \" style=\"max-width: 800px;\" data-swiper-parallax=\"-100\"><p>Mega-increment Visual Studio Code extension is intended to ease parallel independent incrementations and decrementations in various strings. It can be used in writing code for lists, enums, arrays, tests, html and xml tags, csv files, data base examples and tests, date-time iterations, hexadecimal and binary register allocations and many other uses.<br><br><span class=\"text-sky-300\">Click to see extension Readme page at Visual Studio Code marketplace.</span></p></div></a> <a href=\"https://github.com/vladanan/mega-increment/blob/master/PublicAPIDocs.md\" target=\"_blank\" rel=\"noopener noreferrer\" class=\"swiper-slide cursor-alias\"><div class=\"title\" data-swiper-parallax=\"-300\">MACf</div><div class=\"subtitle\" data-swiper-parallax=\"-200\">Mega Increment API</div><div class=\"text text-xs lg:text-lg text-sky-200 rounded-md bg-gradient-to-br from-neutral-950 via-transparent to-stone-950 \" style=\"max-width: 800px;\" data-swiper-parallax=\"-100\"><p>This site is hosting public API for Mega Increment Visual Studio Code extension Core functions.<br><br><span class=\"text-sky-300 font-bold\">Click to check the Docs</span><br><br>Does this API cover all functionalities as GUI for Advanced options in VSCode extension?<br><br>Some of the functionalities are tied to specifics of working space at VSCode editor or GUI Advanced version but all text processing available at extension is also available from this MACf API.</p></div></a></div><br><br><br><br><br><br>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = views.Layout(store, r).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
