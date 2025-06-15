FROM golang:1.24-alpine AS builder

WORKDIR /build

RUN touch config.yaml
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build --ldflags="-s -w" -o /build/main main.go

FROM alpine:3.22

COPY --from=builder /build/main /app/main
COPY --from=builder /build/config.yaml /app/config.yaml

EXPOSE 8080
WORKDIR /app

ENTRYPOINT ["/app/main"]
