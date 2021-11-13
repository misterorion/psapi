FROM golang:1.17-alpine as builder
WORKDIR /app
COPY ./api/go.mod ./
COPY ./api/go.sum ./
RUN go mod download
COPY ./api/*.go .
RUN go build -o ./main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 80
CMD ["/app/main"]