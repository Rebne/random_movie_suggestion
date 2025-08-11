FROM node:20-alpine AS css-builder

WORKDIR /app

COPY . .

RUN npm install
RUN mkdir -p web/static/css
RUN npx tailwindcss -i ./web/styles/styles.css -o ./web/static/css/styles.css

FROM golang:alpine3.20 AS builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@v0.2.778

COPY . .
COPY --from=css-builder /app/web/static/css/styles.css ./web/static/css/styles.css

RUN go mod download
RUN templ generate
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -o /movie-generator

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app /app
COPY --from=builder /movie-generator ./movie-generator

CMD ["./movie-generator"]