@e2e
Feature: End to End test of service broker with org space identity

#  @enable-space-host
  Scenario: Space Host Identity
  Given I create an org and space
  And I install the Conjur service broker

  When I create a service instance for Conjur
  Then the policy for the org and space exists
  And the space host exists
  And the space host api key variable exists

  When I load a secret into Conjur
  And I privilege the org layer to access a secret in Conjur
  And I privilege the space layer to access a secret in Conjur

  And I push the sample app to PCF
  And I start the app
  Then I can retrieve the secret values from the app
