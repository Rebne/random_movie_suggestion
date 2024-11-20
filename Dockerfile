FROM golang:alpine3.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY web/static ./web/static/

COPY web/views/home/*.txt ./web/views/home/

COPY web/views/home/*.go ./web/views/home/

COPY services/*.go ./services/

COPY models/*.go ./models/

COPY helpers/*.go ./helpers/

COPY handlers/*.go ./handlers/

COPY data/*.go ./data/


COPY cmd/*.go ./

COPY .env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /movie-generator

CMD ["/movie-generator"]