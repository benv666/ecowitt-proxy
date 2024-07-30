FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
# Or -ldflags="-extldflags=-static"
ENV CGO_ENABLED=0
RUN go build -o ecoproxy ./cmd/main.go

FROM gcr.io/distroless/static

COPY --from=builder /app/ecoproxy /

ENTRYPOINT ["/ecoproxy"]
