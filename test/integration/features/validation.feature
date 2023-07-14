# file: validation.feature
Feature: request and response validation
  In order use service broker
  As an CF compliant platform
  I need to be assure contract is obeyed

  Scenario: non existing endpoint
    Given I ignore route error
    When I send "POST" request to "/v2/non-existing"
    Then the response code should be 404
    And the response should match json:
      """
      {
        "error": "NotFound"
      }
      """

  Scenario: does not allow POST method on catalog
    Given I ignore route error
    When I send "POST" request to "/v2/catalog"
    Then the response code should be 405
    And the response should match json:
      """
      {
        "error": "MethodNotAllowed"
      }
      """
