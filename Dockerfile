FROM golang:alpine3.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY public ./public/

COPY home_templ.txt ./

COPY *.go ./

COPY .env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /movie-generator

COPY id_data.json ./

CMD ["/movie-generator"]