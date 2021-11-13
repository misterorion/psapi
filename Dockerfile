ARG GO_IMAGE
ARG ALPINE_IMAGE
FROM $GO_IMAGE as builder
WORKDIR /app
COPY ./api/go.mod ./
COPY ./api/go.sum ./
RUN go mod download
COPY ./api/*.go .
RUN go build -o ./main

FROM $ALPINE_IMAGE
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 80
CMD ["/app/main"]