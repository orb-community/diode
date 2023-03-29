# Diode quickstart

This is a basic set of instructions on how to get started using Diode on your local machine using Docker.

## Requirements

You'll need a recent installation of *Docker Community Edition* (or compatible) and *docker-compose*.

You may also need to have [`Git` installed](https://git-scm.com/downloads).

And finally you should be working on a Linux or macOS computer. Diode should technically also work on *Docker for Windows* (with Linux backend) or on the  *Windows Subsystem for Linux (WSL)*, but this guide does not cover this (yet).

## Diode configuration files

Diode requires two configuration files to execute successfully:

* `docker-compose.yml` - to configure and run the Diode containers
* `config.yaml` - to configure the scope of the discovery

We recommend placing both configuration files in the same directory and running all commands from this common directory. For example:

```bash
~ % cd ~
~ % mkdir diode
~ % cd diode
~/diode 
```

### Getting the default Diode `docker-compose.yml`

You can get the default Diode `docker-compose.yml` file by downloading this example from the Diode repository:

```bash
~/diode % curl https://raw.githubusercontent.com/orb-community/diode/develop/docker/docker-compose.yml -o docker-compose.yml
```

### Getting a template `config.yml`

You can get a template of the `config.yml` file by downloading this example from the Diode repository:

```bash
~/diode % curl https://raw.githubusercontent.com/orb-community/diode/develop/docker/config.yml -o config.yml
```

### Updating the `config.yml` for your discovery

The `config.yml` needs to be updated with an inventory of devices to be discovered. The file will look something like this, where the `hosts:` section needs to be populated with the list of devices and their credentials that you want to have discovered.

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
                - url: ssh://1.2.3.4:2021 username=user password=password
                - url: ssh://resolvable.host.name username=user password=password
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

Before running Diode, you should set the `NETBOX_API_HOST`, `NETBOX_API_TOKEN` and `NETBOX_API_PROTOCOL` (`http` or `https`) environment variables to send the discovery output to the correct NetBox instance.

```bash
~/diode % export NETBOX_API_HOST=my.netbox.instance:8000
~/diode % export NETBOX_API_TOKEN=123456789ABCDEF
~/diode % export NETBOX_API_PROTOCOL=http
```

You can now run Diode by executing the following command:

```bash
~/diode % docker compose up
```
