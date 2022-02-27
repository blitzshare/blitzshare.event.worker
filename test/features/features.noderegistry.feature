Feature: Registry Events

  Scenario: Node registry scenario
    Given Node registry event is dispatched
    Then Node registry record is created

  Scenario: Peer registry and deregistry scenario
    Given Peer registry event is dispatched
    And Peer registry record is created
    When Peer deregistry event is dispatched
    And Test Wait for 1 seconds
    Then Peer registry record is deleted


