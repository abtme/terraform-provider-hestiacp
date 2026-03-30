# HestiaCP Terraform Provider - API Reference

This document lists the HestiaCP CLI/API commands available in the HestiaCP control panel, organized by category. The **Terraform Resource** column indicates which resource in this provider implements each command, or `-` if the command is not yet implemented.

---

## Table of Contents

- [User Management](#user-management)
- [Web Domains](#web-domains)
- [DNS](#dns)
- [Databases](#databases)
- [Mail](#mail)
- [SSL / Let's Encrypt](#ssl--lets-encrypt)
- [Backups](#backups)
- [Cron Jobs](#cron-jobs)
- [Firewall](#firewall)
- [File System](#file-system)
- [Access Keys](#access-keys)
- [Remote DNS](#remote-dns)
- [System](#system)
- [Restore / Rebuild / Restart](#restore--rebuild--restart)
- [Search / Logging / Misc](#search--logging--misc)

---

## User Management

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-user` | Add a new user account | `hestiacp_user` (create) |
| `v-add-user-2fa` | Enable two-factor authentication for a user | `-` |
| `v-add-user-composer` | Install Composer for a user | `-` |
| `v-add-user-notification` | Add a notification for a user | `-` |
| `v-add-user-package` | Add a user hosting package | `-` |
| `v-add-user-sftp-jail` | Enable SFTP jail for a user | `-` |
| `v-add-user-sftp-key` | Add an SFTP key for a user | `-` |
| `v-add-user-ssh-key` | Add an SSH key for a user | `hestiacp_ssh_key` (create) — key comment is used as the identifier |
| `v-add-user-wp-cli` | Install WP-CLI for a user | `-` |
| `v-acknowledge-user-notification` | Mark a user notification as acknowledged | `-` |
| `v-change-user-config-value` | Change a specific user configuration value | `-` |
| `v-change-user-contact` | Change a user's contact email address | `hestiacp_user` (update) |
| `v-change-user-language` | Change the UI language for a user | `-` |
| `v-change-user-name` | Change the display name of a user | `-` |
| `v-change-user-ns` | Change the default nameservers for a user | `-` |
| `v-change-user-package` | Change the hosting package assigned to a user | `-` |
| `v-change-user-password` | Change a user's password | `hestiacp_user` (update) |
| `v-change-user-php-cli` | Change the default PHP CLI version for a user | `-` |
| `v-change-user-rkey` | Regenerate the API key for a user | `-` |
| `v-change-user-role` | Change the role (admin/user) of a user | `-` |
| `v-change-user-shell` | Change the default shell for a user | `-` |
| `v-change-user-sort-order` | Change the sort order of resources for a user | `-` |
| `v-change-user-template` | Change the default web template for a user | `-` |
| `v-change-user-theme` | Change the UI theme for a user | `-` |
| `v-check-user-2fa` | Verify a user's two-factor authentication token | `-` |
| `v-check-user-hash` | Check a user's authentication hash | `-` |
| `v-check-user-password` | Verify a user's password | `-` |
| `v-copy-user-package` | Copy an existing user package to a new one | `-` |
| `v-delete-user` | Delete a user account | `hestiacp_user` (delete) |
| `v-delete-user-2fa` | Remove two-factor authentication from a user | `-` |
| `v-delete-user-auth-log` | Delete a user's authentication log | `-` |
| `v-delete-user-backup-exclusions` | Remove backup exclusions for a user | `-` |
| `v-delete-user-ips` | Remove IP addresses assigned to a user | `-` |
| `v-delete-user-log` | Delete a user's activity log | `-` |
| `v-delete-user-notification` | Delete a user notification | `-` |
| `v-delete-user-package` | Delete a user hosting package | `-` |
| `v-delete-user-sftp-jail` | Remove SFTP jail for a user | `-` |
| `v-delete-user-ssh-key` | Remove an SSH key from a user | `hestiacp_ssh_key` (delete) — arg2 is the key comment |
| `v-delete-user-stats` | Delete statistics for a user | `-` |
| `v-get-user-salt` | Retrieve the password salt for a user | `-` |
| `v-get-user-value` | Get a specific configuration value for a user | `-` |
| `v-list-user` | List details for a specific user | `hestiacp_user` (read) |
| `v-list-user-auth-log` | List authentication log entries for a user | `-` |
| `v-list-user-backup` | List details of a specific user backup | `-` |
| `v-list-user-backup-exclusions` | List backup exclusions for a user | `-` |
| `v-list-user-ips` | List IP addresses assigned to a user | `-` |
| `v-list-user-log` | List activity log entries for a user | `-` |
| `v-list-user-notifications` | List all notifications for a user | `-` |
| `v-list-user-ns` | List default nameservers for a user | `-` |
| `v-list-user-package` | List details of a specific user package | `-` |
| `v-list-user-packages` | List all available user packages | `-` |
| `v-list-user-ssh-key` | List SSH keys for a user | `hestiacp_ssh_key` (read) — KEY field is SHA256 fingerprint, not public key |
| `v-list-user-stats` | List resource usage statistics for a user | `-` |
| `v-list-users` | List all users on the server | `-` |
| `v-list-users-stats` | List resource usage statistics for all users | `-` |
| `v-rename-user-package` | Rename an existing user package | `-` |
| `v-suspend-user` | Suspend a user account | `-` |
| `v-unsuspend-user` | Unsuspend a user account | `-` |
| `v-update-user-backup-exclusions` | Update backup exclusions for a user | `-` |
| `v-update-user-cgroup` | Update cgroup resource limits for a user | `-` |
| `v-update-user-counters` | Recalculate resource counters for a user | `-` |
| `v-update-user-disk` | Update disk usage statistics for a user | `-` |
| `v-update-user-package` | Update an existing user hosting package | `-` |
| `v-update-user-quota` | Update disk quota for a user | `-` |
| `v-update-user-stats` | Update resource usage statistics for a user | `-` |

---

## Web Domains

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-domain` | Add a combined web, DNS, and mail domain | `-` |
| `v-add-fastcgi-cache` | Enable FastCGI cache for a web domain | `-` |
| `v-add-web-domain` | Add a new web domain | `hestiacp_web_domain` (create) |
| `v-add-web-domain-alias` | Add an alias (parked domain) to a web domain | `hestiacp_web_domain_alias` (create) |
| `v-add-web-domain-allow-users` | Allow system users to access a web domain | `-` |
| `v-add-web-domain-backend` | Add a backend (PHP-FPM) configuration to a web domain | `-` |
| `v-add-web-domain-ftp` | Add an FTP account to a web domain | `hestiacp_web_domain_ftp` (create) |
| `v-add-web-domain-httpauth` | Add HTTP basic authentication to a web domain | `hestiacp_web_domain_httpauth` (create) |
| `v-add-web-domain-proxy` | Enable proxy (nginx) for a web domain | `-` |
| `v-add-web-domain-redirect` | Add an HTTP redirect rule to a web domain | `hestiacp_web_domain_redirect` (create) |
| `v-add-web-domain-ssl` | Add an SSL certificate to a web domain | `-` |
| `v-add-web-domain-ssl-force` | Force HTTPS redirect for a web domain | `-` |
| `v-add-web-domain-ssl-hsts` | Enable HSTS for a web domain | `-` |
| `v-add-web-domain-ssl-preset` | Apply an SSL preset configuration to a web domain | `-` |
| `v-add-web-domain-stats` | Enable web statistics for a web domain | `-` |
| `v-add-web-domain-stats-user` | Add statistics user credentials for a web domain | `-` |
| `v-add-web-php` | Add a PHP version to the system | `-` |
| `v-change-domain-owner` | Change the owner of a domain | `-` |
| `v-change-web-domain-backend-tpl` | Change the backend template for a web domain | `-` |
| `v-change-web-domain-dirlist` | Toggle directory listing for a web domain | `-` |
| `v-change-web-domain-docroot` | Change the document root for a web domain | `-` |
| `v-change-web-domain-ftp-password` | Change the FTP account password for a web domain | `hestiacp_web_domain_ftp` (update) |
| `v-change-web-domain-ftp-path` | Change the FTP account home path for a web domain | `-` |
| `v-change-web-domain-httpauth` | Change HTTP authentication credentials for a web domain | `hestiacp_web_domain_httpauth` (update) |
| `v-change-web-domain-ip` | Change the IP address assigned to a web domain | `-` |
| `v-change-web-domain-name` | Rename a web domain | `-` |
| `v-change-web-domain-proxy-tpl` | Change the proxy template for a web domain | `-` |
| `v-change-web-domain-sslcert` | Replace the SSL certificate for a web domain | `-` |
| `v-change-web-domain-sslhome` | Change the SSL home directory for a web domain | `-` |
| `v-change-web-domain-stats` | Change the statistics engine for a web domain | `-` |
| `v-change-web-domain-tpl` | Change the web server template for a web domain | `-` |
| `v-delete-domain` | Delete a combined web, DNS, and mail domain | `-` |
| `v-delete-fastcgi-cache` | Disable FastCGI cache for a web domain | `-` |
| `v-delete-web-domain` | Delete a web domain | `hestiacp_web_domain` (delete) |
| `v-delete-web-domain-alias` | Remove an alias from a web domain | `hestiacp_web_domain_alias` (delete) |
| `v-delete-web-domain-allow-users` | Remove user access permission from a web domain | `-` |
| `v-delete-web-domain-backend` | Remove the backend configuration from a web domain | `-` |
| `v-delete-web-domain-ftp` | Remove an FTP account from a web domain | `hestiacp_web_domain_ftp` (delete) |
| `v-delete-web-domain-httpauth` | Remove HTTP basic authentication from a web domain | `hestiacp_web_domain_httpauth` (delete) |
| `v-delete-web-domain-proxy` | Disable proxy for a web domain | `-` |
| `v-delete-web-domain-redirect` | Remove a redirect rule from a web domain | `hestiacp_web_domain_redirect` (delete) |
| `v-delete-web-domains` | Delete all web domains for a user | `-` |
| `v-delete-web-domain-ssl` | Remove the SSL certificate from a web domain | `-` |
| `v-delete-web-domain-ssl-force` | Remove forced HTTPS redirect from a web domain | `-` |
| `v-delete-web-domain-ssl-hsts` | Disable HSTS for a web domain | `-` |
| `v-delete-web-domain-stats` | Disable web statistics for a web domain | `-` |
| `v-delete-web-domain-stats-user` | Remove statistics user credentials from a web domain | `-` |
| `v-delete-web-php` | Remove a PHP version from the system | `-` |
| `v-dump-site` | Export a full site archive | `-` |
| `v-list-default-php` | List the default PHP version configured on the system | `-` |
| `v-list-web-domain` | List details for a specific web domain | `hestiacp_web_domain` (read) |
| `v-list-web-domain-accesslog` | List the access log for a web domain | `-` |
| `v-list-web-domain-errorlog` | List the error log for a web domain | `-` |
| `v-list-web-domains` | List all web domains for a user | `-` |
| `v-list-web-domain-ssl` | List SSL certificate details for a web domain | `-` |
| `v-list-web-stats` | List web traffic statistics | `-` |
| `v-list-web-templates` | List available web server templates | `-` |
| `v-list-web-templates-backend` | List available backend (PHP-FPM) templates | `-` |
| `v-list-web-templates-proxy` | List available proxy (nginx) templates | `-` |
| `v-purge-nginx-cache` | Purge the nginx cache for a web domain | `-` |
| `v-quick-install-app` | Quickly install a web application on a domain | `-` |
| `v-suspend-domain` | Suspend a combined web, DNS, and mail domain | `-` |
| `v-suspend-web-domain` | Suspend a web domain | `-` |
| `v-suspend-web-domains` | Suspend all web domains for a user | `-` |
| `v-unsuspend-domain` | Unsuspend a combined web, DNS, and mail domain | `-` |
| `v-unsuspend-web-domain` | Unsuspend a web domain | `-` |
| `v-unsuspend-web-domains` | Unsuspend all web domains for a user | `-` |
| `v-update-web-domain-disk` | Update disk usage for a web domain | `-` |
| `v-update-web-domain-ssl` | Update the SSL certificate for a web domain | `-` |
| `v-update-web-domain-stat` | Update traffic statistics for a specific web domain | `-` |
| `v-update-web-domain-traff` | Update bandwidth usage for a web domain | `-` |
| `v-update-web-domains-disk` | Update disk usage for all web domains | `-` |
| `v-update-web-domains-stat` | Update traffic statistics for all web domains | `-` |
| `v-update-web-domains-traff` | Update bandwidth usage for all web domains | `-` |
| `v-update-web-templates` | Update web server configuration templates | `-` |

---

## DNS

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-dns-domain` | Add a new DNS zone | `hestiacp_dns_zone` (create) |
| `v-add-dns-on-web-alias` | Create DNS records for a web domain alias | `-` |
| `v-add-dns-record` | Add a DNS record to a zone | `hestiacp_dns_record` (create) |
| `v-change-dns-domain-dnssec` | Enable or disable DNSSEC for a DNS zone | `-` |
| `v-change-dns-domain-exp` | Change the expiry date of a DNS zone | `-` |
| `v-change-dns-domain-ip` | Change the IP address for a DNS zone | `-` |
| `v-change-dns-domain-soa` | Change the SOA record for a DNS zone | `-` |
| `v-change-dns-domain-tpl` | Change the template for a DNS zone | `-` |
| `v-change-dns-domain-ttl` | Change the default TTL for a DNS zone | `-` |
| `v-change-dns-record` | Modify an existing DNS record | `-` |
| `v-change-dns-record-id` | Change the ID of a DNS record | `-` |
| `v-delete-dns-domain` | Delete a DNS zone | `hestiacp_dns_zone` (delete) |
| `v-delete-dns-domains` | Delete all DNS zones for a user | `-` |
| `v-delete-dns-domains-src` | Delete DNS zones matching a source template | `-` |
| `v-delete-dns-on-web-alias` | Remove DNS records created for a web domain alias | `-` |
| `v-delete-dns-record` | Delete a DNS record from a zone | `hestiacp_dns_record` (delete) |
| `v-get-dns-domain-value` | Get a specific value from a DNS zone's configuration | `-` |
| `v-insert-dns-domain` | Insert a DNS zone from a pre-built configuration | `-` |
| `v-insert-dns-record` | Insert a single DNS record directly | `-` |
| `v-insert-dns-records` | Insert multiple DNS records directly | `-` |
| `v-list-dns-domain` | List details for a specific DNS zone | `hestiacp_dns_zone` (read) |
| `v-list-dns-domains` | List all DNS zones for a user | `-` |
| `v-list-dns-records` | List all DNS records for a zone | `hestiacp_dns_record` (read) |
| `v-list-dnssec-public-key` | List the DNSSEC public key for a DNS zone | `-` |
| `v-list-dns-template` | List details of a specific DNS template | `-` |
| `v-list-dns-templates` | List all available DNS templates | `-` |
| `v-suspend-dns-domain` | Suspend a DNS zone | `-` |
| `v-suspend-dns-domains` | Suspend all DNS zones for a user | `-` |
| `v-suspend-dns-record` | Suspend a specific DNS record | `-` |
| `v-sync-dns-cluster` | Synchronize DNS zones to the DNS cluster | `-` |
| `v-unsuspend-dns-domain` | Unsuspend a DNS zone | `-` |
| `v-unsuspend-dns-domains` | Unsuspend all DNS zones for a user | `-` |
| `v-unsuspend-dns-record` | Unsuspend a specific DNS record | `-` |
| `v-update-dns-templates` | Update DNS configuration templates | `-` |

---

## Databases

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-database` | Add a new database | `hestiacp_database` (create) |
| `v-add-database-host` | Add a database server host | `-` |
| `v-add-database-temp-user` | Add a temporary database user | `-` |
| `v-change-database-host-password` | Change the password for a database host connection | `-` |
| `v-change-database-owner` | Change the owner of a database | `-` |
| `v-change-database-password` | Change the password for a database user | `-` |
| `v-change-database-user` | Change the username for a database | `-` |
| `v-delete-database` | Delete a database | `hestiacp_database` (delete) |
| `v-delete-database-host` | Remove a database server host | `-` |
| `v-delete-databases` | Delete all databases for a user | `-` |
| `v-delete-database-temp-user` | Remove a temporary database user | `-` |
| `v-dump-database` | Export a database to a dump file | `-` |
| `v-import-database` | Import a database from a dump file | `-` |
| `v-list-database` | List details for a specific database | `hestiacp_database` (read) |
| `v-list-database-host` | List details for a specific database host | `-` |
| `v-list-database-hosts` | List all configured database hosts | `-` |
| `v-list-databases` | List all databases for a user | `-` |
| `v-list-database-types` | List available database types (MySQL, PostgreSQL, etc.) | `-` |
| `v-suspend-database` | Suspend a database | `-` |
| `v-suspend-database-host` | Suspend a database host | `-` |
| `v-suspend-databases` | Suspend all databases for a user | `-` |
| `v-unsuspend-database` | Unsuspend a database | `-` |
| `v-unsuspend-database-host` | Unsuspend a database host | `-` |
| `v-unsuspend-databases` | Unsuspend all databases for a user | `-` |
| `v-update-database-disk` | Update disk usage for a specific database | `-` |
| `v-update-databases-disk` | Update disk usage for all databases | `-` |

---

## Mail

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-mail-account` | Add a new mail account to a mail domain | `hestiacp_email_account` (create) |
| `v-add-mail-account-alias` | Add an alias to a mail account | `hestiacp_mail_alias` (create) |
| `v-add-mail-account-autoreply` | Set an autoreply message for a mail account | `hestiacp_mail_autoreply` (create) |
| `v-add-mail-account-forward` | Add a forwarding address to a mail account | `hestiacp_mail_forward` (create) |
| `v-add-mail-account-fwd-only` | Set a mail account to forward-only mode | `-` |
| `v-add-mail-domain` | Add a new mail domain | `hestiacp_email_domain` (create) |
| `v-add-mail-domain-antispam` | Enable antispam filtering for a mail domain | `hestiacp_mail_domain_antispam` (create) — enabled by default on new domains |
| `v-add-mail-domain-antivirus` | Enable antivirus scanning for a mail domain | `hestiacp_mail_domain_antivirus` (create) — enabled by default on new domains |
| `v-add-mail-domain-catchall` | Set a catch-all address for a mail domain | `hestiacp_mail_domain_catchall` (create) |
| `v-add-mail-domain-dkim` | Enable DKIM signing for a mail domain | `hestiacp_mail_domain_dkim` (create) — enabled by default on new domains |
| `v-add-mail-domain-reject` | Enable spam rejection for a mail domain | `-` |
| `v-add-mail-domain-smtp-relay` | Configure an SMTP relay for a mail domain | `-` |
| `v-add-mail-domain-ssl` | Add an SSL certificate to a mail domain | `-` |
| `v-add-mail-domain-webmail` | Enable webmail access for a mail domain | `-` |
| `v-change-mail-account-password` | Change the password for a mail account | `-` |
| `v-change-mail-account-quota` | Change the storage quota for a mail account | `-` |
| `v-change-mail-account-rate-limit` | Change the rate limit for a mail account | `-` |
| `v-change-mail-domain-catchall` | Change the catch-all address for a mail domain | `hestiacp_mail_domain_catchall` (update) |
| `v-change-mail-domain-rate-limit` | Change the sending rate limit for a mail domain | `-` |
| `v-change-mail-domain-sslcert` | Replace the SSL certificate for a mail domain | `-` |
| `v-check-mail-account-hash` | Verify the password hash for a mail account | `-` |
| `v-delete-mail-account` | Delete a mail account | `hestiacp_email_account` (delete) |
| `v-delete-mail-account-alias` | Remove an alias from a mail account | `hestiacp_mail_alias` (delete) |
| `v-delete-mail-account-autoreply` | Remove the autoreply message from a mail account | `hestiacp_mail_autoreply` (delete) |
| `v-delete-mail-account-forward` | Remove a forwarding address from a mail account | `hestiacp_mail_forward` (delete) |
| `v-delete-mail-account-fwd-only` | Remove forward-only mode from a mail account | `-` |
| `v-delete-mail-domain` | Delete a mail domain | `hestiacp_email_domain` (delete) |
| `v-delete-mail-domain-antispam` | Disable antispam filtering for a mail domain | `hestiacp_mail_domain_antispam` (delete) |
| `v-delete-mail-domain-antivirus` | Disable antivirus scanning for a mail domain | `hestiacp_mail_domain_antivirus` (delete) |
| `v-delete-mail-domain-catchall` | Remove the catch-all address from a mail domain | `hestiacp_mail_domain_catchall` (delete) |
| `v-delete-mail-domain-dkim` | Disable DKIM signing for a mail domain | `hestiacp_mail_domain_dkim` (delete) |
| `v-delete-mail-domain-reject` | Disable spam rejection for a mail domain | `-` |
| `v-delete-mail-domains` | Delete all mail domains for a user | `-` |
| `v-delete-mail-domain-smtp-relay` | Remove the SMTP relay from a mail domain | `-` |
| `v-delete-mail-domain-ssl` | Remove the SSL certificate from a mail domain | `-` |
| `v-delete-mail-domain-webmail` | Disable webmail access for a mail domain | `-` |
| `v-delete-sys-mail-queue` | Delete the system mail queue | `-` |
| `v-get-mail-account-value` | Get a specific configuration value for a mail account | `-` |
| `v-get-mail-domain-value` | Get a specific configuration value for a mail domain | `-` |
| `v-list-mail-account` | List details for a specific mail account | `hestiacp_email_account` (read) |
| `v-list-mail-account-autoreply` | List the autoreply configuration for a mail account | `hestiacp_mail_autoreply` (read) |
| `v-list-mail-accounts` | List all mail accounts for a mail domain | `-` |
| `v-list-mail-domain` | List details for a specific mail domain | `hestiacp_email_domain` (read) |
| `v-list-mail-domain-dkim` | List DKIM key details for a mail domain | `-` |
| `v-list-mail-domain-dkim-dns` | List DKIM DNS records for a mail domain | `-` |
| `v-list-mail-domains` | List all mail domains for a user | `-` |
| `v-list-mail-domain-ssl` | List SSL certificate details for a mail domain | `-` |
| `v-suspend-mail-account` | Suspend a mail account | `-` |
| `v-suspend-mail-accounts` | Suspend all mail accounts on a mail domain | `-` |
| `v-suspend-mail-domain` | Suspend a mail domain | `-` |
| `v-suspend-mail-domains` | Suspend all mail domains for a user | `-` |
| `v-unsuspend-mail-account` | Unsuspend a mail account | `-` |
| `v-unsuspend-mail-accounts` | Unsuspend all mail accounts on a mail domain | `-` |
| `v-unsuspend-mail-domain` | Unsuspend a mail domain | `-` |
| `v-unsuspend-mail-domains` | Unsuspend all mail domains for a user | `-` |
| `v-update-mail-domain-disk` | Update disk usage for a mail domain | `-` |
| `v-update-mail-domain-ssl` | Update the SSL certificate for a mail domain | `-` |
| `v-update-mail-domains-disk` | Update disk usage for all mail domains | `-` |
| `v-update-mail-templates` | Update mail server configuration templates | `-` |

---

## SSL / Let's Encrypt

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-letsencrypt-domain` | Issue a Let's Encrypt SSL certificate for a domain | `hestiacp_ssl` (create) |
| `v-add-letsencrypt-host` | Issue a Let's Encrypt certificate for the control panel host | `-` |
| `v-add-letsencrypt-user` | Register a user account with Let's Encrypt | `-` |
| `v-delete-letsencrypt-domain` | Remove a Let's Encrypt SSL certificate from a domain | `-` |
| `v-delete-ssl` | Remove an SSL certificate from a domain | `hestiacp_ssl` (delete) |
| `v-delete-web-domain-ssl` | Remove the SSL certificate from a web domain | `-` |
| `v-generate-ssl-cert` | Generate a self-signed SSL certificate | `-` |
| `v-list-letsencrypt-user` | List the Let's Encrypt account details for a user | `-` |
| `v-schedule-letsencrypt-domain` | Schedule Let's Encrypt certificate renewal for a domain | `-` |
| `v-update-host-certificate` | Update the SSL certificate for the control panel host | `-` |
| `v-update-letsencrypt-ssl` | Renew all Let's Encrypt SSL certificates | `-` |

---

## Backups

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-backup-host` | Add a remote backup host | `hestiacp_backup_host` (create) |
| `v-add-backup-host-restic` | Add a Restic remote backup host | `-` |
| `v-backup-user` | Create a backup for a user | `hestiacp_backup` (create) |
| `v-backup-user-config` | Back up the configuration files for a user | `-` |
| `v-backup-user-restic` | Create a Restic backup for a user | `-` |
| `v-backup-users` | Create backups for all users | `-` |
| `v-backup-users-restic` | Create Restic backups for all users | `-` |
| `v-delete-backup-host` | Remove a remote backup host | `hestiacp_backup_host` (delete) |
| `v-delete-backup-host-restic` | Remove a Restic remote backup host | `-` |
| `v-delete-user-backup` | Delete a specific user backup archive | `-` |
| `v-delete-user-backup-restic` | Delete a specific Restic backup snapshot | `-` |
| `v-download-backup` | Download a backup archive | `-` |
| `v-list-backup-host` | List details for a backup host | `hestiacp_backup_host` (read) |
| `v-list-backup-host-restic` | List details for a Restic backup host | `-` |
| `v-list-user-backup` | List details for a specific user backup | `-` |
| `v-list-user-backup-restic` | List details for a specific Restic backup snapshot | `-` |
| `v-list-user-backups` | List all backups for a user | `hestiacp_backup` (read) |
| `v-list-user-backups-restic` | List all Restic backup snapshots for a user | `-` |
| `v-list-user-files-restic` | List files within a Restic backup snapshot | `-` |
| `v-schedule-user-backup` | Schedule a backup for a user | `-` |
| `v-schedule-user-backup-download` | Schedule a backup download for a user | `-` |
| `v-schedule-user-backup-restic` | Schedule a Restic backup for a user | `-` |

---

## Cron Jobs

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-cron-hestia-autoupdate` | Add a cron job for HestiaCP automatic updates | `-` |
| `v-add-cron-job` | Add a new cron job for a user | `hestiacp_cron_job` (create) |
| `v-add-cron-letsencrypt-job` | Add a cron job for Let's Encrypt certificate renewal | `-` |
| `v-add-cron-reports` | Add a cron job to send system reports | `-` |
| `v-add-cron-restart-job` | Add a cron job to restart services | `-` |
| `v-change-cron-job` | Modify an existing cron job | `hestiacp_cron_job` (update) |
| `v-delete-cron-hestia-autoupdate` | Remove the HestiaCP autoupdate cron job | `-` |
| `v-delete-cron-job` | Delete a cron job | `hestiacp_cron_job` (delete) |
| `v-delete-cron-reports` | Remove the cron job for system reports | `-` |
| `v-delete-cron-restart-job` | Remove the service restart cron job | `-` |
| `v-list-cron-job` | List details for a specific cron job | `hestiacp_cron_job` (read) |
| `v-list-cron-jobs` | List all cron jobs for a user | `hestiacp_cron_job` (read) |
| `v-rebuild-cron-jobs` | Rebuild crontab from stored configuration | `-` |
| `v-restart-cron` | Restart the cron daemon | `-` |
| `v-restore-cron-job` | Restore a cron job from a backup | `-` |
| `v-restore-cron-job-restic` | Restore a cron job from a Restic backup | `-` |
| `v-suspend-cron-job` | Suspend a cron job | `-` |
| `v-suspend-cron-jobs` | Suspend all cron jobs for a user | `-` |
| `v-unsuspend-cron-job` | Unsuspend a cron job | `-` |
| `v-unsuspend-cron-jobs` | Unsuspend all cron jobs for a user | `-` |

---

## Firewall

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-firewall-ban` | Ban an IP address via the firewall | `-` |
| `v-add-firewall-chain` | Add a custom firewall chain | `-` |
| `v-add-firewall-ipset` | Add an IP set to the firewall | `-` |
| `v-add-firewall-rule` | Add a firewall rule | `hestiacp_firewall_rule` (create) — arg order: ACTION IPV4_CIDR PORT [PROTOCOL] [COMMENT] [RULE] |
| `v-add-sys-firewall` | Enable the system firewall | `-` |
| `v-change-firewall-rule` | Modify an existing firewall rule | `hestiacp_firewall_rule` (update) |
| `v-delete-firewall-ban` | Remove a firewall ban for an IP address | `-` |
| `v-delete-firewall-chain` | Delete a custom firewall chain | `-` |
| `v-delete-firewall-ipset` | Remove an IP set from the firewall | `-` |
| `v-delete-firewall-rule` | Delete a firewall rule | `hestiacp_firewall_rule` (delete) |
| `v-delete-sys-firewall` | Disable the system firewall | `-` |
| `v-list-firewall` | List all firewall rules and bans | `-` |
| `v-list-firewall-ban` | List all currently banned IP addresses | `-` |
| `v-list-firewall-ipset` | List all IP sets in the firewall | `-` |
| `v-list-firewall-rule` | List a specific firewall rule | `hestiacp_firewall_rule` (read) |
| `v-stop-firewall` | Stop the firewall service | `-` |
| `v-suspend-firewall-rule` | Suspend a firewall rule | `-` |
| `v-unsuspend-firewall-rule` | Unsuspend a firewall rule | `-` |
| `v-update-firewall` | Apply and reload all firewall rules | `-` |
| `v-update-firewall-ipset` | Update an IP set in the firewall | `-` |

---

## File System

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-fs-archive` | Create an archive of a file or directory | `-` |
| `v-add-fs-directory` | Create a directory on the file system | `-` |
| `v-add-fs-file` | Create a new file on the file system | `-` |
| `v-change-fs-file-permission` | Change permissions on a file or directory | `-` |
| `v-check-fs-permission` | Check if a path is accessible with given permissions | `-` |
| `v-copy-fs-directory` | Copy a directory on the file system | `-` |
| `v-copy-fs-file` | Copy a file on the file system | `-` |
| `v-delete-fs-directory` | Delete a directory from the file system | `-` |
| `v-delete-fs-file` | Delete a file from the file system | `-` |
| `v-extract-fs-archive` | Extract an archive on the file system | `-` |
| `v-get-fs-file-type` | Get the MIME type of a file | `-` |
| `v-list-fs-directory` | List the contents of a directory | `-` |
| `v-move-fs-directory` | Move or rename a directory | `-` |
| `v-move-fs-file` | Move or rename a file | `-` |
| `v-open-fs-config` | Open and display a system configuration file | `-` |
| `v-open-fs-file` | Open and display a file | `-` |
| `v-search-fs-object` | Search for a file or directory on the file system | `-` |

---

## Access Keys

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-access-key` | Create a new API access key | `-` |
| `v-check-access-key` | Verify an API access key | `-` |
| `v-check-api-key` | Check the validity of an API key | `-` |
| `v-delete-access-key` | Delete an API access key | `-` |
| `v-generate-api-key` | Generate a new API key | `-` |
| `v-list-access-key` | List details for a specific access key | `-` |
| `v-list-access-keys` | List all API access keys | `-` |
| `v-list-api` | List details for a specific API definition | `-` |
| `v-list-apis` | List all available API definitions | `-` |
| `v-revoke-api-key` | Revoke an existing API key | `-` |

---

## Remote DNS

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-remote-dns-domain` | Add a DNS zone to a remote DNS cluster node | `-` |
| `v-add-remote-dns-host` | Add a remote DNS cluster host | `hestiacp_remote_dns_host` (create) |
| `v-add-remote-dns-record` | Add a DNS record to a remote DNS cluster zone | `-` |
| `v-change-remote-dns-domain-exp` | Change the expiry date on a remote DNS zone | `-` |
| `v-change-remote-dns-domain-soa` | Change the SOA record on a remote DNS zone | `-` |
| `v-change-remote-dns-domain-ttl` | Change the TTL on a remote DNS zone | `-` |
| `v-delete-remote-dns-domain` | Delete a zone from a remote DNS cluster node | `-` |
| `v-delete-remote-dns-domains` | Delete all zones from a remote DNS cluster node | `-` |
| `v-delete-remote-dns-host` | Remove a remote DNS cluster host | `hestiacp_remote_dns_host` (delete) |
| `v-delete-remote-dns-record` | Delete a record from a remote DNS cluster zone | `-` |
| `v-list-remote-dns-hosts` | List all configured remote DNS cluster hosts | `hestiacp_remote_dns_host` (read) |
| `v-suspend-remote-dns-host` | Suspend a remote DNS cluster host | `-` |
| `v-unsuspend-remote-dns-host` | Unsuspend a remote DNS cluster host | `-` |

---

## System

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-add-sys-api-ip` | Add an IP address to the API access whitelist | `-` |
| `v-add-sys-cgroups` | Enable cgroup resource limiting for the system | `-` |
| `v-add-sys-dependencies` | Install required system dependencies | `-` |
| `v-add-sys-filemanager` | Install the file manager component | `-` |
| `v-add-sys-ip` | Add a system IP address | `-` |
| `v-add-sys-pma-sso` | Enable phpMyAdmin single sign-on | `-` |
| `v-add-sys-quota` | Enable disk quotas on the system | `-` |
| `v-add-sys-roundcube` | Install the Roundcube webmail client | `-` |
| `v-add-sys-sftp-jail` | Enable global SFTP jail | `-` |
| `v-add-sys-smtp` | Configure the system SMTP relay | `-` |
| `v-add-sys-smtp-relay` | Add an SMTP relay configuration | `-` |
| `v-add-sys-snappymail` | Install the SnappyMail webmail client | `-` |
| `v-add-sys-ssh-jail` | Enable SSH jail for the system | `-` |
| `v-add-sys-web-terminal` | Install the web terminal component | `-` |
| `v-change-sys-api` | Enable or disable the system API | `-` |
| `v-change-sys-config-value` | Change a value in the main system configuration | `-` |
| `v-change-sys-db-alias` | Change the database host alias | `-` |
| `v-change-sys-demo-mode` | Enable or disable demo mode | `-` |
| `v-change-sys-hestia-ssl` | Change the control panel SSL certificate | `-` |
| `v-change-sys-hostname` | Change the server hostname | `-` |
| `v-change-sys-ip-name` | Change the hostname associated with a system IP | `-` |
| `v-change-sys-ip-nat` | Configure NAT for a system IP address | `-` |
| `v-change-sys-ip-owner` | Change the owner of a system IP address | `-` |
| `v-change-sys-ip-status` | Change the status of a system IP address | `-` |
| `v-change-sys-language` | Change the default system language | `-` |
| `v-change-sys-php` | Change the default PHP version for the system | `-` |
| `v-change-sys-port` | Change the control panel port number | `-` |
| `v-change-sys-release` | Switch to a different HestiaCP release channel | `-` |
| `v-change-sys-service-config` | Update the configuration for a system service | `-` |
| `v-change-sys-timezone` | Change the server timezone | `-` |
| `v-change-sys-webmail` | Change the default webmail client | `-` |
| `v-change-sys-web-terminal-port` | Change the web terminal port | `-` |
| `v-delete-sys-api-ip` | Remove an IP from the API access whitelist | `-` |
| `v-delete-sys-cgroups` | Disable cgroup resource limiting | `-` |
| `v-delete-sys-filemanager` | Remove the file manager component | `-` |
| `v-delete-sys-ip` | Remove a system IP address | `-` |
| `v-delete-sys-pma-sso` | Disable phpMyAdmin single sign-on | `-` |
| `v-delete-sys-quota` | Disable disk quotas on the system | `-` |
| `v-delete-sys-sftp-jail` | Disable global SFTP jail | `-` |
| `v-delete-sys-smtp` | Remove the system SMTP relay configuration | `-` |
| `v-delete-sys-smtp-relay` | Remove the SMTP relay configuration | `-` |
| `v-delete-sys-snappymail` | Remove the SnappyMail webmail client | `-` |
| `v-delete-sys-ssh-jail` | Disable SSH jail for the system | `-` |
| `v-delete-sys-web-terminal` | Remove the web terminal component | `-` |
| `v-export-rrd` | Export RRD performance graph data | `-` |
| `v-get-sys-timezone` | Get the current server timezone | `-` |
| `v-get-sys-timezones` | List all available timezones | `-` |
| `v-import-cpanel` | Import data from a cPanel backup | `-` |
| `v-import-directadmin` | Import data from a DirectAdmin backup | `-` |
| `v-list-sys-clamd-config` | List the ClamAV daemon configuration | `-` |
| `v-list-sys-config` | List the main system configuration | `-` |
| `v-list-sys-cpu-status` | List current CPU usage statistics | `-` |
| `v-list-sys-db-status` | List database service status | `-` |
| `v-list-sys-disk-status` | List disk usage status | `-` |
| `v-list-sys-dns-status` | List DNS service status | `-` |
| `v-list-sys-dovecot-config` | List the Dovecot IMAP/POP3 configuration | `-` |
| `v-list-sys-hestia-autoupdate` | List the HestiaCP autoupdate configuration | `-` |
| `v-list-sys-hestia-ssl` | List the control panel SSL certificate details | `-` |
| `v-list-sys-hestia-updates` | List available HestiaCP updates | `-` |
| `v-list-sys-info` | List general system information | `-` |
| `v-list-sys-interfaces` | List network interfaces on the system | `-` |
| `v-list-sys-ip` | List details for a specific system IP | `-` |
| `v-list-sys-ips` | List all system IP addresses | `-` |
| `v-list-sys-languages` | List all available languages | `-` |
| `v-list-sys-mail-status` | List mail service status | `-` |
| `v-list-sys-memory-status` | List current memory usage statistics | `-` |
| `v-list-sys-mysql-config` | List the MySQL server configuration | `-` |
| `v-list-sys-network-status` | List network interface statistics | `-` |
| `v-list-sys-nginx-config` | List the nginx web server configuration | `-` |
| `v-list-sys-pgsql-config` | List the PostgreSQL server configuration | `-` |
| `v-list-sys-php` | List installed PHP versions | `-` |
| `v-list-sys-php-config` | List the PHP configuration for a given version | `-` |
| `v-list-sys-proftpd-config` | List the ProFTPD configuration | `-` |
| `v-list-sys-rrd` | List available RRD performance graphs | `-` |
| `v-list-sys-services` | List all system services and their status | `-` |
| `v-list-sys-shells` | List available login shells | `-` |
| `v-list-sys-spamd-config` | List the SpamAssassin daemon configuration | `-` |
| `v-list-sys-sshd-port` | List the current SSH daemon port | `-` |
| `v-list-sys-themes` | List available control panel themes | `-` |
| `v-list-sys-users` | List system (OS-level) users | `-` |
| `v-list-sys-vsftpd-config` | List the vsftpd FTP server configuration | `-` |
| `v-list-sys-webmail` | List webmail client configuration | `-` |
| `v-list-sys-web-status` | List web server service status | `-` |
| `v-refresh-sys-theme` | Refresh and reapply the control panel theme | `-` |
| `v-repair-sys-config` | Repair the system configuration | `-` |
| `v-run-cli-cmd` | Execute an arbitrary CLI command via the API | `-` |
| `v-start-service` | Start a system service | `-` |
| `v-stop-service` | Stop a system service | `-` |
| `v-update-sys-defaults` | Reset system configuration to default values | `-` |
| `v-update-sys-hestia` | Update HestiaCP to the latest version | `-` |
| `v-update-sys-hestia-all` | Update HestiaCP and all components | `-` |
| `v-update-sys-hestia-git` | Update HestiaCP from a Git repository | `-` |
| `v-update-sys-ip` | Update system IP configuration | `-` |
| `v-update-sys-ip-counters` | Update traffic counters for system IPs | `-` |
| `v-update-sys-queue` | Process the system task queue | `-` |
| `v-update-sys-rrd` | Update all RRD performance graphs | `-` |
| `v-update-sys-rrd-apache2` | Update the Apache2 RRD performance graph | `-` |
| `v-update-sys-rrd-ftp` | Update the FTP RRD performance graph | `-` |
| `v-update-sys-rrd-httpd` | Update the HTTPD RRD performance graph | `-` |
| `v-update-sys-rrd-la` | Update the load average RRD performance graph | `-` |
| `v-update-sys-rrd-mail` | Update the mail RRD performance graph | `-` |
| `v-update-sys-rrd-mem` | Update the memory RRD performance graph | `-` |
| `v-update-sys-rrd-mysql` | Update the MySQL RRD performance graph | `-` |
| `v-update-sys-rrd-net` | Update the network RRD performance graph | `-` |
| `v-update-sys-rrd-nginx` | Update the nginx RRD performance graph | `-` |
| `v-update-sys-rrd-pgsql` | Update the PostgreSQL RRD performance graph | `-` |
| `v-update-sys-rrd-ssh` | Update the SSH RRD performance graph | `-` |
| `v-update-white-label-logo` | Update the white-label logo for the control panel | `-` |

---

## Restore / Rebuild / Restart

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-rebuild-all` | Rebuild all configuration for a user | `-` |
| `v-rebuild-database` | Rebuild configuration for a specific database | `-` |
| `v-rebuild-databases` | Rebuild configuration for all databases | `-` |
| `v-rebuild-dns-domain` | Rebuild configuration for a specific DNS zone | `-` |
| `v-rebuild-dns-domains` | Rebuild configuration for all DNS zones | `-` |
| `v-rebuild-mail-domain` | Rebuild configuration for a specific mail domain | `-` |
| `v-rebuild-mail-domains` | Rebuild configuration for all mail domains | `-` |
| `v-rebuild-user` | Rebuild all configuration for a specific user | `-` |
| `v-rebuild-users` | Rebuild all configuration for all users | `-` |
| `v-rebuild-web-domain` | Rebuild configuration for a specific web domain | `-` |
| `v-rebuild-web-domains` | Rebuild configuration for all web domains | `-` |
| `v-restart-dns` | Restart the DNS server | `-` |
| `v-restart-ftp` | Restart the FTP server | `-` |
| `v-restart-mail` | Restart the mail server | `-` |
| `v-restart-proxy` | Restart the proxy server | `-` |
| `v-restart-service` | Restart a named system service | `-` |
| `v-restart-system` | Restart the entire system | `-` |
| `v-restart-web` | Restart the web server | `-` |
| `v-restart-web-backend` | Restart the web backend (PHP-FPM) | `-` |
| `v-restore-database` | Restore a database from a backup | `-` |
| `v-restore-database-restic` | Restore a database from a Restic backup | `-` |
| `v-restore-dns-domain` | Restore a DNS zone from a backup | `-` |
| `v-restore-dns-domain-restic` | Restore a DNS zone from a Restic backup | `-` |
| `v-restore-file-restic` | Restore a specific file from a Restic backup | `-` |
| `v-restore-mail-domain` | Restore a mail domain from a backup | `-` |
| `v-restore-mail-domain-restic` | Restore a mail domain from a Restic backup | `-` |
| `v-restore-user` | Restore a full user account from a backup | `-` |
| `v-restore-user-full-restic` | Restore a full user account from a Restic backup | `-` |
| `v-restore-user-restic` | Restore a user from a Restic backup | `-` |
| `v-restore-web-domain` | Restore a web domain from a backup | `-` |
| `v-restore-web-domain-restic` | Restore a web domain from a Restic backup | `-` |
| `v-schedule-user-restore` | Schedule a user restore operation | `-` |
| `v-schedule-user-restore-restic` | Schedule a Restic-based user restore operation | `-` |

---

## Search / Logging / Misc

| Command | Description | Terraform Resource |
|---------|-------------|-------------------|
| `v-generate-password-hash` | Generate a hashed password string | `-` |
| `v-log-action` | Log a custom action entry | `-` |
| `v-log-user-login` | Log a user login event | `-` |
| `v-log-user-logout` | Log a user logout event | `-` |
| `v-search-command` | Search for available CLI commands by keyword | `-` |
| `v-search-domain-owner` | Find the user who owns a given domain | `-` |
| `v-search-object` | Search for any object across all users | `-` |
| `v-search-user-object` | Search for an object within a specific user's account | `-` |
