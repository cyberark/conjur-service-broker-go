package servicebroker

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const expectedCatalogResp = `{
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
                        "free": true,
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
      }`

func Test_server_CatalogGet(t *testing.T) {
	s := &server{}
	w, c := ginTestCtx(t, http.MethodGet, "/v2/catalog", "", false)
	s.CatalogGet(c, CatalogGetParams{})
	require.Empty(t, c.Errors.Errors())
	require.JSONEq(t, expectedCatalogResp, w.Body.String())
	require.Equal(t, http.StatusOK, w.Code)
}
