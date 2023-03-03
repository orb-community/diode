FROM python:3.8.16-slim-buster

WORKDIR /

RUN /usr/local/bin/python -m pip install --upgrade pip && pip install suzieq==0.20.1

COPY ./agent/ /tmp/agent/

RUN SQPATH=$(pip show suzieq |  sed -n 's/Location: \(.*\)/\1/p') \
&& cp -r "/tmp/agent/backend/suzieq/extension/poller/" "$SQPATH/suzieq/" \
&& cp -r "/tmp/agent/backend/suzieq/extension/db/" "$SQPATH/suzieq/" \
&& cp -r "/tmp/agent/docker/suzieq-cfg.yml" "/root/.suzieq/suzieq-cfg.yml" \
&& cp -r "/tmp/agent/docker/inventory.yml" "/root/.suzieq/inventory.yml" \
&& rm -rf /tmp/agent/

ENTRYPOINT [ "bash" ]
