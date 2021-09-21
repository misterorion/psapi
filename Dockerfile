FROM golang:1.17-alpine as builder
WORKDIR /app
COPY src/go.mod ./
COPY src/go.sum ./
RUN go mod download
COPY src/*.go .
RUN go build -o ./main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["/app/main"]