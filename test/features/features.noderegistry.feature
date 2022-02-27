Feature: Registry Events
  Scenario: Node registry scenario
    Given Node registry event is dispatched
    Then Node registry record is created

  Scenario: Peer registry scenario
    Given Peer registry event is dispatched
    Then Peer registry record is created
