package servicebroker

import (
	"github.com/gin-gonic/gin"
)

var catalogResp = Catalog{
	Services: &[]Service{
		{
			Name:        "cyberark-conjur",
			Id:          "c024e536-6dc4-45c6-8a53-127e7f8275ab",
			Description: "An open source security service that provides secrets management, machine-identity based authorization, and more.",
			Bindable:    true,
			Metadata: &Metadata{
				"displayName":         "CyberArk Conjur",
				"imageUrl":            "https://www.conjur.org/img/feature-icons/machine-identity-teal.svg",
				"providerDisplayName": "CyberArk",
				"documentationUrl":    "https://www.conjur.org/api.html",
				"supportUrl":          "https://www.conjur.org/support.html",
				"shareable":           false,
			},
			Plans: []Plan{
				{
					Name:        "community",
					Id:          "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
					Description: "Community service plan",
					Metadata: &Metadata{
						"display_name": "Conjur",
						"bullets": []string{
							"Machine Identity", "Secrets management", "Role-based access control",
						},
					},
				},
			},
		},
	},
}

// CatalogGet get the catalog of services that the service broker offers
// (GET /v2/catalog)
func (*serverImpl) CatalogGet(c *gin.Context, params CatalogGetParams) {
	c.JSON(200, catalogResp)
}
