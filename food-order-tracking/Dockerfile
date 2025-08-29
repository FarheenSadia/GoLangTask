# Stage 1: Build the Go binary
FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o food-app .

# Stage 2: Minimal image for running
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/food-app .
EXPOSE 8080
CMD ["./food-app"]
