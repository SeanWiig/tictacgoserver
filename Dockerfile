FROM golang:1.21 AS builder
WORKDIR /code
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY internal internal
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /code/server /app/server
ENTRYPOINT ["/app/server"]