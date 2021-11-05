FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /go/bin/hello

FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/hello /go/bin/hello
# Run the hello binary.
ENTRYPOINT ["/go/bin/hello"]