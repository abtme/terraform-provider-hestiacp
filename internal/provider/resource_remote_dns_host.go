package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/abtme/terraform-provider-hestiacp/internal/client"
)

var _ resource.Resource = &RemoteDNSHostResource{}

func NewRemoteDNSHostResource() resource.Resource { return &RemoteDNSHostResource{} }

type RemoteDNSHostResource struct{ client *client.Client }

type RemoteDNSHostResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Host      types.String `tfsdk:"host"`
	Port      types.String `tfsdk:"port"`
	User      types.String `tfsdk:"user"`
	Password  types.String `tfsdk:"password"`
	DNSSystem types.String `tfsdk:"dns_system"`
	Type      types.String `tfsdk:"type"`
}

func (r *RemoteDNSHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_dns_host"
}

func (r *RemoteDNSHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Adds a remote DNS cluster host to HestiaCP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"host": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Hostname or IP of the remote DNS server.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"port": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HestiaCP API port on the remote server (e.g. `8083`).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Admin username on the remote HestiaCP server.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Admin password on the remote HestiaCP server.",
			},
			"dns_system": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "DNS system on the remote host (e.g. `bind9`). Defaults to `hestia`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Cluster role: `slave` or `master`. Defaults to `slave`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *RemoteDNSHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type", fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *RemoteDNSHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RemoteDNSHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsSystem := plan.DNSSystem.ValueString()
	if dnsSystem == "" {
		dnsSystem = "hestia"
	}
	hostType := plan.Type.ValueString()
	if hostType == "" {
		hostType = "slave"
	}

	if err := r.client.CreateRemoteDNSHost(
		plan.Host.ValueString(),
		plan.Port.ValueString(),
		plan.User.ValueString(),
		plan.Password.ValueString(),
		dnsSystem,
		hostType,
	); err != nil {
		resp.Diagnostics.AddError("Error creating remote DNS host", err.Error())
		return
	}

	plan.DNSSystem = types.StringValue(dnsSystem)
	plan.Type = types.StringValue(hostType)
	plan.ID = plan.Host
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RemoteDNSHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RemoteDNSHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	h, err := r.client.ReadRemoteDNSHost(state.Host.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading remote DNS host", err.Error())
		return
	}
	if h == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.DNSSystem = types.StringValue(h.DNSSystem)
	state.Type = types.StringValue(h.Type)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RemoteDNSHostResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all fields are RequiresReplace — no in-place updates
}

func (r *RemoteDNSHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RemoteDNSHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteRemoteDNSHost(state.Host.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting remote DNS host", err.Error())
	}
}
