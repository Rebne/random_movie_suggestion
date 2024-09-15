BINARY_NAME=myapp
PORT=8080

build:
	go build -o ${BINARY_NAME} main.go

run:
	go run main.go

templ:
	templ generate -watch -proxy=http://localhost:${PORT}

tailwind:
	tailwindcss -i ./styles.css -o public/styles.css --watch
