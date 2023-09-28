#FROM golang:1.21.1 AS builder
FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

#COPY ca-bundle.crt /etc/ssl/certs/ca-bundle.crt
#COPY ca-bundle.trust.crt /etc/ssl/certs/ca-bundle.trust.crt
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

RUN go get github.com/mattn/go-sqlite3
#RUN
RUN apk add build-base
RUN CGO_ENABLED=1 go build -o app url-shortener/cmd/app

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /build/app

COPY /storage/storage.db /build/storage/storage.db
COPY config/dev.yaml /build/dev.yaml
ENV CONFIG_PATH dev.yaml

CMD ["./app"]