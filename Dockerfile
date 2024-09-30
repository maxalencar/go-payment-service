FROM golang:1.23-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN /usr/local/go/bin/go build -o app ./cmd/app/

FROM alpine

COPY --from=builder /go/src/app/app /usr/local/bin/app
WORKDIR /usr/local/bin

CMD [ "app" ]