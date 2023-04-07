Feature: Provisioning

  Scenario: Provision resource with incorrect HTTP basic auth credentials
    Given my basic auth credentials are incorrect
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77" with body:
    """
    {
    }
    """
    Then the response code should be 401
    And the response should match json:
      """
      { "error": "Unauthorized" }
      """

  Scenario: Provision resource with invalid body - missing keys
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a78" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "not_plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
      }
    }
    """
    Then the response code should be 400
    And the response should match json:
      """
      {
        "error": "ValidationError",
        "description": "request body has an error: doesn't match schema #/components/schemas/ServiceInstanceProvisionRequestBody: property \"plan_id\" is missing"
      }
      """


  Scenario: Provision resource with service broker API 2.15 or better
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
        "organization_name": "my-organization",
        "space_name": "my-space"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
      }
    }
    """
    Then the response code should be 201
    And the response should match json "{}"

  Scenario: Provision resource
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
      }
    }
    """
    Then the response code should be 201
    And the response should match json "{}"

  Scenario: Update resource
    When I send "PATCH" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "parameters": {
      },
      "previous_values": {
        "plan_id": "we-only-have-one-plan"
      }
    }
    """
    Then the response code should be 200
    And the response should match json "{}"

  Scenario: Update resource with invalid body - missing keys
    When I send "PATCH" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "not_service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "parameters": {
      },
      "previous_values": {
        "plan_id": "we-only-have-one-plan"
      }
    }
    """
    Then the response code should be 400
    And the response should match json:
      """
      {
        "error": "ValidationError",
        "description": "request body has an error: doesn't match schema #/components/schemas/ServiceInstanceUpdateRequestBody: property \"service_id\" is missing"
      }
      """

  Scenario: Update resource with invalid body - invalid service ID
    When I send "PATCH" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "XXXXXXX-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "parameters": {
      },
      "previous_values": {
        "plan_id": "we-only-have-one-plan"
      }
    }
    """
    Then the response code should be 400
    And the response should match json:
      """
      {
        "error": "ValidationError",
        "description": "request body has an error: doesn't match schema #/components/schemas/ServiceInstanceUpdateRequestBody: value is not one of the allowed values [\"c024e536-6dc4-45c6-8a53-127e7f8275ab\"]"
      }
      """

  Scenario: Update resource with invalid body - invalid plan ID
    When I send "PATCH" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a77" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "XXXXXXX-fc8b-496f-a715-e9a1b205d05c.community",
      "parameters": {
      },
      "previous_values": {
        "plan_id": "we-only-have-one-plan"
      }
    }
    """
    Then the response code should be 400
    And the response should match json:
      """
      {
        "error": "ValidationError",
        "description": "request body has an error: doesn't match schema #/components/schemas/ServiceInstanceUpdateRequestBody: value is not one of the allowed values [\"3a116ac2-fc8b-496f-a715-e9a1b205d05c.community\"]"
      }
      """

  Scenario: Provision resource with invalid body - invalid service ID
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a78" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "XXXXXXX-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
      }
    }
    """
    Then the response code should be 400
    And the response should match json:
      """
      {
        "error": "ValidationError",
        "description": "request body has an error: doesn't match schema #/components/schemas/ServiceInstanceProvisionRequestBody: value is not one of the allowed values [\"c024e536-6dc4-45c6-8a53-127e7f8275ab\"]"
      }
      """

  Scenario: Provision resource with invalid body - invalid plan ID
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a78" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "XXXXXXX-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
      }
    }
    """
    Then the response code should be 400
    And the response should match json:
    """
    {
      "error": "ValidationError",
      "description": "request body has an error: doesn't match schema #/components/schemas/ServiceInstanceProvisionRequestBody: value is not one of the allowed values [\"3a116ac2-fc8b-496f-a715-e9a1b205d05c.community\"]"
    }
    """

  Scenario: Provision resource with invalid parameters
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
        "parameter1": 1,
        "parameter2": "foo"
      }
    }
    """
    Then the response code should be 400
    And the response should match json:
    """
    {
      "error": "ValidationError",
      "description": "request body has an error: doesn't match schema #/components/schemas/ServiceInstanceProvisionRequestBody: there must be at most 0 properties"
    }
    """
