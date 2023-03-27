FROM golang:1.19-alpine AS builder
ARG SVC
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/orb-community/diode
COPY . .
RUN apk update \
    && apk add make
RUN make $SVC \
    && mv build/diode-$SVC /exe

FROM alpine:latest

RUN mkdir /opt/diode/
# Certificates are needed so that mailing util can work.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
