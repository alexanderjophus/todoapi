FROM golang:1.16-alpine as builder
WORKDIR /go/src/github.com/trelore/todoapi
COPY . .
RUN CGO_ENABLED=0 go build -mod vendor -o app -ldflags="-s -w" .

# RUN chmod +x app

FROM gcr.io/distroless/base
COPY --from=builder /go/src/github.com/trelore/todoapi/app /app
ENTRYPOINT ["./app"]