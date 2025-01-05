FROM golang:latest

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8080

ENTRYPOINT ["air"]
