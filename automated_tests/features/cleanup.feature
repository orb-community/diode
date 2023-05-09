Feature: Cleanup test env

@cleanup
Scenario: cleanup yaml file
  Then remove all the agents .yaml generated on test process

@cleanup
Scenario: cleanup test containers
  Then force remove of all agent containers whose names start with the test prefix