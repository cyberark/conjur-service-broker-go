Feature: Deprovisioning

  Scenario: Deprovision resource with incorrect HTTP basic auth credentials
    Given my basic auth credentials are incorrect
    When I send "DELETE" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77?service_id=service-id-here&plan_id=plan-id-here"
    Then the response code should be 401
    And the response should match json:
      """
      { "error": "unauthorized" }
      """

  Scenario: Deprovision resource
    When I send "DELETE" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77?service_id=service-id-here&plan_id=plan-id-here"
    Then the response code should be 200
    And the response should match json "{}"
