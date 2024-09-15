BINARY_NAME=myapp

build:
	go build -o ${BINARY_NAME} main.go

run:
	go run main.go

# templ:
# 	templ generate -watch -proxy=http://localhost:3000

tailwind:
	tailwindcss -i view/css/input.css -o public/styles.css --watch