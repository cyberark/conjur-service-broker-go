# file: catalog.feature
Feature: get catalog
  In order to know service broker capabilities
  As an CF compliant platform
  I need to be able to request catalog

  Scenario: should get catalog
    When I send "GET" request to "/v2/catalog"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "services": [
            {
                "bindable": true,
                "description": "An open source security service that provides secrets management, machine-identity based authorization, and more.",
                "id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
                "metadata": {
                    "displayName": "CyberArk Conjur",
                    "documentationUrl": "https://www.conjur.org/api.html",
                    "imageUrl": "https://www.conjur.org/img/feature-icons/machine-identity-teal.svg",
                    "providerDisplayName": "CyberArk",
                    "shareable": false,
                    "supportUrl": "https://www.conjur.org/support.html"
                },
                "name": "cyberark-conjur",
                "plans": [
                    {
                        "description": "Community service plan",
                        "id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
                        "metadata": {
                            "bullets": [
                                "Machine Identity",
                                "Secrets management",
                                "Role-based access control"
                            ],
                            "display_name": "Conjur"
                        },
                        "name": "community"
                    }
                ]
            }
        ]
      }
      """