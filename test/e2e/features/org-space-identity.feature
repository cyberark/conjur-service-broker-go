@e2e
Feature: End to End test of service broker with org space identity

  Scenario: Space Host Identity
  Given I create an org and space
  And I create space developer user and login
  And I install the Conjur service broker with space host identity

  When I create a service instance for Conjur
  Then the policies for the org and space exists
  And the space host exists
  And the space host api key variable exists

  When I load the secrets into Conjur
  And I privilege the org layer to access a secret in Conjur
  And I privilege the space layer to access a secret in Conjur

  And I push the sample app to PCF
  And I start the app
  Then I can retrieve the secret values from the app
