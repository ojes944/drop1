FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy

FROM golang:1.21-alpine
WORKDIR /app
COPY --from=builder /app /app
EXPOSE 8080
CMD ["go", "run", "main.go"]
