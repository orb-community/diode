import yaml
from config import TestConfig
from utils import generate_random_string_with_predefined_prefix, find_files
import os
from behave import then

AGENT_FILE_NAME_PREFIX = "test_diode_agent_"
default_path_config_file = "/opt/diode/"
configs = TestConfig.configs()


@then("remove all the agents .yaml generated on test process")
def remove_agent_config_files(context):
    dir_path = configs.get("local_path")
    all_files_generated = find_files(AGENT_FILE_NAME_PREFIX, ".yaml", dir_path)
    if len(all_files_generated) > 0:
        for file in all_files_generated:
            os.remove(file)


class Diode:
    # def __init__(self):
    #     pass

    @classmethod
    def create_config_file(cls, output_type='otlp', output_path_host="0.0.0.0", output_path_port=4317, **kwargs):
        agent = {
            "diode": {"config": {"output_type": output_type, "output_path": f"{output_path_host}:{output_path_port}"}}}
        for key, value in kwargs.items():
            agent['diode'][key].update(value)
        agent = yaml.dump(agent)

        agent_file_name = generate_random_string_with_predefined_prefix(AGENT_FILE_NAME_PREFIX)
        dir_path = configs.get("local_path")
        file_path = f"{dir_path}{agent_file_name}.yaml"
        with open(file_path, "w+") as f:
            f.write(agent)
        return agent, agent_file_name, file_path
