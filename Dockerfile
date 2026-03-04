FROM golang:1.24-alpine as builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
# This creates a file named 'main'
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
# This copies that 'main' file over
COPY --from=builder /app/main .

# 2. COPY YOUR STATIC FILES (This fixes the 404)
COPY --from=builder /app/web-client ./web-client

EXPOSE 8080
# This runs it directly (no 'go' command needed)
CMD ["./main"]
