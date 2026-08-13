package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/abtme/terraform-provider-hestiacp/internal/provider"
)

// testAccProtoV6ProviderFactories is used by acceptance tests.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hestiacp": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	for _, env := range []string{"HESTIACP_URL", "HESTIACP_ACCESS_KEY"} {
		if v := os.Getenv(env); v == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}
}

// ── User ──────────────────────────────────────────────────────────────────

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccUserConfig("tfacctest", "acc@test.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hestiacp_user.test", "username", "tfacctest"),
					resource.TestCheckResourceAttr("hestiacp_user.test", "email", "acc@test.com"),
				),
			},
			// Import
			{
				ResourceName:            "hestiacp_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update email
			{
				Config: testAccUserConfig("tfacctest", "updated@test.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hestiacp_user.test", "email", "updated@test.com"),
				),
			},
		},
	})
}

func testAccUserConfig(username, email string) string {
	return `
provider "hestiacp" {}

resource "hestiacp_user" "test" {
  username   = "` + username + `"
  password   = "AccT3stP@ss!"
  email      = "` + email + `"
  package    = "default"
}
`
}

// ── Web Domain ────────────────────────────────────────────────────────────

func TestAccWebDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebDomainConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hestiacp_web_domain.test", "domain", "tfacc-test.example.com"),
				),
			},
		},
	})
}

func testAccWebDomainConfig() string {
	return `
provider "hestiacp" {}

resource "hestiacp_user" "domain_owner" {
  username = "tfaccdomainowner"
  password = "AccT3stP@ss!"
  email    = "domainowner@test.com"
}

resource "hestiacp_web_domain" "test" {
  user   = hestiacp_user.domain_owner.username
  domain = "tfacc-test.example.com"
}
`
}

// ── Database ──────────────────────────────────────────────────────────────

func TestAccDatabaseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hestiacp_database.test", "db_name", "testdb"),
					resource.TestCheckResourceAttr("hestiacp_database.test", "db_type", "mysql"),
				),
			},
		},
	})
}

func testAccDatabaseConfig() string {
	return `
provider "hestiacp" {}

resource "hestiacp_user" "db_owner" {
  username = "tfaccdbowner"
  password = "AccT3stP@ss!"
  email    = "dbowner@test.com"
}

resource "hestiacp_database" "test" {
  user        = hestiacp_user.db_owner.username
  db_name     = "testdb"
  db_user     = "testdbuser"
  db_password = "DbAccP@ss!"
  db_type     = "mysql"
}
`
}
