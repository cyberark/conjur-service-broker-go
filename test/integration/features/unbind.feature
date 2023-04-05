Feature: Unbinding

  Scenario: Unbind with incorrect HTTP basic auth credentials
    Given I send "PUT" request to "/v2/service_instances/6b40649e-331b-424d-afa0-6d569f016f51/service_bindings/5e7a43f2-b3fc-4591-ab19-783389c2cb63" with body:
    """
    {
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "bind_resource": {
        "app_guid": "5e7a43f2-b3fc-4591-ab19-783389c2cb63"
      },
      "parameters": {
        "parameter1-name-here": 1,
        "parameter2-name-here": "parameter2-value-here"
      }
    }
    """
    And my basic auth credentials are incorrect
    When I send "DELETE" request to "/v2/service_instances/6b40649e-331b-424d-afa0-6d569f016f51/service_bindings/5e7a43f2-b3fc-4591-ab19-783389c2cb63?service_id=c024e536-6dc4-45c6-8a53-127e7f8275ab&plan_id=3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
    Then the response code should be 401
    And the response should match json:
      """
      { "error": "Unauthorized" }
      """

  Scenario: Unbind where binding does not exist
    When I send "DELETE" request to "/v2/service_instances/fe837829-2174-4c7a-8686-d3635e38b145/service_bindings/e765e6d3-3264-417f-a726-172e3c364911?service_id=c024e536-6dc4-45c6-8a53-127e7f8275ab&plan_id=3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
    Then the response code should be 410
    And the response should match json:
      """
      { "error": "Gone", "description": "host doesn't exists" }
      """

  Scenario: Successful unbinding
    Given I send "PUT" request to "/v2/service_instances/6b40649e-331b-424d-afa0-6d569f016f51/service_bindings/5e7a43f2-b3fc-4591-ab19-783389c2cb64" with body:
    """
    {
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "bind_resource": {
        "app_guid": "5e7a43f2-b3fc-4591-ab19-783389c2cb64"
      },
      "parameters": {
        "parameter1-name-here": 1,
        "parameter2-name-here": "parameter2-value-here"
      }
    }
    """
    And the response code should be 201
    And I create conjur client
    And conjur credentials are valid
    When I send "DELETE" request to "/v2/service_instances/6b40649e-331b-424d-afa0-6d569f016f51/service_bindings/5e7a43f2-b3fc-4591-ab19-783389c2cb64?service_id=c024e536-6dc4-45c6-8a53-127e7f8275ab&plan_id=3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
    And the response code should be 200
    And conjur credentials are invalid

#  Scenario: Unbind with incorrect Conjur credentials
#    Given I make a bind request with body:
#    """
#    {
#      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
#      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
#      "bind_resource": {
#        "app_guid": "5e7a43f2-b3fc-4591-ab19-783389c2cb63"
#      },
#      "parameters": {
#        "parameter1-name-here": 1,
#        "parameter2-name-here": "parameter2-value-here"
#      }
#    }
#    """
#    And I use a service broker with a bad Conjur API key
#    When I make a corresponding unbind request
#    Then the HTTP response status code is "403"
#    And the JSON should be {}
#
#  Scenario: Unbind with Conjur unavailable
#    Given I make a bind request with body:
#    """
#    {
#      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
#      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
#      "bind_resource": {
#        "app_guid": "5e7a43f2-b3fc-4591-ab19-783389c2cb63"
#      },
#      "parameters": {
#        "parameter1-name-here": 1,
#        "parameter2-name-here": "parameter2-value-here"
#      }
#    }
#    """
#    And I use a service broker with a bad Conjur URL
#    When I make a corresponding unbind request
#    Then the HTTP response status code is "500"
#    And the JSON should be {}
