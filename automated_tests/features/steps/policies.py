from utils import generate_random_string_with_predefined_prefix, threading_wait_until, yaml_to_dict
from behave import step, when, then
from hamcrest import *
import requests
import yaml
from deepdiff import DeepDiff
from random import sample

POLICY_NAME_PREFIX = "test_diode_policy_"


@step("{amount_of_policies} policies are applied to the agent")
def create_and_apply_policies_to_agent(context, amount_of_policies):
    create_and_apply_policy_to_agent(context, amount_of_policies)


@step("{amount_of_policies} policy is applied to the agent")
def create_and_apply_policy_to_agent(context, amount_of_policies):
    assert_that(amount_of_policies.isnumeric(), equal_to(True), f"Amount of policies must be an integer. "
                                                                f"It was: {amount_of_policies}")
    for i in range(int(amount_of_policies)):
        policy_name = f"{generate_random_string_with_predefined_prefix(POLICY_NAME_PREFIX, 10)}"
        policy = Policy(policy_name)
        policy.full_test_policy()
        policy_yaml = policy.yaml()
        create_policy_response = create_policy(policy_yaml, context.port)
        ddiff = DeepDiff(yaml_to_dict(create_policy_response), policy.policy[policy_name])
        assert_that(ddiff, equal_to({}), f"Policy created is different of policy set in body request. {ddiff}")
        get_policy_response = get_policy(policy_name, context.port)
        ddiff = DeepDiff(yaml_to_dict(get_policy_response), policy.policy[policy_name])
        assert_that(ddiff, equal_to({}), f"Policy returned in /policies/{policy_name} is different of policy set in "
                                         f"body request. {ddiff}")
        context.policies[policy_name] = policy.policy[policy_name]


@when("{amount_of_policies} policy is deleted from agent")
def remove_policy_from_agent(context, amount_of_policies):
    amount_of_policies = int(amount_of_policies)
    assert_that(amount_of_policies, less_than_or_equal_to(len(list(context.policies.keys()))),
                f"Unable to remove {amount_of_policies} policies, "
                f"because only {len(list(context.policies.keys()))} "
                f"policies are applied")
    policies_to_remove = sample(list(context.policies.keys()), amount_of_policies)
    for policy_name in policies_to_remove:
        remove_policy(policy_name, context.port)
        context.policies.pop(policy_name)


@then("policies route shows {amount_of_policies} policies applied")
def check_amount_of_policies_applied_plural(context, amount_of_policies):
    check_amount_of_policies_applied(context, amount_of_policies)


@then("policies route shows {amount_of_policies} policy applied")
def check_amount_of_policies_applied(context, amount_of_policies):
    succeed, all_policies = wait_until_policies_applied(context.port, int(amount_of_policies))
    assert_that(len(all_policies), equal_to(int(amount_of_policies)),
                f"Unexpected amount of policies applied to the agent. Policies applied: {all_policies}")


class Policy:

    def __init__(self, name, kind="discovery", backend_type="suzieq"):
        self.policy = {name: {"kind": kind, "backend": backend_type, "config": {"netbox": {}},
                              "data": {"inventory": {"sources": [], "devices": [], "auths": [], "namespaces": []}}}}
        self.config = self.policy[name]['config']
        self.inventory_sources = self.policy[name]["data"]["inventory"]["sources"]
        self.inventory_devices = self.policy[name]["data"]["inventory"]["devices"]
        self.inventory_auths = self.policy[name]["data"]["inventory"]["auths"]
        self.inventory_namespaces = self.policy[name]["data"]["inventory"]["namespaces"]

    def full_test_policy(self):
        source_name = f"{generate_random_string_with_predefined_prefix(POLICY_NAME_PREFIX, 4)}_source"
        source = {"name": source_name, "hosts": [{"url": "ssh://test.diode.agent username=test1"},
                                                 {"url": "ssh://other.test.agent username=test2"}]}
        self.inventory_sources.append(source)

        device_name = f"{generate_random_string_with_predefined_prefix(POLICY_NAME_PREFIX, 4)}_devices"

        device = {"name": device_name, "transport": "ssh", "ignore-known-hosts": True, "slow-host": True}
        self.inventory_devices.append(device)

        auths = [{"name": "suzieq-u", "username": "test1", "password": "test@123"},
                 {"name": "suzieq-x", "username": "test2", "password": "12345678"}]
        for auth in auths:
            self.inventory_auths.append(auth)

        namespace = {"name": f"{generate_random_string_with_predefined_prefix(POLICY_NAME_PREFIX, 4)}_namespace",
                     "source": source_name, "device": device_name}
        self.inventory_namespaces.append(namespace)

        return self.policy

    def yaml(self):
        return yaml.dump(self.policy)


def create_policy(yaml_request, port, api_url="http://localhost", expected_status_code=201):
    """

    Creates a new policy in diode agent

    :param (dict) yaml_request: policy yaml
    :param (int) port: port where agent is running
    :param (str) api_url: policy api path
    :param (int) expected_status_code: code to be returned on response
    :return:  of policy creation

    """

    headers_request = {'Content-type': 'application/x-yaml'}

    response = requests.post(f"{api_url}:{port}/api/v1/policies", data=yaml_request, headers=headers_request)
    try:
        response_json = response.json()
    except ValueError:
        response_json = response.text
    assert_that(response.status_code, equal_to(expected_status_code),
                'Request to create policy failed with status=' + str(response.status_code) + ': ' + str(response_json))

    return response_json


def get_policies(port, api_url="http://localhost", expected_status_code=200):
    """

   Get all policies

    :param (int) port: port where agent is running
    :param (str) api_url: diode api path
    :param (int) expected_status_code: code to be returned on response
    :return:  all policies applied to agent

    """

    response = requests.get(f"{api_url}:{port}/api/v1/policies")
    try:
        response_json = response.json()
    except ValueError:
        response_json = response.text
    assert_that(response.status_code, equal_to(expected_status_code),
                'Request to get policies failed with status=' + str(response.status_code) + ': ' + str(response_json))

    return response_json


def get_policy(policy_name, port, api_url="http://localhost", expected_status_code=200):
    """

    Get policy information

    :param (str) policy_name: name of the policy to be fetched
    :param (int) port: port where agent is running
    :param (str) api_url: diode api path
    :param (int) expected_status_code: code to be returned on response
    :return:  policy response

    """

    response = requests.get(f"{api_url}:{port}/api/v1/policies/{policy_name}")
    try:
        response_json = response.json()
    except ValueError:
        response_json = response.text
    assert_that(response.status_code, equal_to(expected_status_code), f"Request to get policy {policy_name} "
                                                                      f"failed with status="
                                                                      f"{str(response.status_code)}:"
                                                                      f"{str(response_json)}")

    return response_json


def remove_policy(policy_name, port, api_url="http://localhost", expected_status_code=200):
    """

    Remove policy

    :param (str) policy_name: name of the policy to be fetched
    :param (int) port: port where agent is running
    :param (str) api_url: diode api path
    :param (int) expected_status_code: code to be returned on response
    :return:  policy response

    """

    response = requests.delete(f"{api_url}:{port}/api/v1/policies/{policy_name}")
    try:
        response_json = response.json()
    except ValueError:
        response_json = response.text
    assert_that(response.status_code, equal_to(expected_status_code), f"Request to delete policy {policy_name} "
                                                                      f"failed with status="
                                                                      f"{str(response.status_code)}:"
                                                                      f"{str(response_json)}")

    return response_json


@threading_wait_until
def wait_until_policies_applied(port, amount_of_policies, api_url="http://localhost", expected_status_code=200,
                                event=None):
    all_policies = get_policies(port, api_url, expected_status_code)
    if len(all_policies) == amount_of_policies:
        event.set()
    return event.is_set(), all_policies
