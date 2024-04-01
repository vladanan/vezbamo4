build:
	@go build -o bin src/main.go

go_live_reload:
	@air

templ_live_reload:
	@air2 -c .air.templ.toml

templ_watch_hot_reload:
	@templ generate --watch --proxy="http://0.0.0.0:10000" --cmd="go run ./src/main.go"

tailwind_watch:
	@tailwindcss -i views/input.css -o assets/output.css --watch

templ_tail:
	@templ generate --watch --proxy="http://0.0.0.0:10000" --cmd="go run ./src/main.go & tailwindcss -i views/input.css -o assets/output.css"
