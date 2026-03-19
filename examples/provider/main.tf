terraform {
  required_providers {
    hestiacp = {
      source  = "registry.terraform.io/your-org/hestiacp"
      version = "~> 0.1"
    }
  }
}

provider "hestiacp" {
  # url        = "https://myserver.com:8083"   # or set HESTIACP_URL
  # access_key = "ACCESSKEY:SECRETKEY"          # or set HESTIACP_ACCESS_KEY
}

# ── User ──────────────────────────────────────────────────────────────────
resource "hestiacp_user" "alice" {
  username   = "alice"
  password   = "S3cur3P@ss!"
  email      = "alice@example.com"
  first_name = "Alice"
  package    = "default"
}

# ── Web domain ────────────────────────────────────────────────────────────
resource "hestiacp_web_domain" "example" {
  user   = hestiacp_user.alice.username
  domain = "example.com"
}

# ── DNS zone ──────────────────────────────────────────────────────────────
resource "hestiacp_dns_zone" "example" {
  user   = hestiacp_user.alice.username
  domain = "example.com"
  ip     = "1.2.3.4"
}

# ── DNS records ───────────────────────────────────────────────────────────
resource "hestiacp_dns_record" "a_root" {
  user   = hestiacp_user.alice.username
  domain = hestiacp_dns_zone.example.domain
  record = "@"
  type   = "A"
  value  = "1.2.3.4"
}

resource "hestiacp_dns_record" "a_www" {
  user   = hestiacp_user.alice.username
  domain = hestiacp_dns_zone.example.domain
  record = "www"
  type   = "CNAME"
  value  = "example.com."
}

resource "hestiacp_dns_record" "mx" {
  user     = hestiacp_user.alice.username
  domain   = hestiacp_dns_zone.example.domain
  record   = "@"
  type     = "MX"
  value    = "mail.example.com."
  priority = 10
}

# ── Database ──────────────────────────────────────────────────────────────
resource "hestiacp_database" "app" {
  user        = hestiacp_user.alice.username
  db_name     = "appdb"
  db_user     = "appuser"
  db_password = "DbP@ss123"
  db_type     = "mysql"
}

# ── Email domain ──────────────────────────────────────────────────────────
resource "hestiacp_email_domain" "example" {
  user   = hestiacp_user.alice.username
  domain = "example.com"
}

# ── Email account ─────────────────────────────────────────────────────────
resource "hestiacp_email_account" "info" {
  user     = hestiacp_user.alice.username
  domain   = hestiacp_email_domain.example.domain
  account  = "info"
  password = "M@ilP@ss!"
  quota    = 1024 # MB; 0 = unlimited
}

# ── SSL (Let's Encrypt) ───────────────────────────────────────────────────
resource "hestiacp_ssl" "example" {
  user    = hestiacp_user.alice.username
  domain  = hestiacp_web_domain.example.domain
  aliases = "www.example.com"
}

# ── Backup ────────────────────────────────────────────────────────────────
resource "hestiacp_backup" "alice" {
  user = hestiacp_user.alice.username
}
