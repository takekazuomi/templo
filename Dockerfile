FROM golang:1.15 AS builder
WORKDIR /tmp
RUN go get -v github.com/takekazuomi/templo

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/bin/templo /go/bin/templo

ENTRYPOINT ["/go/bin/templo"]