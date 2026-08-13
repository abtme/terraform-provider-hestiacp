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

var _ resource.Resource = &DNSZoneResource{}

func NewDNSZoneResource() resource.Resource { return &DNSZoneResource{} }

type DNSZoneResource struct{ client *client.Client }

type DNSZoneResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
	IP     types.String `tfsdk:"ip"`
}

func (r *DNSZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (r *DNSZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a DNS zone in HestiaCP.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":   schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username that owns this zone. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"ip":     schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Default IP for zone SOA record."},
		},
	}
}

func (r *DNSZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSZoneResourceModel
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

	if err := r.client.CreateDNSZone(user, plan.Domain.ValueString(), ip); err != nil {
		resp.Diagnostics.AddError("Error creating DNS zone", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", user, plan.Domain.ValueString()))
	plan.IP = types.StringValue(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSZoneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	z, err := r.client.ReadDNSZone(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS zone", err.Error())
		return
	}
	if z == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.IP = types.StringValue(z.IP)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSZoneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSZoneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDNSZone(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting DNS zone", err.Error())
	}
}
