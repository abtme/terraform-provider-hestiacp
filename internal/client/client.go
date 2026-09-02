package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ErrNotFound is returned by Read functions when HestiaCP reports the object does not exist.
var ErrNotFound = errors.New("object does not exist")

// ReturnCode maps HestiaCP numeric return codes to error strings.
var ReturnCode = map[int]string{
	0:  "OK",
	1:  "E_ARGS: not enough arguments",
	2:  "E_INVALID: object or argument is not valid",
	3:  "E_NOTEXIST: object does not exist",
	4:  "E_EXISTS: object already exists",
	5:  "E_SUSPENDED: object is already suspended",
	6:  "E_UNSUSPENDED: object is already unsuspended",
	7:  "E_INUSE: object is in use by another object",
	8:  "E_LIMIT: hosting package limit reached",
	9:  "E_PASSWORD: wrong password",
	10: "E_FORBIDDEN: access denied",
	11: "E_DISABLED: subsystem is disabled",
	12: "E_PARSING: configuration is broken",
	13: "E_DISK: not enough disk space",
	14: "E_LA: server too busy",
	15: "E_CONNECT: connection failed",
	17: "E_DB: database server not responding",
	19: "E_UPDATE: update failed",
	20: "E_RESTART: service restart failed",
}

// Client is the HestiaCP API client.
// Auth uses ACCESS_KEY:SECRET_KEY (v1.6+ preferred method).
type Client struct {
	baseURL      string // e.g. https://myserver.com:8083
	accessKey    string // ACCESS_KEY:SECRET_KEY
	username     string // default HestiaCP user for resource operations
	httpClient   *http.Client
	sslSemaphore chan struct{} // limits concurrent SSL/LE certificate operations
}

// New creates a new HestiaCP client.
// sslConcurrency caps how many SSL issuance operations run simultaneously
// (Let's Encrypt and HestiaCP both struggle with many concurrent requests).
func New(baseURL, accessKey, username string) *Client {
	return &Client{
		baseURL:   strings.TrimSuffix(baseURL, "/"),
		accessKey: accessKey,
		username:  username,
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
		sslSemaphore: make(chan struct{}, 1), // max 1 concurrent SSL operation (sequential)
	}
}

// DefaultUser returns the configured default HestiaCP username.
func (c *Client) DefaultUser() string {
	return c.username
}

// apiURL returns the full API endpoint.
func (c *Client) apiURL() string {
	return c.baseURL + "/api/"
}

// retryBackoff is the sequence of delays between retry attempts.
// Max attempts = 1 (initial) + len(retryBackoff) = 6 total.
var retryBackoff = []time.Duration{
	5 * time.Second,
	10 * time.Second,
	20 * time.Second,
	40 * time.Second,
	80 * time.Second,
}

// isRetryableCode returns true for HestiaCP return codes that may be transient.
func isRetryableCode(rc string) bool {
	var code int
	fmt.Sscanf(rc, "%d", &code)
	return code == 2 || code == 3 // E_INVALID or E_NOTEXIST
}

// do executes a POST to the HestiaCP API with form-encoded body.
// cmd  = v-add-user, v-list-web-domains, etc.
// args = positional arg1..argN values.
// Returns the raw response body. Automatically retries on E_INVALID (2) and
// E_NOTEXIST (3) with exponential backoff (5/10/20/40/80s, max 6 attempts).
func (c *Client) do(cmd string, args ...string) (string, error) {
	rc, err := c.doOnce(cmd, args...)
	if err != nil {
		return rc, err
	}
	for _, delay := range retryBackoff {
		if !isRetryableCode(rc) {
			break
		}
		time.Sleep(delay)
		rc, err = c.doOnce(cmd, args...)
		if err != nil {
			return rc, err
		}
	}
	return rc, nil
}

// debugHTTP reports whether HESTIACP_DEBUG_HTTP is set, enabling verbose
// request/response logging to stderr (visible via TF_LOG=debug or above,
// since the plugin SDK forwards plugin stderr into Terraform/OpenTofu's own
// log stream). The access key is always redacted.
func debugHTTP() bool {
	return os.Getenv("HESTIACP_DEBUG_HTTP") != ""
}

// doOnce makes a single POST to the HestiaCP API without retry.
func (c *Client) doOnce(cmd string, args ...string) (string, error) {
	form := url.Values{}
	form.Set("hash", c.accessKey)
	form.Set("cmd", cmd)
	form.Set("returncode", "yes")

	for i, a := range args {
		form.Set(fmt.Sprintf("arg%d", i+1), a)
	}

	if debugHTTP() {
		redacted := url.Values{}
		for k, v := range form {
			redacted[k] = v
		}
		redacted.Set("hash", "[redacted]")
		fmt.Fprintf(os.Stderr, "HESTIACP_DEBUG_HTTP: request cmd=%s form=%q\n", cmd, redacted.Encode())
	}

	resp, err := c.httpClient.Post(c.apiURL(), "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		if debugHTTP() {
			fmt.Fprintf(os.Stderr, "HESTIACP_DEBUG_HTTP: request failed: %v\n", err)
		}
		return "", fmt.Errorf("HestiaCP API request failed: %w", err)
	}
	defer resp.Body.Close()
	if debugHTTP() {
		fmt.Fprintf(os.Stderr, "HESTIACP_DEBUG_HTTP: response status=%s\n", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HestiaCP response: %w", err)
	}

	if debugHTTP() {
		fmt.Fprintf(os.Stderr, "HESTIACP_DEBUG_HTTP: response body=%q\n", strings.TrimSpace(string(body)))
	}

	return strings.TrimSpace(string(body)), nil
}

// doJSON executes a POST to the HestiaCP API and unmarshals a JSON response.
func (c *Client) doJSON(cmd string, out interface{}, args ...string) error {
	form := url.Values{}
	form.Set("hash", c.accessKey)
	form.Set("cmd", cmd)
	form.Set("returncode", "no")

	for i, a := range args {
		form.Set(fmt.Sprintf("arg%d", i+1), a)
	}
	// request JSON output
	form.Set("arg"+fmt.Sprintf("%d", len(args)+1), "json")

	resp, err := c.httpClient.Post(c.apiURL(), "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("HestiaCP API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HestiaCP response: %w", err)
	}

	if err := json.Unmarshal(body, out); err != nil {
		bodyStr := strings.TrimSpace(string(body))
		if strings.HasPrefix(bodyStr, "Error:") && strings.Contains(bodyStr, "doesn't exist") {
			return ErrNotFound
		}
		return fmt.Errorf("failed to parse HestiaCP JSON response (%s): %w", bodyStr, err)
	}
	return nil
}

// reloadWebServer forces a real, synchronous, graceful nginx reload via the
// HestiaCP API, bypassing HestiaCP's own deferred restart machinery for
// callers where that machinery is known to be unreliable.
//
// This only matters for the SSL/Let's Encrypt path. Traced on srv2
// 2026-09-02: v-add-web-domain, its alias/redirect siblings, and their
// delete counterparts all call v-restart-web/v-restart-proxy with an empty
// restart mode at the end of every invocation regardless of success —
// which (confirmed live: a brand-new domain served HestiaCP's correct
// per-domain "Coming Soon" placeholder immediately after v-add-web-domain
// returned, with zero extra calls) already reloads nginx synchronously and
// correctly. Those call sites do NOT need any help from this function.
//
// v-add-letsencrypt-domain is different: internally it calls
// v-add-web-domain-ssl/v-update-web-domain-ssl with restart mode
// "updatessl" specifically (see v-add-letsencrypt-domain around its
// "updatessl" calls), which queues "v-restart-proxy ssl" onto the
// 2-minute cron queue instead of reloading inline, and even once that
// queue drains resolves to "service nginx upgrade" (nginx's graceful
// binary-upgrade path) — this combination is what caused the multi-day SSL
// issuance failure root-caused 2026-09-01 (see
// docs/nginx-restart-reliability.md if present, or the project history).
// A full `systemctl restart` was used to recover live in that incident,
// but nothing in that investigation showed a plain reload was
// insufficient — only that the queued/"upgrade" path was. A graceful
// reload is standard, well-established practice for picking up renewed
// certs (it's literally certbot's own default deploy-hook command) and
// avoids dropping connections on every other domain sharing this nginx,
// so use that here rather than a full restart.
//
// v-restart-service's own arg2 ("restart mode") is format-validated
// server-side to one of: "", "yes", "no", "ssl", "reload", "updatessl",
// "scheduled" (main.sh:is_restart_format_valid) — anything outside that
// set is rejected with E_INVALID. Passing "" here hits the
// `systemctl reload-or-restart nginx` branch, which — since nginx's own
// systemd unit defines ExecReload — is a plain graceful reload, not a
// restart; confirmed live (MainPID/ActiveEnterTimestamp unchanged across
// the call, unlike the full-restart case where both changed).
//
// Never defer this: the calling Create/Delete only returns success once
// this reload itself has succeeded, rather than trusting HestiaCP's own
// queued/deferred path for an operation whose own success depends on it.
func (c *Client) reloadWebServer() error {
	rc, err := c.do("v-restart-service", "nginx")
	if err != nil {
		return err
	}
	return checkRC(rc)
}

// checkRC converts a HestiaCP numeric return-code string into an error.
// rc = "0" means success; anything else is an error.
// allowedExtra are codes treated as non-fatal (e.g. 3=not-found on delete).
func checkRC(rc string, allowedExtra ...int) error {
	if rc == "0" {
		return nil
	}
	var code int
	fmt.Sscanf(rc, "%d", &code)
	for _, a := range allowedExtra {
		if code == a {
			return nil
		}
	}
	msg, ok := ReturnCode[code]
	if !ok {
		msg = "unknown error"
	}
	return fmt.Errorf("HestiaCP error %d: %s", code, msg)
}

// ── User ────────────────────────────────────────────────────────────────────

type User struct {
	Email     string `json:"CONTACT"`
	Package   string `json:"PACKAGE"`
	Name      string `json:"NAME"`
	Shell     string `json:"SHELL"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateUser(username, password, email, pkg, firstName string) error {
	rc, err := c.do("v-add-user", username, password, email, pkg, firstName)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ReadUser(username string) (*User, error) {
	var out map[string]User
	if err := c.doJSON("v-list-user", &out, username); err != nil {
		return nil, err
	}
	u, ok := out[username]
	if !ok {
		return nil, nil // not found
	}
	return &u, nil
}

func (c *Client) UpdateUserPassword(username, newPassword string) error {
	rc, err := c.do("v-change-user-password", username, newPassword)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) UpdateUserEmail(username, newEmail string) error {
	rc, err := c.do("v-change-user-contact", username, newEmail)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteUser(username string) error {
	rc, err := c.do("v-delete-user", username)
	if err != nil {
		return err
	}
	return checkRC(rc, 3) // 3 = already gone
}

// ── Web Domain ──────────────────────────────────────────────────────────────

type WebDomain struct {
	IP        string `json:"IP"`
	Template  string `json:"TEMPLATE"`
	SSL       string `json:"SSL"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateWebDomain(user, domain, ip string) error {
	rc, err := c.do("v-add-web-domain", user, domain, ip)
	if err != nil {
		return err
	}
	return checkRC(rc, 4) // 4 = already exists, treat as success
}

func (c *Client) ReadWebDomain(user, domain string) (*WebDomain, error) {
	var out map[string]WebDomain
	if err := c.doJSON("v-list-web-domain", &out, user, domain); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	d, ok := out[domain]
	if !ok {
		return nil, nil
	}
	return &d, nil
}

func (c *Client) DeleteWebDomain(user, domain string) error {
	rc, err := c.do("v-delete-web-domain", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── DNS Zone ────────────────────────────────────────────────────────────────

type DNSZone struct {
	IP        string `json:"IP"`
	Template  string `json:"TEMPLATE"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateDNSZone(user, domain, ip string) error {
	rc, err := c.do("v-add-dns-domain", user, domain, ip)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ReadDNSZone(user, domain string) (*DNSZone, error) {
	var out map[string]DNSZone
	if err := c.doJSON("v-list-dns-domain", &out, user, domain); err != nil {
		return nil, err
	}
	z, ok := out[domain]
	if !ok {
		return nil, nil
	}
	return &z, nil
}

func (c *Client) DeleteDNSZone(user, domain string) error {
	rc, err := c.do("v-delete-dns-domain", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── DNS Record ──────────────────────────────────────────────────────────────

type DNSRecord struct {
	Record   string `json:"RECORD"`
	Type     string `json:"TYPE"`
	Value    string `json:"VALUE"`
	Priority string `json:"PRIORITY"`
}

func (c *Client) CreateDNSRecord(user, domain, record, recType, value string, priority int) error {
	rc, err := c.do("v-add-dns-record", user, domain, record, recType, value, fmt.Sprintf("%d", priority))
	if err != nil {
		return err
	}
	return checkRC(rc)
}

// ListDNSRecords returns all DNS records for a zone.
func (c *Client) ListDNSRecords(user, domain string) (map[string]DNSRecord, error) {
	var out map[string]DNSRecord
	if err := c.doJSON("v-list-dns-records", &out, user, domain); err != nil {
		return nil, err
	}
	return out, nil
}

// FindDNSRecordID finds the numeric ID of a record matching record+type+value.
func (c *Client) FindDNSRecordID(user, domain, record, recType, value string) (string, error) {
	records, err := c.ListDNSRecords(user, domain)
	if err != nil {
		return "", err
	}
	for id, r := range records {
		if r.Record == record && r.Type == recType && r.Value == value {
			return id, nil
		}
	}
	return "", nil // not found
}

func (c *Client) DeleteDNSRecord(user, domain, id string) error {
	rc, err := c.do("v-delete-dns-record", user, domain, id)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Database ────────────────────────────────────────────────────────────────

type Database struct {
	DbUser    string `json:"DBUSER"`
	Host      string `json:"HOST"`
	Type      string `json:"TYPE"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateDatabase(user, dbName, dbUser, dbPass, dbType string) error {
	rc, err := c.do("v-add-database", user, dbName, dbUser, dbPass, dbType)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ReadDatabase(user, dbName string) (*Database, error) {
	fullName := user + "_" + dbName
	var out map[string]Database
	if err := c.doJSON("v-list-database", &out, user, fullName); err != nil {
		return nil, err
	}
	db, ok := out[fullName]
	if !ok {
		return nil, nil
	}
	return &db, nil
}

func (c *Client) DeleteDatabase(user, dbName string) error {
	fullName := user + "_" + dbName
	rc, err := c.do("v-delete-database", user, fullName)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Domain ─────────────────────────────────────────────────────────────

type MailDomain struct {
	Antispam  string `json:"ANTISPAM"`
	Antivirus string `json:"ANTIVIRUS"`
	DKIM      string `json:"DKIM"`
	Catchall  string `json:"CATCHALL"`
	Suspended string `json:"SUSPENDED"`
	SSL       string `json:"SSL"`
}

func (c *Client) CreateMailDomain(user, domain string) error {
	rc, err := c.do("v-add-mail-domain", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 4) // 4 = already exists, treat as success
}

func (c *Client) ReadMailDomain(user, domain string) (*MailDomain, error) {
	var out map[string]MailDomain
	if err := c.doJSON("v-list-mail-domain", &out, user, domain); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	md, ok := out[domain]
	if !ok {
		return nil, nil
	}
	return &md, nil
}

func (c *Client) DeleteMailDomain(user, domain string) error {
	rc, err := c.do("v-delete-mail-domain", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Account ────────────────────────────────────────────────────────────

type MailAccount struct {
	Quota     string `json:"QUOTA"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateMailAccount(user, domain, account, password string, quota int) error {
	rc, err := c.do("v-add-mail-account", user, domain, account, password, fmt.Sprintf("%d", quota))
	if err != nil {
		return err
	}
	return checkRC(rc, 4) // 4 = already exists, treat as success
}

func (c *Client) ReadMailAccount(user, domain, account string) (*MailAccount, error) {
	var out map[string]MailAccount
	if err := c.doJSON("v-list-mail-account", &out, user, domain, account); errors.Is(err, ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	ma, ok := out[account]
	if !ok {
		return nil, nil
	}
	return &ma, nil
}

func (c *Client) DeleteMailAccount(user, domain, account string) error {
	rc, err := c.do("v-delete-mail-account", user, domain, account)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── SSL ─────────────────────────────────────────────────────────────────────

func (c *Client) CreateSSL(user, domain, aliases string) error {
	c.sslSemaphore <- struct{}{} // acquire — ensures only one SSL op at a time
	defer func() { <-c.sslSemaphore }()

	// uses Let's Encrypt via v-add-letsencrypt-domain
	// Only pass aliases arg when non-empty; never pass the restart arg — passing
	// "yes" as arg3 when aliases is empty shifts it into the aliases position and
	// HestiaCP rejects "yes" as a non-existent web domain (E_NOTEXIST).
	args := []string{user, domain}
	if aliases != "" {
		args = append(args, aliases)
	}
	rc, err := c.do("v-add-letsencrypt-domain", args...)
	if err != nil {
		return err
	}
	if err := checkRC(rc, 4); err != nil { // 4 = already exists, treat as success
		return err
	}
	// This is the specific case that motivated reloadWebServer: HestiaCP's
	// own post-issuance restart for this call is queued/unreliable — see its
	// doc comment. Do not remove this without re-reading that history.
	return c.reloadWebServer()
}

func (c *Client) DeleteSSL(user, domain string) error {
	// Not "v-delete-ssl" - that command does not exist on HestiaCP (confirmed
	// live 2026-09-02: the API returns E_FORBIDDEN for it, meaning every
	// hestiacp_ssl destroy was silently failing). The real command is
	// v-delete-web-domain-ssl, options "USER DOMAIN [RESTART]".
	rc, err := c.do("v-delete-web-domain-ssl", user, domain)
	if err != nil {
		return err
	}
	if err := checkRC(rc, 3); err != nil {
		return err
	}
	return c.reloadWebServer()
}

func (c *Client) CreateMailSSL(user, domain string) error {
	c.sslSemaphore <- struct{}{} // acquire — ensures only one SSL op at a time
	defer func() { <-c.sslSemaphore }()

	// Delete first so we always copy the current web cert (not a stale copy).
	c.do("v-delete-mail-domain-ssl", user, domain) //nolint:errcheck
	certDir := fmt.Sprintf("/home/%s/conf/web/%s/ssl", user, domain)
	// Retry is handled automatically by do() on E_INVALID/E_NOTEXIST.
	rc, err := c.do("v-add-mail-domain-ssl", user, domain, certDir)
	if err != nil {
		return err
	}
	return checkRC(rc, 4)
}

func (c *Client) DeleteMailSSL(user, domain string) error {
	rc, err := c.do("v-delete-mail-domain-ssl", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3) // 3 = not found, treat as success
}

// ── Backup ──────────────────────────────────────────────────────────────────

func (c *Client) CreateBackup(user string) error {
	rc, err := c.do("v-backup-user", user)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

type Backup struct {
	Size string `json:"SIZE"`
	Date string `json:"DATE"`
	Time string `json:"TIME"`
}

func (c *Client) ListBackups(user string) (map[string]Backup, error) {
	var out map[string]Backup
	if err := c.doJSON("v-list-user-backups", &out, user); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) DeleteBackup(user, backup string) error {
	rc, err := c.do("v-delete-user-backup", user, backup)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Cron Job ─────────────────────────────────────────────────────────────────

type CronJob struct {
	Min       string `json:"MIN"`
	Hour      string `json:"HOUR"`
	Day       string `json:"DAY"`
	Month     string `json:"MONTH"`
	Wday      string `json:"WDAY"`
	Command   string `json:"CMD"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateCronJob(user, min, hour, day, month, wday, command string) error {
	rc, err := c.do("v-add-cron-job", user, min, hour, day, month, wday, command)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ListCronJobs(user string) (map[string]CronJob, error) {
	var out map[string]CronJob
	if err := c.doJSON("v-list-cron-jobs", &out, user); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) FindCronJobID(user, min, hour, day, month, wday, command string) (string, error) {
	jobs, err := c.ListCronJobs(user)
	if err != nil {
		return "", err
	}
	for id, j := range jobs {
		if j.Command == command && j.Min == min && j.Hour == hour &&
			j.Day == day && j.Month == month && j.Wday == wday {
			return id, nil
		}
	}
	return "", nil
}

func (c *Client) UpdateCronJob(user, id, min, hour, day, month, wday, command string) error {
	rc, err := c.do("v-change-cron-job", user, id, min, hour, day, month, wday, command)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteCronJob(user, id string) error {
	rc, err := c.do("v-delete-cron-job", user, id)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Firewall Rule ─────────────────────────────────────────────────────────────

type FirewallRule struct {
	Action    string `json:"ACTION"`
	Protocol  string `json:"PROTOCOL"`
	Port      string `json:"PORT"`
	IP        string `json:"IP"`
	Comment   string `json:"COMMENT"`
	Rule      string `json:"RULE"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateFirewallRule(action, protocol, port, ip, comment, rule string) error {
	rc, err := c.do("v-add-firewall-rule", action, ip, port, protocol, comment, rule)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ListFirewallRules() (map[string]FirewallRule, error) {
	var out map[string]FirewallRule
	if err := c.doJSON("v-list-firewall", &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) FindFirewallRuleID(action, protocol, port, ip string) (string, error) {
	rules, err := c.ListFirewallRules()
	if err != nil {
		return "", err
	}
	for id, r := range rules {
		if r.Action == action && r.Protocol == protocol && r.Port == port && r.IP == ip {
			return id, nil
		}
	}
	return "", nil
}

func (c *Client) UpdateFirewallRule(id, action, protocol, port, ip, comment, rule string) error {
	rc, err := c.do("v-change-firewall-rule", id, action, ip, port, protocol, comment, rule)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteFirewallRule(id string) error {
	rc, err := c.do("v-delete-firewall-rule", id)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Web Domain Alias ──────────────────────────────────────────────────────────

func (c *Client) CreateWebDomainAlias(user, domain, alias string) error {
	rc, err := c.do("v-add-web-domain-alias", user, domain, alias)
	if err != nil {
		return err
	}
	return checkRC(rc, 4) // 4 = already exists, alias already set
}

func (c *Client) DeleteWebDomainAlias(user, domain, alias string) error {
	rc, err := c.do("v-delete-web-domain-alias", user, domain, alias)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Web Domain FTP ────────────────────────────────────────────────────────────

func (c *Client) CreateWebDomainFTP(user, domain, ftpUser, ftpPassword, ftpPath string) error {
	rc, err := c.do("v-add-web-domain-ftp", user, domain, ftpUser, ftpPassword, ftpPath)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) UpdateWebDomainFTPPassword(user, domain, ftpUser, password string) error {
	rc, err := c.do("v-change-web-domain-ftp-password", user, domain, ftpUser, password)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteWebDomainFTP(user, domain, ftpUser string) error {
	rc, err := c.do("v-delete-web-domain-ftp", user, domain, ftpUser)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Web Domain Redirect ───────────────────────────────────────────────────────

func (c *Client) CreateWebDomainRedirect(user, domain, redirect, httpCode string) error {
	rc, err := c.do("v-add-web-domain-redirect", user, domain, redirect, httpCode)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteWebDomainRedirect(user, domain string) error {
	rc, err := c.do("v-delete-web-domain-redirect", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── SSH Key ───────────────────────────────────────────────────────────────────

type SSHKey struct {
	Fingerprint string `json:"KEY"` // SHA256:... fingerprint, not the public key
	Date        string `json:"DATE"`
}

func (c *Client) CreateSSHKey(user, key string) error {
	rc, err := c.do("v-add-user-ssh-key", user, key)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ListSSHKeys(user string) (map[string]SSHKey, error) {
	var out map[string]SSHKey
	if err := c.doJSON("v-list-user-ssh-key", &out, user); err != nil {
		return nil, err
	}
	return out, nil
}

// SSHKeyComment extracts the comment (3rd field) from a public key line.
// e.g. "ssh-ed25519 AAAA... my-key" → "my-key"
// HestiaCP uses the comment as the key identifier.
func SSHKeyComment(key string) string {
	parts := strings.Fields(key)
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func (c *Client) DeleteSSHKey(user, id string) error {
	rc, err := c.do("v-delete-user-ssh-key", user, id)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Forward ──────────────────────────────────────────────────────────────

func (c *Client) CreateMailForward(user, domain, account, forward string) error {
	rc, err := c.do("v-add-mail-account-forward", user, domain, account, forward)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteMailForward(user, domain, account, forward string) error {
	rc, err := c.do("v-delete-mail-account-forward", user, domain, account, forward)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Alias ────────────────────────────────────────────────────────────────

func (c *Client) CreateMailAlias(user, domain, account, alias string) error {
	rc, err := c.do("v-add-mail-account-alias", user, domain, account, alias)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteMailAlias(user, domain, account, alias string) error {
	rc, err := c.do("v-delete-mail-account-alias", user, domain, account, alias)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Autoreply ────────────────────────────────────────────────────────────

type MailAutoreply struct {
	Message string `json:"MSG"`
	Status  string `json:"STATUS"`
}

func (c *Client) CreateMailAutoreply(user, domain, account, message string) error {
	rc, err := c.do("v-add-mail-account-autoreply", user, domain, account, message)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ReadMailAutoreply(user, domain, account string) (*MailAutoreply, error) {
	var out map[string]MailAutoreply
	if err := c.doJSON("v-list-mail-account-autoreply", &out, user, domain, account); err != nil {
		return nil, err
	}
	a, ok := out[account]
	if !ok {
		return nil, nil
	}
	return &a, nil
}

func (c *Client) DeleteMailAutoreply(user, domain, account string) error {
	rc, err := c.do("v-delete-mail-account-autoreply", user, domain, account)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Web Domain HTTP Auth ──────────────────────────────────────────────────────

func (c *Client) CreateWebDomainHTTPAuth(user, domain, authUser, authPassword string) error {
	rc, err := c.do("v-add-web-domain-httpauth", user, domain, authUser, authPassword)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) UpdateWebDomainHTTPAuth(user, domain, authUser, authPassword string) error {
	rc, err := c.do("v-change-web-domain-httpauth", user, domain, authUser, authPassword)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteWebDomainHTTPAuth(user, domain, authUser string) error {
	rc, err := c.do("v-delete-web-domain-httpauth", user, domain, authUser)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Domain Catchall ──────────────────────────────────────────────────────

func (c *Client) CreateMailDomainCatchall(user, domain, email string) error {
	rc, err := c.do("v-add-mail-domain-catchall", user, domain, email)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) UpdateMailDomainCatchall(user, domain, email string) error {
	rc, err := c.do("v-change-mail-domain-catchall", user, domain, email)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteMailDomainCatchall(user, domain string) error {
	rc, err := c.do("v-delete-mail-domain-catchall", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Domain DKIM ──────────────────────────────────────────────────────────

func (c *Client) CreateMailDomainDKIM(user, domain string) error {
	rc, err := c.do("v-add-mail-domain-dkim", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 4) // 4 = already enabled
}

func (c *Client) DeleteMailDomainDKIM(user, domain string) error {
	rc, err := c.do("v-delete-mail-domain-dkim", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Domain Antispam ──────────────────────────────────────────────────────

func (c *Client) CreateMailDomainAntispam(user, domain string) error {
	rc, err := c.do("v-add-mail-domain-antispam", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 4) // 4 = already enabled
}

func (c *Client) DeleteMailDomainAntispam(user, domain string) error {
	rc, err := c.do("v-delete-mail-domain-antispam", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Mail Domain Antivirus ─────────────────────────────────────────────────────

func (c *Client) CreateMailDomainAntivirus(user, domain string) error {
	rc, err := c.do("v-add-mail-domain-antivirus", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 4) // 4 = already enabled
}

func (c *Client) DeleteMailDomainAntivirus(user, domain string) error {
	rc, err := c.do("v-delete-mail-domain-antivirus", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Remote DNS Host ───────────────────────────────────────────────────────────

type RemoteDNSHost struct {
	Host      string `json:"HOST"`
	Port      string `json:"PORT"`
	User      string `json:"USER"`
	DNSSystem string `json:"SYSTEM"`
	Type      string `json:"TYPE"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateRemoteDNSHost(host, port, user, password, dnsSystem, hostType string) error {
	rc, err := c.do("v-add-remote-dns-host", host, port, user, password, dnsSystem, hostType)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ListRemoteDNSHosts() (map[string]RemoteDNSHost, error) {
	var out map[string]RemoteDNSHost
	if err := c.doJSON("v-list-remote-dns-hosts", &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ReadRemoteDNSHost(host string) (*RemoteDNSHost, error) {
	hosts, err := c.ListRemoteDNSHosts()
	if err != nil {
		return nil, err
	}
	h, ok := hosts[host]
	if !ok {
		return nil, nil
	}
	return &h, nil
}

func (c *Client) DeleteRemoteDNSHost(host string) error {
	rc, err := c.do("v-delete-remote-dns-host", host)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}

// ── Backup Host ───────────────────────────────────────────────────────────────

type BackupHost struct {
	Host      string `json:"HOST"`
	Username  string `json:"USERNAME"`
	Path      string `json:"PATH"`
	Port      string `json:"PORT"`
	Type      string `json:"TYPE"`
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateBackupHost(hostType, host, username, password, path, port string) error {
	rc, err := c.do("v-add-backup-host", hostType, host, username, password, path, port)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ReadBackupHost(hostType, host string) (*BackupHost, error) {
	var out map[string]BackupHost
	if err := c.doJSON("v-list-backup-host", &out, hostType, host); err != nil {
		return nil, err
	}
	h, ok := out[host]
	if !ok {
		return nil, nil
	}
	return &h, nil
}

func (c *Client) DeleteBackupHost(hostType, host string) error {
	rc, err := c.do("v-delete-backup-host", hostType, host)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
}
