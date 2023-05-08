from hamcrest import *
from random import choice, choices
import socket
import threading
import string
from datetime import datetime
import os


def return_port_by_availability(available=True, unavailable_ports=list(), time_to_wait=5):
    """

    :param (bool) available: Is the port on which agent must try to run available?. Default: True.
    :param (list) unavailable_ports: list of ports that are already in use, so not available to bind. Default: None.
    :param (int) time_to_wait: seconds that threading must wait after run the agent
    :return: (int) port number
    """

    assert_that(available, any_of(equal_to(True), equal_to(False)), "Unexpected value for 'available' parameter")

    if not available:
        assert_that(len(unavailable_ports), greater_than(0), "No port is unavailable")
        unavailable_port = choice(unavailable_ports)
        return unavailable_port
    else:
        available_port = None
        retries = 0
        while available_port is None and retries < 10:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.bind(('', 0))
            addr = s.getsockname()
            s.close()
            retries += 1
            if addr[1] not in unavailable_ports:
                available_port = addr[1]
                return available_port
    assert_that(available_port, is_not(None), "Unable to find an available port")
    threading.Event().wait(time_to_wait)
    return available_port


def generate_random_string_with_predefined_prefix(string_prefix, n_random=10):
    """
    :param (str) string_prefix: prefix to identify object created by tests
    :param (int) n_random: amount of random characters
    :return: random_string_with_predefined_prefix
    """
    random_string_with_predefined_prefix = string_prefix + random_string(n_random)
    return random_string_with_predefined_prefix


def random_string(k=10, mode='mixed'):
    """
    Generates a string composed of k (int) random letters lowercase and uppercase mixed

    :param (int) k: sets the length of the randomly generated string. Default:10
    :param(str) mode: define if the letters will be lowercase, uppercase or mixed.Default: mixed. Options: lower and
    upper.
    :return: (str) string consisting of k random letters lowercase and uppercase mixed.
    """
    assert_that(mode, any_of("mixed", "upper", "lower"), "Invalid string mode")
    if mode == 'mixed':
        return ''.join(choices(string.ascii_letters, k=k))
    elif mode == 'lower':
        return ''.join(choices(string.ascii_lowercase, k=k))
    else:
        return ''.join(choices(string.ascii_uppercase, k=k))


def threading_wait_until(func):
    def wait_event(*args, wait_time=1, timeout=30, start_func_value=False, **kwargs):
        event = threading.Event()
        func_value = start_func_value
        start = datetime.now().timestamp()
        time_running = 0
        while not event.is_set() and time_running < int(timeout):
            func_value = func(*args, event=event, **kwargs)
            event.wait(wait_time)
            time_running = datetime.now().timestamp() - start
        return func_value

    return wait_event


def find_files(prefix, suffix, path):
    """
    Find files that match with prefix and suffix condition

    :param prefix: string with which the file should start with
    :param suffix: string with which the file should end with
    :param path: directory where the files will be searched
    :return: (list) path to all files that matches
    """
    result = list()
    for root, dirs, files in os.walk(path):
        for name in files:
            if name.startswith(prefix) and name.endswith(suffix):
                result.append(os.path.join(root, name))
    return result
