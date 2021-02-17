FROM golang:1.15 AS builder
WORKDIR /tmp
RUN go get -v github.com/takekazuomi/templo

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/bin/templo .
CMD ["./templo"]