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

var _ resource.Resource = &FirewallRuleResource{}

func NewFirewallRuleResource() resource.Resource { return &FirewallRuleResource{} }

type FirewallRuleResource struct{ client *client.Client }

type FirewallRuleResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Action   types.String `tfsdk:"action"`
	Protocol types.String `tfsdk:"protocol"`
	Port     types.String `tfsdk:"port"`
	IP       types.String `tfsdk:"ip"`
	Comment  types.String `tfsdk:"comment"`
	Rule     types.String `tfsdk:"rule"`
}

func (r *FirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (r *FirewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HestiaCP firewall rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Rule action: `ACCEPT` or `DROP`.",
			},
			"protocol": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("TCP"),
				MarkdownDescription: "Protocol: `TCP`, `UDP`, or `ICMP`. Defaults to `TCP`.",
			},
			"port": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Port number or range, e.g. `80`, `443`, `8080:8090`. Use `0` for ICMP.",
			},
			"ip": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Source IP or CIDR to match. Empty string means any source.",
			},
			"comment": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Optional label for the rule (no spaces).",
			},
			"rule": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Optional rule position number. Lower numbers are evaluated first.",
			},
		},
	}
}

func (r *FirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateFirewallRule(
		plan.Action.ValueString(),
		plan.Protocol.ValueString(),
		plan.Port.ValueString(),
		plan.IP.ValueString(),
		plan.Comment.ValueString(),
		plan.Rule.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error creating firewall rule", err.Error())
		return
	}

	id, err := r.client.FindFirewallRuleID(
		plan.Action.ValueString(),
		plan.Protocol.ValueString(),
		plan.Port.ValueString(),
		plan.IP.ValueString(),
	)
	if err != nil || id == "" {
		id = fmt.Sprintf("%s/%s/%s", plan.Action.ValueString(), plan.Protocol.ValueString(), plan.Port.ValueString())
	}
	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := r.client.ListFirewallRules()
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall rules", err.Error())
		return
	}

	rule, ok := rules[state.ID.ValueString()]
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Action = types.StringValue(rule.Action)
	state.Protocol = types.StringValue(rule.Protocol)
	state.Port = types.StringValue(rule.Port)
	state.IP = types.StringValue(rule.IP)
	state.Comment = types.StringValue(rule.Comment)
	state.Rule = types.StringValue(rule.Rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FirewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateFirewallRule(
		state.ID.ValueString(),
		plan.Action.ValueString(),
		plan.Protocol.ValueString(),
		plan.Port.ValueString(),
		plan.IP.ValueString(),
		plan.Comment.ValueString(),
		plan.Rule.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error updating firewall rule", err.Error())
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteFirewallRule(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall rule", err.Error())
	}
}
