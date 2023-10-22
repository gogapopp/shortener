FROM golang:1.21.0-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main ./cmd/shortener/main.go
EXPOSE 8080
CMD [ "/app/main", "-d", "host=db user=postgres database=postgres port=5432 password=123" ]
