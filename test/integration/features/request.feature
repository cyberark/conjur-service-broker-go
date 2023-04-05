Feature: Service Broker request configuration

Scenario: A request is sent without the required X-Broker-API-Version header
  Given my request doesn't include the X-Broker-API-Version header
  When I send "GET" request to "/v2/catalog"
  Then the response code should be 412
  And the response should match json:
      """
      {
        "error": "ValidationError",
        "description": "X-Broker-API-Version value is required but missing"
      }
      """
