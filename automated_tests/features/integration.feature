@integration @AUTORETRY
Feature: Diode integration tests
  Suite to test against diode

  Scenario: run diode agent in default port
    Given that a diode configuration file exist with default configuration and port default
    When the diode agent is run using existing configuration file
    Then the diode agent container is running

  Scenario: run diode agent in non-default port
    Given that a diode configuration file exist with default configuration and port available
    When the diode agent is run using existing configuration file
    Then the diode agent container is running



  Scenario: apply one policy to diode agent

  Scenario: remove one policy to diode agent