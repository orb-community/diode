from behave.model_core import Status


def before_scenario(context, scenario):
    context.containers_id_port = dict()


def after_scenario(context, scenario):
    if scenario.status != Status.failed:
        context.execute_steps('''
        Then stop agent container
        Then remove agent container
        ''')
