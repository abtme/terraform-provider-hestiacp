package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/abtme/terraform-provider-hestiacp/internal/client"
)

// ═══════════════════════════════════════════════════════════════════════════
// WEB DOMAIN ALIAS
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &WebDomainAliasResource{}

func NewWebDomainAliasResource() resource.Resource { return &WebDomainAliasResource{} }

type WebDomainAliasResource struct{ client *client.Client }

type WebDomainAliasResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
	Alias  types.String `tfsdk:"alias"`
}

func (r *WebDomainAliasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_domain_alias"
}

func (r *WebDomainAliasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a web domain alias (parked domain) in HestiaCP.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":   schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain": schema.StringAttribute{Required: true, MarkdownDescription: "Primary web domain the alias is attached to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"alias":  schema.StringAttribute{Required: true, MarkdownDescription: "Alias domain name, e.g. `www.example.com`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *WebDomainAliasResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *WebDomainAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebDomainAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateWebDomainAlias(user, plan.Domain.ValueString(), plan.Alias.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating web domain alias", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", user, plan.Domain.ValueString(), plan.Alias.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebDomainAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebDomainAliasResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {}

func (r *WebDomainAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebDomainAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteWebDomainAlias(state.User.ValueString(), state.Domain.ValueString(), state.Alias.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting web domain alias", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// WEB DOMAIN FTP
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &WebDomainFTPResource{}

func NewWebDomainFTPResource() resource.Resource { return &WebDomainFTPResource{} }

type WebDomainFTPResource struct{ client *client.Client }

type WebDomainFTPResourceModel struct {
	ID          types.String `tfsdk:"id"`
	User        types.String `tfsdk:"user"`
	Domain      types.String `tfsdk:"domain"`
	FTPUser     types.String `tfsdk:"ftp_user"`
	FTPPassword types.String `tfsdk:"ftp_password"`
	FTPPath     types.String `tfsdk:"ftp_path"`
}

func (r *WebDomainFTPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_domain_ftp"
}

func (r *WebDomainFTPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an FTP account for a HestiaCP web domain.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":         schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain":       schema.StringAttribute{Required: true, MarkdownDescription: "Web domain to attach FTP to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"ftp_user":     schema.StringAttribute{Required: true, MarkdownDescription: "FTP username.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"ftp_password": schema.StringAttribute{Required: true, Sensitive: true, MarkdownDescription: "FTP password."},
			"ftp_path":     schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(""), MarkdownDescription: "FTP chroot path. Defaults to the domain's document root."},
		},
	}
}

func (r *WebDomainFTPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *WebDomainFTPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebDomainFTPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateWebDomainFTP(
		user, plan.Domain.ValueString(),
		plan.FTPUser.ValueString(), plan.FTPPassword.ValueString(),
		plan.FTPPath.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error creating web domain FTP", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", user, plan.Domain.ValueString(), plan.FTPUser.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainFTPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebDomainFTPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebDomainFTPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state WebDomainFTPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.FTPPassword.Equal(state.FTPPassword) {
		if err := r.client.UpdateWebDomainFTPPassword(
			state.User.ValueString(), state.Domain.ValueString(),
			state.FTPUser.ValueString(), plan.FTPPassword.ValueString(),
		); err != nil {
			resp.Diagnostics.AddError("Error updating FTP password", err.Error())
			return
		}
	}

	plan.ID = state.ID
	plan.User = state.User
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainFTPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebDomainFTPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteWebDomainFTP(state.User.ValueString(), state.Domain.ValueString(), state.FTPUser.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting web domain FTP", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// WEB DOMAIN REDIRECT
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &WebDomainRedirectResource{}

func NewWebDomainRedirectResource() resource.Resource { return &WebDomainRedirectResource{} }

type WebDomainRedirectResource struct{ client *client.Client }

type WebDomainRedirectResourceModel struct {
	ID       types.String `tfsdk:"id"`
	User     types.String `tfsdk:"user"`
	Domain   types.String `tfsdk:"domain"`
	Redirect types.String `tfsdk:"redirect"`
	HTTPCode types.String `tfsdk:"http_code"`
}

func (r *WebDomainRedirectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_domain_redirect"
}

func (r *WebDomainRedirectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an HTTP redirect for a HestiaCP web domain.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":      schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain":    schema.StringAttribute{Required: true, MarkdownDescription: "Web domain to redirect.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"redirect":  schema.StringAttribute{Required: true, MarkdownDescription: "Target URL to redirect to."},
			"http_code": schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("301"), MarkdownDescription: "HTTP redirect code. Defaults to `301`."},
		},
	}
}

func (r *WebDomainRedirectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *WebDomainRedirectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebDomainRedirectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateWebDomainRedirect(
		user, plan.Domain.ValueString(),
		plan.Redirect.ValueString(), plan.HTTPCode.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error creating web domain redirect", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", user, plan.Domain.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainRedirectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebDomainRedirectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebDomainRedirectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state WebDomainRedirectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateWebDomainRedirect(
		state.User.ValueString(), state.Domain.ValueString(),
		plan.Redirect.ValueString(), plan.HTTPCode.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error updating web domain redirect", err.Error())
		return
	}

	plan.ID = state.ID
	plan.User = state.User
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainRedirectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebDomainRedirectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteWebDomainRedirect(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting web domain redirect", err.Error())
	}
}
