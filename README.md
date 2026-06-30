# terraform-provider-hestiacp

A Terraform provider for [HestiaCP](https://hestiacp.com/) built with the
[Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

## Requirements

| Tool      | Version  |
|-----------|----------|
| Go        | ≥ 1.21   |
| Terraform | ≥ 1.3.0  |

## Resources

| Resource                      | Description                              |
|-------------------------------|------------------------------------------|
| `hestiacp_user`               | User account                             |
| `hestiacp_web_domain`         | Web domain / virtual host                |
| `hestiacp_dns_zone`           | DNS zone                                 |
| `hestiacp_dns_record`         | Individual DNS record (A, CNAME, MX …)  |
| `hestiacp_database`           | MySQL / PostgreSQL database + user       |
| `hestiacp_email_domain`       | Mail domain                              |
| `hestiacp_email_account`      | Mailbox                                  |
| `hestiacp_ssl`                | Let's Encrypt SSL certificate            |
| `hestiacp_backup`             | User backup snapshot                     |

## Authentication

HestiaCP v1.6+ uses `ACCESS_KEY:SECRET_KEY` pairs.

Create a key on your server:

```bash
v-add-access-key admin '*' terraform json
```

Then configure the provider:

```hcl
provider "hestiacp" {
  url        = "https://myserver.com:8083"
  access_key = "ACCESSKEY:SECRETKEY"
}
```

Or via environment variables:

```bash
export HESTIACP_URL="https://myserver.com:8083"
export HESTIACP_ACCESS_KEY="ACCESSKEY:SECRETKEY"
```

> **Important:** Before the provider can reach the API from a remote machine,
> whitelist the runner's IP under *Server Settings → API* in the HestiaCP UI,
> or set the allowed list to `allow-all` for testing.

## Quick Start

```hcl
terraform {
  required_providers {
    hestiacp = {
      source  = "registry.terraform.io/abtme/hestiacp"
      version = "~> 0.1"
    }
  }
}

provider "hestiacp" {}   # reads HESTIACP_URL + HESTIACP_ACCESS_KEY

resource "hestiacp_user" "alice" {
  username   = "alice"
  password   = "S3cur3P@ss!"
  email      = "alice@example.com"
  package    = "default"
}

resource "hestiacp_web_domain" "site" {
  user   = hestiacp_user.alice.username
  domain = "example.com"
}

resource "hestiacp_ssl" "site" {
  user   = hestiacp_user.alice.username
  domain = hestiacp_web_domain.site.domain
}
```

## Local Development

```bash
# Build and install into ~/.terraform.d/plugins/…
make install

# Point a test config at the local binary
cat > ~/.terraformrc <<'EOF'
provider_installation {
  dev_overrides {
    "registry.terraform.io/your-org/hestiacp" = "/home/<you>/.terraform.d/plugins/…"
  }
  direct {}
}
EOF

terraform -chdir=examples/provider init
terraform -chdir=examples/provider apply
```

## Running Tests

Unit tests (no live server needed):

```bash
make test
```

Acceptance tests (requires a real HestiaCP instance):

```bash
export HESTIACP_URL="https://myserver.com:8083"
export HESTIACP_ACCESS_KEY="ACCESSKEY:SECRETKEY"
make testacc
```

## HestiaCP API Return Codes

| Code | Meaning            |
|------|--------------------|
| 0    | OK                 |
| 1    | E_ARGS             |
| 2    | E_INVALID          |
| 3    | E_NOTEXIST         |
| 4    | E_EXISTS           |
| 8    | E_LIMIT            |
| 10   | E_FORBIDDEN        |
| 11   | E_DISABLED         |

Full list in `internal/client/client.go`.

## Publishing to the Terraform Registry

1. Tag a release: `git tag v0.1.0 && git push origin v0.1.0`
2. GoReleaser (`.goreleaser.yml`) will build cross-platform binaries automatically via GitHub Actions.
3. Follow the [Terraform Registry publishing guide](https://developer.hashicorp.com/terraform/registry/providers/publishing).
