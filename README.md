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


## Config file example 
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
