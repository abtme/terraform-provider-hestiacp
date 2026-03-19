package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
	baseURL    string // e.g. https://myserver.com:8083
	accessKey  string // ACCESS_KEY:SECRET_KEY
	httpClient *http.Client
}

// New creates a new HestiaCP client.
func New(baseURL, accessKey string) *Client {
	return &Client{
		baseURL:   strings.TrimSuffix(baseURL, "/"),
		accessKey: accessKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// apiURL returns the full API endpoint.
func (c *Client) apiURL() string {
	return c.baseURL + "/api/"
}

// do executes a POST to the HestiaCP API with form-encoded body.
// cmd  = v-add-user, v-list-web-domains, etc.
// args = positional arg1..argN values.
// Returns the raw response body.
func (c *Client) do(cmd string, args ...string) (string, error) {
	form := url.Values{}
	form.Set("hash", c.accessKey)
	form.Set("cmd", cmd)
	form.Set("returncode", "yes")

	for i, a := range args {
		form.Set(fmt.Sprintf("arg%d", i+1), a)
	}

	resp, err := c.httpClient.Post(c.apiURL(), "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("HestiaCP API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HestiaCP response: %w", err)
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
		return fmt.Errorf("failed to parse HestiaCP JSON response (%s): %w", string(body), err)
	}
	return nil
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
	Email    string `json:"EMAIL"`
	Package  string `json:"PACKAGE"`
	Name     string `json:"NAME"`
	Shell    string `json:"SHELL"`
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
	return checkRC(rc)
}

func (c *Client) ReadWebDomain(user, domain string) (*WebDomain, error) {
	var out map[string]WebDomain
	if err := c.doJSON("v-list-web-domain", &out, user, domain); err != nil {
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
	Suspended string `json:"SUSPENDED"`
}

func (c *Client) CreateMailDomain(user, domain string) error {
	rc, err := c.do("v-add-mail-domain", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) ReadMailDomain(user, domain string) (*MailDomain, error) {
	var out map[string]MailDomain
	if err := c.doJSON("v-list-mail-domain", &out, user, domain); err != nil {
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
	return checkRC(rc)
}

func (c *Client) ReadMailAccount(user, domain, account string) (*MailAccount, error) {
	var out map[string]MailAccount
	if err := c.doJSON("v-list-mail-account", &out, user, domain, account); err != nil {
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
	// uses Let's Encrypt via v-add-letsencrypt-domain
	// aliases = comma-separated list e.g. "www.example.com,mail.example.com"
	rc, err := c.do("v-add-letsencrypt-domain", user, domain, aliases, "yes")
	if err != nil {
		return err
	}
	return checkRC(rc)
}

func (c *Client) DeleteSSL(user, domain string) error {
	rc, err := c.do("v-delete-ssl", user, domain)
	if err != nil {
		return err
	}
	return checkRC(rc, 3)
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
