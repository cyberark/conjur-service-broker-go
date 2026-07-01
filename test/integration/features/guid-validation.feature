Feature: GUID format validation on provision and bind endpoints

Scenario: Provision is rejected when context.organization_guid contains a path separator
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a90" with body:
    """
    {
      "context": {
        "organization_guid": "tenant/../victim-org",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {}
    }
    """
    Then the response code should be 400

  Scenario: Provision is rejected when context.organization_guid is a plain string
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a91" with body:
    """
    {
      "context": {
        "organization_guid": "not-a-uuid",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {}
    }
    """
    Then the response code should be 400

  Scenario: Provision is rejected when context.space_guid contains an embedded slash
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a92" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
        "space_guid": "victim-space/injected"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {}
    }
    """
    Then the response code should be 400

  Scenario: Provision is rejected when context.space_guid is a Conjur resource ID
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a93" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
        "space_guid": "myaccount:policy:root"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {}
    }
    """
    Then the response code should be 400

  Scenario: Provision is rejected when top-level organization_guid is not a UUID
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a94" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "../../escape",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {}
    }
    """
    Then the response code should be 400

  Scenario: Provision is rejected when top-level space_guid is not a UUID
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a95" with body:
    """
    {
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "victim-tenant/cf/spaces",
      "parameters": {}
    }
    """
    Then the response code should be 400

  Scenario: Bind is rejected when context.organization_guid contains a path separator
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a96/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f3" with body:
    """
    {
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id":    "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "context": {
        "organization_guid": "attacker-org/victim-org",
        "space_guid":        "8c56f85c-c16e-4158-be79-5dac74f970de"
      }
    }
    """
    Then the response code should be 400

  Scenario: Bind is rejected when context.space_guid contains a path separator
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a96/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f4" with body:
    """
    {
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id":    "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
        "space_guid":        "attacker/victim-space"
      }
    }
    """
    Then the response code should be 400

  Scenario: Bind is rejected when context.organization_guid is a Conjur resource ID
    When I send "PUT" request to "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32a96/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f5" with body:
    """
    {
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id":    "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "context": {
        "organization_guid": "dev:policy:cf",
        "space_guid":        "8c56f85c-c16e-4158-be79-5dac74f970de"
      }
    }
    """
    Then the response code should be 400
