package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-hestiacp/internal/client"
)

// Ensure HestiacpProvider satisfies the provider.Provider interface.
var _ provider.Provider = &HestiacpProvider{}

// HestiacpProvider defines the provider implementation.
type HestiacpProvider struct {
	version string
}

// HestiacpProviderModel holds the provider configuration values.
type HestiacpProviderModel struct {
	URL       types.String `tfsdk:"url"`
	AccessKey types.String `tfsdk:"access_key"`
	Username  types.String `tfsdk:"username"`
}

// New returns a provider.Provider factory function.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HestiacpProvider{version: version}
	}
}

func (p *HestiacpProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hestiacp"
	resp.Version = p.version
}

func (p *HestiacpProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for HestiaCP — manage users, web domains, DNS, databases, email, SSL and backups via the HestiaCP REST API.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Base URL of the HestiaCP server including port, e.g. `https://myserver.com:8083`. May also be set via the `HESTIACP_URL` environment variable.",
				Optional:            true,
			},
			"access_key": schema.StringAttribute{
				MarkdownDescription: "HestiaCP access key in `ACCESS_KEY:SECRET_KEY` format. May also be set via the `HESTIACP_ACCESS_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Default HestiaCP username for resource operations. May also be set via the `HESTIACP_USERNAME` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *HestiacpProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config HestiacpProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve URL — config takes precedence over env var.
	url := os.Getenv("HESTIACP_URL")
	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing HestiaCP URL",
			"Set the url provider attribute or the HESTIACP_URL environment variable.",
		)
	}

	// Resolve access key.
	accessKey := os.Getenv("HESTIACP_ACCESS_KEY")
	if !config.AccessKey.IsNull() {
		accessKey = config.AccessKey.ValueString()
	}
	if accessKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Missing HestiaCP Access Key",
			"Set the access_key provider attribute or the HESTIACP_ACCESS_KEY environment variable.",
		)
	}

	// Resolve default username.
	username := os.Getenv("HESTIACP_USERNAME")
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c := client.New(url, accessKey, username)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *HestiacpProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewWebDomainResource,
		NewDNSZoneResource,
		NewDNSRecordResource,
		NewDatabaseResource,
		NewEmailDomainResource,
		NewEmailAccountResource,
		NewSSLResource,
		NewMailSSLResource,
		NewBackupResource,
		NewCronJobResource,
		NewFirewallRuleResource,
		NewWebDomainAliasResource,
		NewWebDomainFTPResource,
		NewWebDomainRedirectResource,
		NewSSHKeyResource,
		NewMailForwardResource,
		NewMailAliasResource,
		NewMailAutoreplyResource,
		NewWebDomainHTTPAuthResource,
		NewMailDomainCatchallResource,
		NewMailDomainDKIMResource,
		NewMailDomainAntispamResource,
		NewMailDomainAntivirusResource,
		NewRemoteDNSHostResource,
		NewBackupHostResource,
	}
}

func (p *HestiacpProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
