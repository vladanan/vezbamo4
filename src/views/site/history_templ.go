// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package site

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/gorilla/sessions"
	"github.com/vladanan/vezbamo4/src/views"
	"net/http"
)

func History(store sessions.Store, r *http.Request) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"relative top-4 m-5 text-blue-300\"><p class=\"text-xl text-blue-300\">History about project</p><p class=\"text-xl text-blue-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(views.Translate(store, r, "UnWelcome"))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `src/views/site/history.templ`, Line: 15, Col: 78}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><p>[This page is work in progress]</p><br><p>Few years ago first version of Vezbamo project was written in PHP and used only for assignments for first grade primary and college mathematics. https://vezbamo1.onrender.com/</p><br><p>Some time later second version was implemented in MySQL, Node.JS/Espress.JS, EJS and Bootstrap. It's not fully funictional any more but is still alive and is maintained as historical milestone for Vezbamo project. https://vezbamo04.vercel.app/</p><br><p>Third implementation is with Next.JS/React.JS @ Vercel, PostgreSQL @ Supabase, Tailwind CSS. It's hosting API for Mega Increment Visual Studio Code extension as well as admin dahsboard for Vezbamo. https://vezbamo.vercel.app/</p><br><p>Newest version has Go on backend (PostgreSQL @ Supabase) and HTMX+Tailwind on frontend with a little bit of React. It also implements i18n and auth becoming complete and functional full-stack project with services for both Vezbamo and Mega Increment.</p><br><p>copy o vezbamo pitanja/mi/custom/telemetry apis, uxv</p><br><p>sva komunikacija na sh sr en en-US</p></div><div class=\"m-3 p-2 md:max-w-full bg-white rounded-2xl shadow-lg\">Izvori za fotografije na sajtu:<ul><li class=\"mx-5 list-disc font-bold\">tests and questions:</li><li class=\"mx-5 list-disc\">Image by <a href=\"https://pixabay.com/users/firmbee-663163/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=620822\">Firmbee</a> from <a href=\"https://pixabay.com/photos/office-business-accountant-620822/\">Pixabay</a></li><li class=\"mx-5 list-disc\">Image by <a href=\"https://pixabay.com/users/stocksnap-894430/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=2569523\">StockSnap</a> from <a href=\"https://pixabay.com/photos/people-man-guy-millenials-busy-2569523/\">Pixabay</a></li><li class=\"mx-5 list-disc\">Photo by <a href=\"https://unsplash.com/@zacharytnelson?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText\">Zachary Nelson</a> on <a href=\"https://unsplash.com/photos/photo-of-three-men-jumping-on-ground-near-bare-trees-during-daytime-98Elr-LIvD8\">Unsplash</a></li><li class=\"mx-5 list-disc\">Image by <a href=\"https://pixabay.com/users/absolutvision-6158753/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=7774314\">Gino Crescoli</a> from <a href=\"https://pixabay.com/illustrations/translation-keyboard-computer-7774314/\">Pixabay</a></li><li class=\"mx-5 list-disc font-bold\">assignments:</li><li class=\"mx-5 list-disc\">Image by <a href=\"https://pixabay.com/users/pexels-2286921/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=1866497\">Pexels</a> from <a href=\"https://pixabay.com/photos/abacus-classroom-count-counter-1866497/\">Pixabay</a></li><li class=\"mx-5 list-disc\"><a href=\"https://giphy.com/gifs/americasgottalent-agt-americas-got-talent-l0MYOpH6uxygkOCvS\">kombinacije</a></li><li class=\"mx-5 list-disc\"><a href=\"https://i.gifer.com/ZJlt.gif\">patke</a></li><li class=\"mx-5 list-disc\"><a href=\"https://wifflegif.com/gifs/294265-ellen-page-interview-ellen-page-juggling-gif\">permutacije link 1</a> <a href=\"https://64.media.tumblr.com/a438cfdca27894e9b2a7b5460a3cef26/tumblr_msw2tr1h3s1qdtt0bo2_r1_250.gif\">permutacije link 2</a></li><li class=\"mx-5 list-disc\"><a href=\"https://skiphursh.tumblr.com/post/166431151274\">varijacije</a></li><li class=\"mx-5 list-disc font-bold\">questions api:</li><li class=\"mx-5 list-disc\">Photo by <a href=\"https://unsplash.com/@windows?utm_content=creditCopyText&amp;utm_medium=referral&amp;utm_source=unsplash\">Windows</a> on <a href=\"https://unsplash.com/photos/men-and-women-sitting-and-standing-while-staring-at-laptop-p74ndnYWRY4?utm_content=creditCopyText&amp;utm_medium=referral&amp;utm_source=unsplash\">Unsplash</a></li><li class=\"mx-5 list-disc font-bold\">user portal:</li><li class=\"mx-5 list-disc\">Image by <a href=\"https://pixabay.com/users/makieni777-16975393/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=5331883\">makieni777</a> from <a href=\"https://pixabay.com/photos/cat-work-technology-pet-cute-pc-5331883/\">Pixabay</a></li><li class=\"mx-5 list-disc font-bold\">custom apis:</li><li class=\"mx-5 list-disc\">Image by <a href=\"https://pixabay.com/users/absolutvision-6158753/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=7773520\">Gino Crescoli</a> from <a href=\"https://pixabay.com/illustrations/computer-language-html-computer-web-7773520/\">Pixabay</a></li></ul></div><br><br><br><br><br><br>")
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
