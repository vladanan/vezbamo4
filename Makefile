# 0) osnovni build pre deploy
build:
	@go build -o bin src/main.go

go_live_reload:
	@air

templ_live_reload:
	@air2 -c .air.templ.toml

# 1) motri .go fajlove ali ne i templ itd. i da trigeruje templ generate izmenom fake.templ fajla
devgo:
	@air -c .air.go_and_triger_templ.toml

# 2) temp gnerate da prati templ fajlove i da radi hot reload u browseru
devtg:
	@templ generate --watch --proxy="http://0.0.0.0:10000" --cmd="go run ./src/main.go"

# 3) tailwind da prati izmene u scc
devtw:
	@tailwindcss -i views/input.css -o assets/output.css --watch

# 2+3) kombinovani temp i tailwind ali zeza sa bojama
devtt:
	@templ generate --watch --proxy="http://0.0.0.0:10000" --cmd="go run ./src/main.go & tailwindcss -i views/input.css -o assets/output.css"
