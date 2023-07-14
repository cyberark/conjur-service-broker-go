Feature: Deprovisioning

  Scenario: Deprovision resource with incorrect HTTP basic auth credentials
    Given my basic auth credentials are incorrect
    When I send "DELETE" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77?service_id=c024e536-6dc4-45c6-8a53-127e7f8275ab&plan_id=3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
    Then the response code should be 401
    And the response should match json:
      """
      { "error": "Unauthorized" }
      """

  Scenario: Deprovision resource with incorrect plan id
    When I send "DELETE" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77?service_id=c024e536-6dc4-45c6-8a53-127e7f8275ab&plan_id=invalid"
    Then the response code should be 400
    And the response should match json:
      """
      { "error": "ValidationError", "description": "parameter \"plan_id\" in query has an error: value is not one of the allowed values [\"3a116ac2-fc8b-496f-a715-e9a1b205d05c.community\"]" }
      """

  Scenario: Deprovision resource
    When I send "DELETE" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77?service_id=c024e536-6dc4-45c6-8a53-127e7f8275ab&plan_id=3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
    Then the response code should be 200
    And the response should match json "{}"
