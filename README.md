# Diode quickstart

This is a basic set of instructions on how to get started using Diode on your local machine using Docker.

## Requirements

### Tools

You'll need a recent installation of *Docker Community Edition* (or compatible) and *docker-compose*.

You may also need to have [`Git` installed](https://git-scm.com/downloads).

And finally you should be working on a Linux or macOS computer. Diode should technically also work on *Docker for Windows* (with Linux backend) or on the  *Windows Subsystem for Linux (WSL)*, but this guide does not cover this (yet).

> ❗️ Warning for M1 Mac users
> 
> Discovery currently does not work on M1 devices. We're working on a fix for this.

> ❗️ Docker and Docker Compose versions
> 
> It is better to update your Docker and Docker Compose to the latest versions before testing. Also note that on some systems the ``docker-compose`` and ``docker compose`` commands can report different versions so double check which one you're using!

### NetBox Instance

You will also need an instance of NetBox that the discovery results can be pushed into. The easiest way to do this if you don't already have a NetBox instance running is to use [NetBox Docker](https://github.com/netbox-community/netbox-docker).

You'll clone the ``netbox-docker`` repo, and make some small changes to the repo to get everything set-up:

```bash
git clone https://github.com/netbox-community/netbox-docker.git
cd netbox-docker
mv docker-compose.override.yml.example docker-compose.override.yml
```

Now edit ``docker-compose.override.yml`` to include your desired super user settings to look something like this:

```
version: '3.7'
services:
  netbox:
    ports:
      - "8000:8080"
    # If you want the Nginx unit status page visible from the
    # outside of the container add the following port mapping:
    # - "8001:8081"
    # healthcheck:
      # Time for which the health check can fail after the container is started.
      # This depends mostly on the performance of your database. On the first start,
      # when all tables need to be created the start_period should be higher than on
      # subsequent starts. For the first start after major version upgrades of NetBox
      # the start_period might also need to be set higher.
      # Default value in our docker-compose.yml is 60s
      # start_period: 90s
    environment:
      SKIP_SUPERUSER: "false"
      SUPERUSER_API_TOKEN: "YOUR API TOKEN"
      SUPERUSER_EMAIL: "YOUR EMAIL"
      SUPERUSER_NAME: "YOUR SUPERUSER_USERNAME"
      SUPERUSER_PASSWORD: "YOUR SUPERUSER_PASSWORD"
```

You can now start ``netbox-docker`` using ``docker-compose``. Note that you'll need to wait for the NetBox service to become healthy which you can monitor with ``docker-compose ps``. Once healthy you can log in to your NetBox instance on ``127.0.0.1:8000`` using the super user credentials you entered above.

```bash
docker-compose up
```

## Diode configuration files

Diode requires two configuration files to execute successfully:

* `docker-compose.yml` - to configure and run the Diode containers
* `config.yaml` - to configure the scope of the discovery

We recommend placing both configuration files in the same directory and running all commands from this common directory. For example:

```bash
cd ~
mkdir diode
cd diode
```

### Getting the default Diode `docker-compose.yml`

You can get the default Diode `docker-compose.yml` file by downloading this example from the Diode repository:

```bash
curl https://raw.githubusercontent.com/orb-community/diode/develop/docker/docker-compose.yml -o docker-compose.yml
```

### Getting a template `config.yml`

You can get a template of the `config.yml` file by downloading this example from the Diode repository:

```bash
curl https://raw.githubusercontent.com/orb-community/diode/develop/docker/config.yml -o config.yml
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
export NETBOX_API_HOST=my.netbox.instance:8000
export NETBOX_API_TOKEN=123456789ABCDEF
export NETBOX_API_PROTOCOL=http
export TAG=develop
```

You can now run Diode by executing the following command:

```bash
docker compose up
```
