import docker
import threading
from behave import given, when, then
from diode_config_file import Diode, default_path_config_file, AGENT_FILE_NAME_PREFIX
from hamcrest import *
from utils import return_port_by_availability, threading_wait_until
from config import TestConfig

configs = TestConfig.configs()


@given("that a diode configuration file exist with default configuration and port {port}")
def create_diode_config_file(context, port):
    assert_that(port, any_of("available", "unavailable", "default"), "Invalid option for port")

    if port == "default":
        agent, context.agent_file_name, file_path = Diode.create_config_file()
        context.port_to_use = 10911
    else:
        availability = {"available": True, "unavailable": False}
        context.port_availability = availability[port]
        context.port_to_use = return_port_by_availability(context.port_availability)
        agent, context.agent_file_name, file_path = Diode.create_config_file(config={"port": context.port_to_use})


@when('the diode agent is run using existing configuration file')
def run_agent_using_config_file(context):
    agent_image = configs.get("agent_image")
    dir_path = configs.get("local_path")
    file_path = f"{dir_path}"
    volume = [f"{file_path}:{default_path_config_file}"]
    command = f"run -c {default_path_config_file}{context.agent_file_name}.yaml"
    context.container_id = run_agent_container(agent_image, context.agent_file_name, command=command, volumes=volume)
    assert_that(context.container_id, is_not(None), "Failed to run agent container. Container id is None.")
    context.containers_id_port.update({context.container_id: context.port_to_use})


@then('the diode agent container is {status}')
def validate_container_status(context, status):
    expected_status, container = wait_container_status(context.container_id, status)
    assert_that(expected_status, is_(True), f"Container status fail with: {container.status}. Expected was: {status}")


@then("force remove of all agent containers whose names start with the test prefix")
def remove_all_orb_agent_test_containers(context):
    docker_client = docker.from_env()
    containers = docker_client.containers.list(all=True)
    for container in containers:
        test_container = container.name.startswith(AGENT_FILE_NAME_PREFIX)
        if test_container is True:
            container.remove(force=True)


@then("stop agent container")
def stop_orb_agent_container(context):
    for container_id in context.containers_id_port.keys():
        stop_container(container_id)


@then("remove agent container")
def remove_orb_agent_container(context):
    for container_id in context.containers_id_port.keys():
        remove_container(container_id)
    context.containers_id = {}


def run_agent_container(container_image, container_name, env_vars=None, detach=True, command=None, network_mode="host",
                        volumes=None, time_to_wait=5):
    """
    Run agent container

    :param (str) container_image: that will be used for running the container
    :param (str) container_name: base of container name
    :param (dict) env_vars: that will be passed to the container context
    :param (bool) detach: If true, detach from the exec command. Default: True
    :param (str or list) command: The command to run in the container. Default: None
    :param (dtr) network_mode: One of:
                                    bridge: Create a new network stack for the container on the bridge network.
                                    none: No networking for this container.
                                    container:<name|id> Reuse another container’s network stack.
                                    host: Use the host network stack.
                                    Default:'host'
    :param (list) volumes: list of strings which each one of its elements specifies a mount volume. Default: None
    :param (int) time_to_wait: seconds that threading must wait after run the agent. Default: 5
    :returns: (str) the container ID
    """
    if volumes is None:
        volumes = list()
    if volumes is None:
        volumes = []
    client = docker.from_env()
    container = client.containers.run(container_image, environment=env_vars, name=container_name, detach=detach,
                                      command=command, network_mode=network_mode, volumes=volumes)
    threading.Event().wait(time_to_wait)
    return container.id


def stop_container(container_id):
    """

    :param container_id: agent container ID
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    container.stop()


def remove_container(container_id, force_remove=False):
    """

    :param container_id: agent container ID
    :param force_remove: if True, similar to docker rm -f. Default: False
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    container.remove(force=force_remove)


@threading_wait_until
def wait_container_status(container_id, status, event=None):
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    if container.status == status:
        event.set()
    return event.is_set(), container
