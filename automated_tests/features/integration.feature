@integration @AUTORETRY
Feature: Diode integration tests

  @smoke
  Scenario: run diode agent in default port
    Given that a diode configuration file exist with default configuration and port default
    When the diode agent is run using existing configuration file
    Then the diode agent container is running

  @smoke
  Scenario: run diode agent in non-default port
    Given that a diode configuration file exist with default configuration and port available
    When the diode agent is run using existing configuration file
    Then the diode agent container is running

  @smoke
  Scenario: run diode agent in unavailable port
    Given that a diode agent is already running
      And that a diode configuration file exist with default configuration and port unavailable
    When the diode agent is run using existing configuration file
    Then the diode agent container is exited

  @smoke
  Scenario: apply one policy to diode agent
    Given that a diode agent is already running
    When 1 policy is applied to the agent
    Then policies route shows 1 policy applied

  @smoke
  Scenario: apply multiple policies to diode agent
    Given that a diode agent is already running
    When 30 policies are applied to the agent
    Then policies route shows 30 policies applied

  @smoke
  Scenario: remove one policy from diode agent
    Given that a diode agent is already running
      And 3 policies are applied to the agent
    When 1 policy is deleted from agent
    Then policies route shows 2 policies applied

  @smoke
  Scenario: remove all policies from diode agent
    Given that a diode agent is already running
      And 1 policies are applied to the agent
    When 1 policy is deleted from agent
    Then policies route shows 0 policies applied