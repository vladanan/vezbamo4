// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Footer() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"relative top-10 left-2 text-blue-300 text-xs text-left\">Vezbamo: 4.0.0-alpha.22.1</p><p class=\"relative top-10 left-2 text-blue-400 text-xs text-left\">Copyright <span><a href=\"auto_login_admin\">&copy;</a></span> Vladan Anđelković 2022-2024.</p><a href=\"dashboard\" class=\"absolute top-10 right-2 text-blue-500 text-xs text-left\">dashboard</a> <a href=\"auto_login_user\" class=\"absolute top-14 right-2 text-red-500 text-xs text-left\">AUTO user login</a> <a href=\"komponents\" class=\"absolute top-20 right-2 text-green-500 text-xs text-left\">test komponents</a><!--\n  <a\n    href=\"https://onedrive.live.com/edit?id=D1C8EFB22DF66B8F!sb3e8a975419448a99ae3522050054c98&resid=D1C8EFB22DF66B8F!sb3e8a975419448a99ae3522050054c98&cid=d1c8efb22df66b8f&ithint=file%2Cpptx&redeem=aHR0cHM6Ly8xZHJ2Lm1zL3AvYy9kMWM4ZWZiMjJkZjY2YjhmL0VYV3A2TE9VUWFsSW11TlNJRkFGVEpnQlhJbVdJaTgxXzZxb2tZdEN3UElDbEE_ZT00OlVjOHdBdCZhdD05&migratedtospo=true&wdo=2\"\n    target=\"_blank\"\n\t\trel=\"noopener noreferrer\"\n    class=\"absolute top-20 right-2 text-blue-400 text-xs text-left hover:bg-sky-900 hover:text-blue-300 px-2 pb-1 border border-red-100 rounded-md\"\n  >tmp link za milicinu prezentaciju<span class=\"text-purple-300 text-lg mt-2\">&nbsp;&#8635;</span></a>\n  -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
