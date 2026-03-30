package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-hestiacp/internal/client"
)

// ── Mail Domain Catchall ──────────────────────────────────────────────────────

var _ resource.Resource = &MailDomainCatchallResource{}

func NewMailDomainCatchallResource() resource.Resource { return &MailDomainCatchallResource{} }

type MailDomainCatchallResource struct{ client *client.Client }

type MailDomainCatchallResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
	Email  types.String `tfsdk:"email"`
}

func (r *MailDomainCatchallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_domain_catchall"
}

func (r *MailDomainCatchallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sets a catch-all email address for a HestiaCP mail domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"user": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Mail domain to configure.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Catch-all destination email address.",
			},
		},
	}
}

func (r *MailDomainCatchallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailDomainCatchallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailDomainCatchallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}

	if err := r.client.CreateMailDomainCatchall(user, plan.Domain.ValueString(), plan.Email.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating mail domain catchall", err.Error())
		return
	}

	plan.User = types.StringValue(user)
	plan.ID = types.StringValue(user + ":" + plan.Domain.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailDomainCatchallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailDomainCatchallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	md, err := r.client.ReadMailDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail domain", err.Error())
		return
	}
	if md == nil || md.Catchall == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Email = types.StringValue(md.Catchall)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailDomainCatchallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MailDomainCatchallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateMailDomainCatchall(state.User.ValueString(), state.Domain.ValueString(), plan.Email.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating mail domain catchall", err.Error())
		return
	}

	plan.User = state.User
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailDomainCatchallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailDomainCatchallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMailDomainCatchall(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting mail domain catchall", err.Error())
	}
}

// ── Mail Domain DKIM ──────────────────────────────────────────────────────────

var _ resource.Resource = &MailDomainDKIMResource{}

func NewMailDomainDKIMResource() resource.Resource { return &MailDomainDKIMResource{} }

type MailDomainDKIMResource struct{ client *client.Client }

type MailDomainDKIMResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
}

func (r *MailDomainDKIMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_domain_dkim"
}

func (r *MailDomainDKIMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Enables DKIM signing for a HestiaCP mail domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"user": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Mail domain to enable DKIM on.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *MailDomainDKIMResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailDomainDKIMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailDomainDKIMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}

	if err := r.client.CreateMailDomainDKIM(user, plan.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error enabling mail domain DKIM", err.Error())
		return
	}

	plan.User = types.StringValue(user)
	plan.ID = types.StringValue(user + ":" + plan.Domain.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailDomainDKIMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailDomainDKIMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	md, err := r.client.ReadMailDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail domain", err.Error())
		return
	}
	if md == nil || md.DKIM != "yes" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailDomainDKIMResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// no updatable fields beyond RequiresReplace ones
}

func (r *MailDomainDKIMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailDomainDKIMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMailDomainDKIM(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error disabling mail domain DKIM", err.Error())
	}
}

// ── Mail Domain Antispam ──────────────────────────────────────────────────────

var _ resource.Resource = &MailDomainAntispamResource{}

func NewMailDomainAntispamResource() resource.Resource { return &MailDomainAntispamResource{} }

type MailDomainAntispamResource struct{ client *client.Client }

type MailDomainAntispamResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
}

func (r *MailDomainAntispamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_domain_antispam"
}

func (r *MailDomainAntispamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Enables antispam filtering for a HestiaCP mail domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"user": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Mail domain to enable antispam on.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *MailDomainAntispamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailDomainAntispamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailDomainAntispamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}

	if err := r.client.CreateMailDomainAntispam(user, plan.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error enabling mail domain antispam", err.Error())
		return
	}

	plan.User = types.StringValue(user)
	plan.ID = types.StringValue(user + ":" + plan.Domain.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailDomainAntispamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailDomainAntispamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	md, err := r.client.ReadMailDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail domain", err.Error())
		return
	}
	if md == nil || md.Antispam != "yes" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailDomainAntispamResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *MailDomainAntispamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailDomainAntispamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMailDomainAntispam(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error disabling mail domain antispam", err.Error())
	}
}

// ── Mail Domain Antivirus ─────────────────────────────────────────────────────

var _ resource.Resource = &MailDomainAntivirusResource{}

func NewMailDomainAntivirusResource() resource.Resource { return &MailDomainAntivirusResource{} }

type MailDomainAntivirusResource struct{ client *client.Client }

type MailDomainAntivirusResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
}

func (r *MailDomainAntivirusResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_domain_antivirus"
}

func (r *MailDomainAntivirusResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Enables antivirus scanning for a HestiaCP mail domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"user": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Mail domain to enable antivirus on.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *MailDomainAntivirusResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailDomainAntivirusResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailDomainAntivirusResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}

	if err := r.client.CreateMailDomainAntivirus(user, plan.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error enabling mail domain antivirus", err.Error())
		return
	}

	plan.User = types.StringValue(user)
	plan.ID = types.StringValue(user + ":" + plan.Domain.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailDomainAntivirusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailDomainAntivirusResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	md, err := r.client.ReadMailDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail domain", err.Error())
		return
	}
	if md == nil || md.Antivirus != "yes" {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailDomainAntivirusResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *MailDomainAntivirusResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailDomainAntivirusResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMailDomainAntivirus(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error disabling mail domain antivirus", err.Error())
	}
}
