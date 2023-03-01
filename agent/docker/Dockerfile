FROM python:3.8.16-slim-buster

RUN /usr/local/bin/python -m pip install --upgrade pip && pip install suzieq==0.20.1

COPY backend/suzieq/extension /tmp/extension

ENTRYPOINT [ "bash" ]
