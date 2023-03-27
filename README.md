# Diode

## Building diode service docker image:

`SERVICE=service make build_docker`

## Running diode service docker image

This an example [docker-compose.yml](https://github.com/orb-community/diode/blob/develop/docker/docker-compose.yml) to run diode service [image](https://hub.docker.com/r/orbcommunity/diode-service/tags):
```
version: '3.7'

# Docker Services
services:

  service:
    image: orbcommunity/diode-service:IMAGE_TAG
    ports:
      - "4317:4317"
    environment:
      - DIODE_SERVICE_NETBOX_ENDPOINT=NETBOX_API_HOST
      - DIODE_SERVICE_NETBOX_TOKEN=NETBOX_API_TOKEN
```

You need to change some items like IMAGE_TAG, NETBOX_API_HOST, NETBOX_API_TOKEN to adapt your environment

## Building diode agent docker image:

`make agent`

## Running diode agent docker image

```sh
docker run -v /usr/local/diode:/opt/diode/ --net=host orbcommunity/diode-agent:develop run -c /opt/diode/config.yaml
```
Diode agent can be configured as different output formats:

- FILE 
The agent output will be on a one or more files that will be created on the folder especified in output_path option

- HTTP
The agent outputt will be on a POST directly on Netbox Ingest plug-in API, that uses output_path as API endpoint and output_auth as API token

- OTLP
The agent output will be on a OTLP request to the diode service port on the endpoint especified in output_path option

## Config file example using FILE output
```yaml
diode:
  config:
    output_type: file
    output_path: "/opt/diode"
  policies:      
    discovery_1:
      kind: discovery
      backend: suzieq
      data:
        config:
          poller:
            run-once-update: true       
        inventory: 
          sources:
            - name: default_inventory
              hosts:
                - url: ssh://192.168.0.4 username=user password=password
          devices:
            - name: default_devices
              transport: ssh
              ignore-known-hosts: true
              slow-host: true
          namespaces:
            - name: default_namespace
              source: default_inventory
              device: default_devices
```

## Config file example using HTTP output
```yaml
diode:
  config:
    output_type: http
    output_path: "https://your-netbox-url/api/plugins/ingest/ingest/"
    output_auth: "Token xxxxxxxxxxxxxx-your-netbox-token-xxxxxxxxxxxxxxxx"
  policies:      
    discovery_1:
      kind: discovery
      backend: suzieq
      data:
        config:
          poller:
            run-once-update: true       
        inventory: 
          sources:
            - name: default_inventory
              hosts:
                - url: ssh://192.168.0.4 username=user password=password
          devices:
            - name: default_devices
              transport: ssh
              ignore-known-hosts: true
              slow-host: true
          namespaces:
            - name: default_namespace
              source: default_inventory
              device: default_devices
```


## Config file example using OTLP output
```yaml
diode:
  config:
    output_type: otlp
    output_path: "your-diode-service-instance:4317"
  policies:      
    discovery_1:
      kind: discovery
      backend: suzieq
      data:
        config:
          poller:
            run-once-update: true       
        inventory: 
          sources:
            - name: default_inventory
              hosts:
                - url: ssh://192.168.0.4 username=user password=password
          devices:
            - name: default_devices
              transport: ssh
              ignore-known-hosts: true
              slow-host: true
          namespaces:
            - name: default_namespace
              source: default_inventory
              device: default_devices
```

