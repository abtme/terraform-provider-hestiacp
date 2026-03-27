# Plan: Publishing the HestiaCP Terraform Provider

## Phase 1 — Decisions & Prerequisites
> Nothing else unblocks until these are done.

| # | Action | Notes |
|---|---|---|
| 1.1 | **Choose GitHub org/username** | This becomes `registry.terraform.io/<org>/hestiacp` — permanent, users write it in their configs |
| 1.2 | **Generate a GPG signing key** | The Registry requires all release binaries to be GPG-signed. Key must exist before Registry signup |

### 1.2 — GPG Key Generation

```bash
# Generate a new key (RSA 4096 or ed25519)
gpg --full-generate-key
# Select: RSA and RSA, 4096 bits, no expiry, real name, GitHub email

# Export public key — paste this into the Terraform Registry UI
gpg --armor --export YOUR_KEY_ID

# Export private key — add this as a GitHub secret
gpg --armor --export-secret-keys YOUR_KEY_ID

# Note your fingerprint
gpg --fingerprint YOUR_KEY_ID
```

Store the private key and passphrase in a password manager.

---

## Phase 2 — Code Fixes

| # | Action | Why |
|---|---|---|
| 2.1 | Fix `go 1.25.0` → `go 1.22.0` in `go.mod` | Go 1.25 doesn't exist — GoReleaser will fail |
| 2.2 | Replace `your-org` in 10+ files | `go.mod`, `main.go`, `GNUmakefile`, all resource files, examples, README |
| 2.3 | Implement `ImportState` on all 9 resources | The user acceptance test already has `ImportState: true` — it will panic without this |
| 2.4 | Create `.gitignore` | Exclude the committed binary, `dist/`, `.terraform/` |
| 2.5 | Delete stray `{internal` directory | Junk artifact from a shell brace-expansion gone wrong |
| 2.6 | Create `.goreleaser.yml` | Needed for multi-platform release builds + GPG signing |
| 2.7 | Create GitHub Actions workflows | `release.yml` (fires on `v*` tag) + `test.yml` (CI on push/PR) |
| 2.8 | Run `tfplugindocs`, create per-resource example files | Registry requires a `docs/` directory |

### 2.3 — ImportState ID Conventions

| Resource | Import ID format | Example |
|---|---|---|
| `hestiacp_user` | `<username>` | `alice` |
| `hestiacp_web_domain` | `<user>/<domain>` | `alice/example.com` |
| `hestiacp_dns_zone` | `<user>/<domain>` | `alice/example.com` |
| `hestiacp_dns_record` | `<user>/<domain>/<record_id>` | `alice/example.com/42` |
| `hestiacp_database` | `<user>/<db_name>` | `alice/appdb` |
| `hestiacp_email_domain` | `<user>/<domain>` | `alice/example.com` |
| `hestiacp_email_account` | `<user>/<domain>/<account>` | `alice/example.com/info` |
| `hestiacp_ssl` | `<user>/<domain>` | `alice/example.com` |
| `hestiacp_backup` | not importable (it's a trigger, not a managed object) | — |

> **Note:** `hestiacp_email_account` Create currently sets ID to `account@domain`. Change it to `user/domain/account` for consistency with the import format. No breaking change since nothing is published yet.

### 2.6 — `.goreleaser.yml`

```yaml
version: 2
before:
  hooks:
    - go mod tidy
builds:
  - env:
      - CGO_ENABLED=0
    mod_timestamp: "{{ .CommitTimestamp }}"
    flags:
      - -trimpath
    ldflags:
      - "-s -w -X main.version={{.Version}}"
    goos:
      - freebsd
      - windows
      - linux
      - darwin
    goarch:
      - amd64
      - "386"
      - arm
      - arm64
    ignore:
      - goos: darwin
        goarch: "386"
    binary: "{{ .ProjectName }}_v{{ .Version }}"
archives:
  - format: zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
checksum:
  name_template: "{{ .ProjectName }}_{{ .Version }}_SHA256SUMS"
  algorithm: sha256
signs:
  - artifacts: checksum
    args:
      - "--batch"
      - "--local-user"
      - "{{ .Env.GPG_FINGERPRINT }}"
      - "--output"
      - "${signature}"
      - "--detach-sign"
      - "${artifact}"
release:
  draft: true
changelog:
  skip: true
```

### 2.7 — GitHub Actions Workflows

**`.github/workflows/release.yml`** — fires on `v*` tags:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - name: Import GPG key
        uses: crazy-max/ghaction-import-gpg@v6
        id: import_gpg
        with:
          gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}
          passphrase: ${{ secrets.GPG_PASSPHRASE }}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GPG_FINGERPRINT: ${{ steps.import_gpg.outputs.fingerprint }}
```

**`.github/workflows/test.yml`** — fires on push/PR to main:

```yaml
name: Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - run: go mod tidy
      - run: go test ./... -v -timeout 120s
      - run: go build ./...
```

---

## Phase 3 — GitHub Repository Setup

| # | Action |
|---|---|
| 3.1 | Create a **public** GitHub repo named exactly `terraform-provider-hestiacp` |
| 3.2 | Add secrets: `GPG_PRIVATE_KEY` and `GPG_PASSPHRASE` (Settings → Secrets → Actions) |
| 3.3 | Push code to `main`, verify the Test workflow passes |

---

## Phase 4 — Terraform Registry

| # | Action |
|---|---|
| 4.1 | Sign in to [registry.terraform.io](https://registry.terraform.io) with your GitHub account |
| 4.2 | Register your GPG **public** key: account → Manage → GPG Keys → Add |
| 4.3 | Publish the provider: Publish → Provider → select the repo |

> The Registry requires the repo to be public and named `terraform-provider-hestiacp`.

---

## Phase 5 — First Release

| # | Action |
|---|---|
| 5.1 | `git tag v0.1.0 && git push github v0.1.0` |
| 5.2 | Watch GoReleaser build all platforms in GitHub Actions |
| 5.3 | Un-draft the GitHub Release |
| 5.4 | Verify `registry.terraform.io/<org>/hestiacp` shows `v0.1.0` (may take a few minutes; click Resync if needed) |

---

## Phase 6 — Verify & Test

| # | Action |
|---|---|
| 6.1 | `terraform init` from a fresh config using the Registry source — confirms end-to-end install |
| 6.2 | Run acceptance tests with `TF_ACC=1` against a real HestiaCP server |

### 6.1 — Fresh install test

```hcl
# test-install/main.tf
terraform {
  required_providers {
    hestiacp = {
      source  = "registry.terraform.io/<GITHUB_ORG>/hestiacp"
      version = "~> 0.1"
    }
  }
}

provider "hestiacp" {}
```

```bash
export HESTIACP_URL="https://your-server:8083"
export HESTIACP_ACCESS_KEY="YOURKEY:YOURSECRET"
terraform -chdir=test-install init
```

### 6.2 — Acceptance tests

```bash
export TF_ACC=1
export HESTIACP_URL="https://your-server:8083"
export HESTIACP_ACCESS_KEY="YOURKEY:YOURSECRET"
make testacc
```

> The `hestiacp_ssl` test requires your test domain to have valid DNS pointing to the server (Let's Encrypt requirement).

---

## Open Questions

1. What is your GitHub username/org? (replaces `your-org` everywhere)
2. Do you already have a GPG key to use, or do you need to generate one?
3. Keep GitLab CI alongside GitHub Actions, or move fully to GitHub?
