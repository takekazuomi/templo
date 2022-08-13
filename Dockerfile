FROM golang:1.19 AS builder

WORKDIR /app
COPY . .

RUN go mod download \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./templo.go

FROM scratch

WORKDIR /app/
COPY --from=builder /app/templo /go/bin/templo

ENTRYPOINT ["/go/bin/templo"]
