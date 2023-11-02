FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o chargen ./cmd/chargen/main.go
EXPOSE 8080
CMD ["./chargen"]
