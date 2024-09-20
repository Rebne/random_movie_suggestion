ROOT=cmd
BINARY_NAME=myapp
PORT=8080

build:
	go build -o ${BINARY_NAME} ${ROOT}/.

run:
	go run ${ROOT}/main.go

templ:
	templ generate -watch -proxy=http://localhost:${PORT}

css:
	tailwindcss -i ./web/styles/styles.css -o ./web/static/css/styles.css --watch

tailwind:
	tailwindcss -i ./web/styles/styles.css -o ./web/static/css/styles.css