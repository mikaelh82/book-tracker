FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o book-tracker ./main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/book-tracker .
EXPOSE 8080
CMD ["./book-tracker"]