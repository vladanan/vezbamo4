# 0) basic build pre deploy
build:
	@go build -o bin src/main.go

go_live_reload:
	@air

templ_live_reload:
	@air2 -c .air.templ.toml

# 1) watch .go files but no templ etc. then triggers templ generate by changing fake.templ file
devgo:
	@air -c .air.go_and_triger_templ.toml

# 2) temp gnerate to watch templ files and to perform hot reload in browser
devtg:
	@templ generate --watch --proxy="http://0.0.0.0:10000" --cmd="go run ./src/main.go"

# 3) tailwind to watch changes in views folder for html, templ, js
devtw:
	@tailwindcss -i assets/input.css -o assets/output.css --watch

# 2+3) NE RADI: combined templ and tailwind but tailwind is not working well with some classes
devtt:
	@templ generate --watch --proxy="http://0.0.0.0:10000" --cmd="go run ./src/main.go & tailwindcss -i input.css -o assets/output.css"
