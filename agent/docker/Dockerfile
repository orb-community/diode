FROM python:3.8.16-slim-buster

RUN /usr/local/bin/python -m pip install --upgrade pip && pip install suzieq==0.20.1

WORKDIR /tmp

ADD ./agent/ /tmp/agent/

RUN SQPATH=$(pip show suzieq |  sed -n 's/Location: \(.*\)/\1/p') \
&& cp -r "/tmp/agent/backend/suzieq/extension/poller/" "$SQPATH/suzieq/" \
&& cp -r "/tmp/agent/backend/suzieq/extension/db/" "$SQPATH/suzieq/"

ENTRYPOINT [ "bash" ]
