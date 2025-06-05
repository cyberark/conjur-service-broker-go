
package api

import (
	"io"
	"net/http"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

// Client is an interface generated for "github.com/cyberark/conjur-api-go/conjurapi.Client".
type Client interface {
	AddSecret(string, string) error
	AddSecretRequest(string, string) (*http.Request, error)
	Authenticate(authn.LoginPair) ([]byte, error)
	AuthenticateReader(authn.LoginPair) (io.ReadCloser, error)
	AuthenticateRequest(authn.LoginPair) (*http.Request, error)
	AuthenticatorStatus(string, string) (*conjurapi.AuthenticatorStatusResponse, error)
	AuthenticatorStatusRequest(string, string) (*http.Request, error)
	ChangeCurrentUserPassword(string) ([]byte, error)
	ChangeUserPassword(string, string, string) ([]byte, error)
	ChangeUserPasswordRequest(string, string, string) (*http.Request, error)
	CheckPermission(string, string) (bool, error)
	CheckPermissionForRole(string, string, string) (bool, error)
	CheckPermissionForRoleRequest(string, string, string) (*http.Request, error)
	CheckPermissionRequest(string, string) (*http.Request, error)
	CreateHost(string, string) (conjurapi.HostFactoryHostResponse, error)
	CreateHostRequest(string, string) (*http.Request, error)
	CreateHostWithAnnotations(string, string, map[string]string) (conjurapi.HostFactoryHostResponse, error)
	CreateIssuer(conjurapi.Issuer) (conjurapi.Issuer, error)
	CreateToken(string, string, []string, int) ([]conjurapi.HostFactoryTokenResponse, error)
	CreateTokenRequest(string) (*http.Request, error)
	DeleteIssuer(string, bool) error
	DeleteToken(string) error
	DeleteTokenRequest(string) (*http.Request, error)
	DryRunPolicy(conjurapi.PolicyMode, string, io.Reader) (*conjurapi.DryRunPolicyResponse, error)
	EnableAuthenticator(string, string, bool) error
	EnableAuthenticatorRequest(string, string, bool) (*http.Request, error)
	EnterpriseServerInfo() (*conjurapi.EnterpriseInfoResponse, error)
	FetchPolicy(string, bool, uint, uint) ([]byte, error)
	ForceRefreshToken() error
	GetAuthenticator() conjurapi.Authenticator
	GetConfig() conjurapi.Config
	GetHttpClient() *http.Client
	GetTelemetryHeader() string
	InternalAuthenticate() ([]byte, error)
	Issuer(string) (conjurapi.Issuer, error)
	Issuers() ([]conjurapi.Issuer, error)
	JWTAuthenticate(string, string) ([]byte, error)
	JWTAuthenticateRequest(string, string) (*http.Request, error)
	ListOidcProviders() ([]conjurapi.OidcProvider, error)
	ListOidcProvidersRequest() (*http.Request, error)
	LoadPolicy(conjurapi.PolicyMode, string, io.Reader) (*conjurapi.PolicyResponse, error)
	LoadPolicyRequest(conjurapi.PolicyMode, string, io.Reader, bool) (*http.Request, error)
	Login(string, string) ([]byte, error)
	LoginRequest(string, string) (*http.Request, error)
	NeedsTokenRefresh() bool
	OidcAuthenticate(string, string, string) ([]byte, error)
	OidcAuthenticateRequest(string, string, string) (*http.Request, error)
	OidcTokenAuthenticate(string) ([]byte, error)
	OidcTokenAuthenticateRequest(string) (*http.Request, error)
	PermittedRoles(string, string) ([]string, error)
	PermittedRolesRequest(string, string) (*http.Request, error)
	PublicKeys(string, string) ([]byte, error)
	PublicKeysRequest(string, string) (*http.Request, error)
	PurgeCredentials() error
	RefreshToken() error
	Resource(string) (map[string]interface{}, error)
	ResourceExists(string) (bool, error)
	ResourceIDs(*conjurapi.ResourceFilter) ([]string, error)
	ResourceRequest(string) (*http.Request, error)
	Resources(*conjurapi.ResourceFilter) ([]map[string]interface{}, error)
	ResourcesCount(*conjurapi.ResourceFilter) (*conjurapi.ResourcesCount, error)
	ResourcesCountRequest(*conjurapi.ResourceFilter) (*http.Request, error)
	ResourcesRequest(*conjurapi.ResourceFilter) (*http.Request, error)
	RetrieveBatchSecrets([]string) (map[string][]byte, error)
	RetrieveBatchSecretsRequest([]string, bool) (*http.Request, error)
	RetrieveBatchSecretsSafe([]string) (map[string][]byte, error)
	RetrieveSecret(string) ([]byte, error)
	RetrieveSecretReader(string) (io.ReadCloser, error)
	RetrieveSecretRequest(string) (*http.Request, error)
	RetrieveSecretWithVersion(string, int) ([]byte, error)
	RetrieveSecretWithVersionReader(string, int) (io.ReadCloser, error)
	RetrieveSecretWithVersionRequest(string, int) (*http.Request, error)
	Role(string) (map[string]interface{}, error)
	RoleExists(string) (bool, error)
	RoleMembers(string) ([]map[string]interface{}, error)
	RoleMembersRequest(string) (*http.Request, error)
	RoleMemberships(string) ([]map[string]interface{}, error)
	RoleMembershipsAll(string) ([]string, error)
	RoleMembershipsRequest(string) (*http.Request, error)
	RoleMembershipsRequestWithOptions(string, bool) (*http.Request, error)
	RoleRequest(string) (*http.Request, error)
	RootRequest() (*http.Request, error)
	RotateAPIKey(string) ([]byte, error)
	RotateAPIKeyReader(string) (io.ReadCloser, error)
	RotateAPIKeyRequest(string) (*http.Request, error)
	RotateCurrentRoleAPIKey() ([]byte, error)
	RotateCurrentRoleAPIKeyRequest(string, string) (*http.Request, error)
	RotateCurrentUserAPIKey() ([]byte, error)
	RotateCurrentUserAPIKeyRequest(string, string) (*http.Request, error)
	RotateHostAPIKey(string) ([]byte, error)
	RotateUserAPIKey(string) ([]byte, error)
	ServerInfoRequest() (*http.Request, error)
	ServerVersion() (string, error)
	ServerVersionFromRoot() (string, error)
	SetAuthenticator(conjurapi.Authenticator)
	SetHttpClient(*http.Client)
	SubmitRequest(*http.Request) (*http.Response, error)
	UpdateIssuer(string, conjurapi.IssuerUpdate) (conjurapi.Issuer, error)
	VerifyMinServerVersion(string) error
	WhoAmI() ([]byte, error)
	WhoAmIRequest() (*http.Request, error)
}
