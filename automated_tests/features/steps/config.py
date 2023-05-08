import configparser
from hamcrest import *
import os

LOCAL_AGENT_CONTAINER_NAME = "diode-agent-int-test"


class TestConfig:
    _configs = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def configs(cls):
        if cls._configs is None:
            cls._configs = _read_configs()
        return cls._configs


def _read_configs():
    parser = configparser.ConfigParser()
    parser.read('./features/config.ini')
    configs = parser['test_config']

    local_path = configs.get("local_path", os.getcwd())  # local_path is required if user will use docker to test,
    # otherwise the function will map the local path.
    assert_that(os.path.exists(local_path), equal_to(True), f"Invalid path: {local_path}.")
    configs['local_path'] = f"{local_path}/"

    agent_image_name = configs.get('agent_docker_image', 'orbcommunity/diode-agent')
    agent_image_tag = configs.get('agent_docker_tag', 'develop')
    agent_image = f"{agent_image_name}:{agent_image_tag}"
    configs['agent_image'] = agent_image

    return configs
