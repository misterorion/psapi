FROM golang:alpine as builder
WORKDIR /app
COPY ./api/go.mod ./
COPY ./api/go.sum ./
RUN go mod download
COPY ./api/*.go .
RUN go build -o ./main

FROM alpine:3
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["/app/main"]