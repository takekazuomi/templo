FROM golang:1.19 AS builder

WORKDIR /app
COPY . .
RUN go mod download && go build ./templo.go

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/templo /go/bin/templo

ENTRYPOINT ["/go/bin/templo"]
