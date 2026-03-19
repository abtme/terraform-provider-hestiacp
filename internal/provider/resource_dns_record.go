package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-hestiacp/internal/client"
)

var _ resource.Resource = &DNSRecordResource{}

func NewDNSRecordResource() resource.Resource { return &DNSRecordResource{} }

type DNSRecordResource struct{ client *client.Client }

type DNSRecordResourceModel struct {
	ID       types.String `tfsdk:"id"`
	User     types.String `tfsdk:"user"`
	Domain   types.String `tfsdk:"domain"`
	Record   types.String `tfsdk:"record"`
	Type     types.String `tfsdk:"type"`
	Value    types.String `tfsdk:"value"`
	Priority types.Int64  `tfsdk:"priority"`
}

func (r *DNSRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *DNSRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a single DNS record within a HestiaCP DNS zone.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":   schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain": schema.StringAttribute{Required: true, MarkdownDescription: "DNS zone the record belongs to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"record": schema.StringAttribute{Required: true, MarkdownDescription: "Record name, e.g. `www`, `@`, `mail`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"type":   schema.StringAttribute{Required: true, MarkdownDescription: "Record type: A, AAAA, CNAME, MX, TXT, NS, SRV, etc.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"value":  schema.StringAttribute{Required: true, MarkdownDescription: "Record value / IP address."},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(10),
				MarkdownDescription: "Record priority (used for MX/SRV records). Defaults to `10`.",
			},
		},
	}
}

func (r *DNSRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateDNSRecord(
		plan.User.ValueString(),
		plan.Domain.ValueString(),
		plan.Record.ValueString(),
		plan.Type.ValueString(),
		plan.Value.ValueString(),
		int(plan.Priority.ValueInt64()),
	); err != nil {
		resp.Diagnostics.AddError("Error creating DNS record", err.Error())
		return
	}

	// Read back the assigned numeric ID from HestiaCP for future deletes.
	id, err := r.client.FindDNSRecordID(
		plan.User.ValueString(), plan.Domain.ValueString(),
		plan.Record.ValueString(), plan.Type.ValueString(), plan.Value.ValueString(),
	)
	if err != nil || id == "" {
		id = fmt.Sprintf("%s/%s/%s/%s", plan.User.ValueString(), plan.Domain.ValueString(), plan.Record.ValueString(), plan.Type.ValueString())
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	records, err := r.client.ListDNSRecords(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS records", err.Error())
		return
	}

	rec, ok := records[state.ID.ValueString()]
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Record = types.StringValue(rec.Record)
	state.Type = types.StringValue(rec.Type)
	state.Value = types.StringValue(rec.Value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// record/type/domain/user all have RequiresReplace; only value & priority
	// can change in-place. Simplest approach: delete + recreate.
	var plan, state DNSRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_ = r.client.DeleteDNSRecord(state.User.ValueString(), state.Domain.ValueString(), state.ID.ValueString())

	if err := r.client.CreateDNSRecord(
		plan.User.ValueString(), plan.Domain.ValueString(),
		plan.Record.ValueString(), plan.Type.ValueString(),
		plan.Value.ValueString(), int(plan.Priority.ValueInt64()),
	); err != nil {
		resp.Diagnostics.AddError("Error updating DNS record", err.Error())
		return
	}

	id, _ := r.client.FindDNSRecordID(
		plan.User.ValueString(), plan.Domain.ValueString(),
		plan.Record.ValueString(), plan.Type.ValueString(), plan.Value.ValueString(),
	)
	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDNSRecord(state.User.ValueString(), state.Domain.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting DNS record", err.Error())
	}
}
