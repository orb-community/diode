# Diode

Diode is part of the NetBoxLabs suite of network monitoring and assurance tools, it provides a comprehensive set of features and capabilities, including network discovery, monitoring, and assurance. With Diode, users can easily discover and monitor their networks, set thresholds and alerting levels, and detect devices that are not compliant with network policies. The data collected by Diode can be used to generate insights and reports on the network state and performance, and the data can be integrated into other applications via the API. Diode is designed to maximize network performance and provide network administrators with the best possible visibility and control over their networks. With Diode, users can rest assured that their networks are secure and running optimally.

The <b>diode-agent</b> is responsible for collecting data from the network and transmitting it to the <b>diode-server</b>. This data can be used to generate insights and reports on the network state and performance. The <b>diode-server</b> will then store the data and make it available to the users of the <b>netbox</b> application. 

The <b>diode-agent</b> will be able to scan the network for devices and their configuration. It will also be able to monitor the network for changes and alert the users if any devices or network parameters fall outside of the configured thresholds. The agent will also be able to detect devices that are not compliant with the network policies and alert the users. 

The <b>diode-server</b> will be responsible for collecting, storing and delivering data from the diode-agent to the users of the netbox application. It will also provide a web interface for users to configure the diode-agent and set the thresholds for the network parameters. The <b>diode-server</b> also provide an way for users to integrate the <b>diode-agent</b> data into their netbox or netbox cloud instances.


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
