FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /go/bin/app

FROM golang:alpine

COPY --from=builder /go/bin/app /go/bin/app

RUN apk add --no-cache ca-certificates openssl

ENTRYPOINT ["/go/bin/app"]