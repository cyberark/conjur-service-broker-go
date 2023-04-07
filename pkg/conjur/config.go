package conjur

// Config is the part of config that is specific to conjur
type Config struct {
	// CONJUR_VERSION: the version of Conjur enterprise, currently only version '5' is supported. Any other, non-empty value would raise an error.
	ConjurVersion uint32 `env:"CONJUR_VERSION" envDefault:"5"`
	// CONJUR_ACCOUNT: the account name for the Conjur instance you are connecting to.
	ConjurAccount string `env:"CONJUR_ACCOUNT,required"`
	// CONJUR_APPLIANCE_URL: the URL of the Conjur appliance instance you are connecting to. When using an HA Conjur master cluster, this should be the URL of the master load balancer.
	ConjurApplianceURL string `env:"CONJUR_APPLIANCE_URL,required"`
	// CONJUR_FOLLOWER_URL (HA only): If using high availability, this should be the URL of a load balancer for the cluster's Follower instances. This is the URL that applications use to communicate with Conjur.
	ConjurFollowerURL string `env:"CONJUR_FOLLOWER_URL"`
	// CONJUR_POLICY: the Policy branch where new Host identities should be added. The Conjur identity specified in CONJUR_AUTHN_LOGIN must have create and update permissions on this policy branch.
	// NOTE: The CONJUR_POLICY is optional, but is strongly recommended. If this value is not specified, the Service Broker uses the root Conjur policy.
	// NOTE: If you use multiple CloudFoundry foundations, this policy branch should include an identifier for the foundation to distinguish applications deployed in each foundation. For example, if you have both a production and development foundation, then your policy branches for each Conjur Service Broker might be cf/prod and cf/dev.
	ConjurPolicy string `env:"CONJUR_POLICY" envDefault:"root"`
	// CONJUR_AUTHN_LOGIN: the identity of a Conjur Host (of the form host/host-id) with create and update privileges on CONJUR_POLICY. This account is used to add and remove Hosts from Conjur policy as apps are deployed to or removed from the platform.
	//
	// If you are using Enterprise Conjur, you should add an annotation on the Service Broker Host in policy to indicate which platform the Service Broker is used on. The policy you load should similar to:
	//
	// - !host
	//  id: cf-service-broker
	//  annotations:
	//    platform: cloudfoundry
	// You may elect to set platform to cloudfoundry or to pivotalcloudfoundry, for example. This annotation is used to set annotations on Hosts added by the Service Broker, so that they show in the Conjur UI with the appropriate platform logo.
	//
	// NOTE: The CONJUR_AUTHN_LOGIN value for the Host created in policy above is host/cf-service-broker.
	ConjurAuthNLogin string `env:"CONJUR_AUTHN_LOGIN,required,unset"`
	// CONJUR_AUTHN_API_KEY: the API Key of the Conjur Host whose identity you have provided in CONJUR_AUTHN_LOGIN.
	ConjurAuthNAPIKey string `env:"CONJUR_AUTHN_API_KEY,required,unset"`
	// CONJUR_SSL_CERTIFICATE: the PEM-encoded x509 CA certificate chain for Conjur. This is required if your Conjur installation uses SSL (e.g. Conjur Enterprise).
	//
	// This value may be obtained by running the command:
	//
	// $ openssl s_client -showcerts -servername [CONJUR_DNS_NAME] \
	//    -connect [CONJUR_DNS_NAME]:443 < /dev/null 2> /dev/null \
	//    | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p'
	// -----BEGIN CERTIFICATE-----
	// ...
	// -----END CERTIFICATE-----
	ConjurSSLCertificate string `env:"CONJUR_SSL_CERTIFICATE"`
}
