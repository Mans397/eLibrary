FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./main.go

RUN ls -la /app

EXPOSE 8080

CMD ["./main"]