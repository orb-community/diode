# Diode quickstart

This is a basic set of instructions on how to get started using Diode on your local machine using Docker.

## Requirements

You'll need a recent installation of *Docker Community Edition* (or compatible) and *docker-compose*.

You may also need to have [`Git` installed](https://git-scm.com/downloads).

And finally you should be working on a Linux or macOS computer. Diode should technically also work on *Docker for Windows* (with Linux backend) or on the  *Windows Subsystem for Linux (WSL)*, but this guide does not cover this (yet).

## Diode configuration files

Diode requires two configuration files to execute successfully:

* `docker-compose.yml` - to configure the Diode containers
* `config.yaml` - to configure the discovery scope

We recommend placing both configuration files in the same directory and running all commands from this directory:

```bash
~ % cd ~
~ % mkdir diode
~ % cd diode
~/diode 
```

### Getting the default `docker-compose.yml`

You can get a working `docker-compose.yml` file by downloading an example from the Diode repository:

```bash
~/diode % curl https://raw.githubusercontent.com/orb-community/diode/develop/docker/docker-compose.yml -o docker-compose.yml
```

### Getting an example `config.yml`

You can get a working `config.yml` file by downloading an example from the Diode repository:

```bash
~/diode % curl https://raw.githubusercontent.com/orb-community/diode/develop/docker/config.yml -o config.yml
```

### Updating the `config.yml` for your discovery

The `config.yml` should look something like this and the `hosts:` section should be updated with a list of devices (and their credentials) to be discovered. 

```yaml
diode:
  config:
    output_type: otlp
    output_path: "127.0.0.1:4317"
  policies:  
    discovery_1:
      kind: discovery
      backend: suzieq
      data:   
        inventory: 
          sources:
            - name: default_inventory
              hosts:
                - url: ssh://192.168.0.4:2021 username=user password=password
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

The `inventory:` section of the `config.yml` follows the SuzieQ Inventory File Format. Please refer to the SuzieQ [documentation](https://suzieq.readthedocs.io/en/latest/inventory/) for additional details.

## Running Diode

Before running Diode, you should set the `NETBOX_API_HOST` and `NETBOX_API_TOKEN` environment variables to the send the discovery information to the correct NetBox instance.

```bash
~/diode % export NETBOX_API_HOST=my.netbox.instance:8000
~/diode % export NETBOX_API_TOKEN=123456789ABCDEF
```

You can now run Diode by executing the following command:

```bash
~/diode % docker compose up
```
