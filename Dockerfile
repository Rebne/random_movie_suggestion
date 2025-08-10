FROM node:20-alpine AS css-builder

WORKDIR /app

COPY package.json package-lock.json ./
COPY web/styles ./web/styles
COPY tailwind.config.js ./

RUN npm install
RUN npx tailwindcss -i ./web/styles/styles.css -o ./styles.css

FROM golang:alpine3.20

WORKDIR /app

COPY . .

COPY --from=css-builder /app/styles.css ./web/static/css/styles.css

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /movie-generator

CMD ["/movie-generator"]