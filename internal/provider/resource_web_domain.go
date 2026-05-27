package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-hestiacp/internal/client"
)

var _ resource.Resource = &WebDomainResource{}
var _ resource.ResourceWithImportState = &WebDomainResource{}

func NewWebDomainResource() resource.Resource { return &WebDomainResource{} }

type WebDomainResource struct{ client *client.Client }

type WebDomainResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
	IP     types.String `tfsdk:"ip"`
}

func (r *WebDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_domain"
}

func (r *WebDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a web domain (virtual host) in HestiaCP.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":   schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username that owns this domain. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain": schema.StringAttribute{Required: true, MarkdownDescription: "Domain name, e.g. `example.com`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"ip":     schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Server IP address to bind the domain to. Defaults to `0.0.0.0`."},
		},
	}
}

func (r *WebDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	ip := plan.IP.ValueString()
	if ip == "" {
		ip = "0.0.0.0"
	}

	if err := r.client.CreateWebDomain(user, plan.Domain.ValueString(), ip); err != nil {
		resp.Diagnostics.AddError("Error creating web domain", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", user, plan.Domain.ValueString()))
	plan.IP = types.StringValue(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d, err := r.client.ReadWebDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading web domain", err.Error())
		return
	}
	if d == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.IP = types.StringValue(d.IP)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update — IP changes require replace (handled by RequiresReplace on domain).
// No mutable fields beyond what triggers replace, so Update is a no-op.
func (r *WebDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WebDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// ImportState accepts "user/domain".
func (r *WebDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: user/domain")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[1])...)
}

func (r *WebDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteWebDomain(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting web domain", err.Error())
	}
}
