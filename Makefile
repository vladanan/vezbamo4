build:
	@go build -o bin src/main.go

dev:
	@air2 -c .air.templ.toml