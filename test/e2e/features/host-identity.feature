@e2e
Feature: End to End test of service broker with host identity

  Scenario: Service broker functions correctly with PCF
    Given I create an org and space
    And I create space developer user and login
    And I install the Conjur service broker

    When I create a service instance for Conjur
    Then the policies for the org and space exists

    When I load the secrets into Conjur
    And I privilege the org group to access a secret in Conjur
    And I privilege the space group to access a secret in Conjur

    And I push the sample app to PCF
    And I privilege the app host to access a secret in Conjur
    And I start the app
    Then the binding host exists

    And I can retrieve the secret values from the app

    When I remove the service instance
    Then the policies for the org and space exists

    # Redeploy and run app, maintaining existing policy
    When I create a service instance for Conjur
    Then the policies for the org and space exists

    When I push the sample app to PCF
    # The app host will have a new binding ID, so we need to grant
    # permissions for the app specific secret again, but not for the
    # org/space secrets:
    And I privilege the app host to access a secret in Conjur
    And I start the app
    Then I can retrieve the secret values from the app

    When I remove the service instance
    Then the policies for the org and space exists
