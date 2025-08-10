FROM node:20-alpine AS css-builder

WORKDIR /app

COPY package.json package-lock.json ./
COPY web/styles ./web/styles
COPY tailwind.config.js ./

RUN npm install
RUN npx tailwindcss -i ./web/styles/styles.css -o ./styles.css

FROM golang:alpine3.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY web/static ./web/static/
COPY --from=css-builder /app/styles.css ./web/static/css/styles.css

COPY web/views/home/*.txt ./web/views/home/

COPY web/views/home/*.go ./web/views/home/

COPY services/*.go ./services/

COPY models/*.go ./models/

COPY helpers/*.go ./helpers/

COPY handlers/*.go ./handlers/

COPY data/ ./data/

COPY cmd/*.go ./

COPY .env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /movie-generator

CMD ["/movie-generator"]