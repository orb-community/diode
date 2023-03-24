# Diode

## Building diode service docker image:

`SERVICE=service make build_docker`


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
