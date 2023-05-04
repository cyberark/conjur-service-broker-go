package servicebroker

const (
	bindBody = `{
    "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
    "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "bind_resource": {
        "app_guid": "bb841d2b-8287-47a9-ac8f-eef4c16106f2"
      },
      "context": {
          "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
          "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "parameters": {
        "parameter1": 1,
        "parameter2": "foo"
      }
}`
	emptyBindResp = `{"credentials":{"account":"","appliance_url":"","authn_api_key":"","authn_login":"","ssl_certificate":"","version":0}}`
	provisionBody = `{
			"context": {
				"organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
				"space_guid":        "8c56f85c-c16e-4158-be79-5dac74f970de",
				"organization_name": "my-organization",
				"space_name":        "my-space"
			},
			"service_id":        "c024e536-6dc4-45c6-8a53-127e7f8275ab",
			"plan_id":           "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
			"organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
			"space_guid":        "8c56f85c-c16e-4158-be79-5dac74f970de",
			"parameters":        {}
		}`
)
