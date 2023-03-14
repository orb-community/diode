FROM golang:1.19-alpine AS builder

WORKDIR /go/src/github.com/orb-community/diode
COPY go.mod .
RUN go mod tidy
COPY . .
RUN apk update && apk add make build-base git
RUN mkdir /tmp/build && CGO_ENABLED=1 make agent_bin && mv build/diode-agent /tmp/build/diode-agent

FROM python:3.8.16-slim-buster

WORKDIR /

RUN /usr/local/bin/python -m pip install --upgrade pip && pip install suzieq==0.20.1

COPY ./agent/ /tmp/agent/

RUN SQPATH=$(pip show suzieq |  sed -n 's/Location: \(.*\)/\1/p') \
&& cp -r "/tmp/agent/backend/suzieq/extension/poller/" "$SQPATH/suzieq/" \
&& cp -r "/tmp/agent/backend/suzieq/extension/db/" "$SQPATH/suzieq/" \
&& mkdir /root/.suzieq \
&& cp -r /tmp/agent/docker/suzieq-cfg.yml /root/.suzieq/suzieq-cfg.yml \
&& cp -r /tmp/agent/docker/inventory.yml /root/.suzieq/inventory.yml \
&& rm -rf /tmp/agent/

COPY --from=builder /tmp/build/diode-agent /usr/local/bin/diode-agent

ENTRYPOINT [ "bash" ]
