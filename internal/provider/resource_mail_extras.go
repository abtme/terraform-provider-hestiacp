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

// ═══════════════════════════════════════════════════════════════════════════
// MAIL FORWARD
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &MailForwardResource{}

func NewMailForwardResource() resource.Resource { return &MailForwardResource{} }

type MailForwardResource struct{ client *client.Client }

type MailForwardResourceModel struct {
	ID      types.String `tfsdk:"id"`
	User    types.String `tfsdk:"user"`
	Domain  types.String `tfsdk:"domain"`
	Account types.String `tfsdk:"account"`
	Forward types.String `tfsdk:"forward"`
}

func (r *MailForwardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_forward"
}

func (r *MailForwardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an email forward for a HestiaCP mail account.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":    schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain":  schema.StringAttribute{Required: true, MarkdownDescription: "Mail domain.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"account": schema.StringAttribute{Required: true, MarkdownDescription: "Mailbox name (part before the @).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"forward": schema.StringAttribute{Required: true, MarkdownDescription: "Email address to forward to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *MailForwardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailForwardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailForwardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateMailForward(user, plan.Domain.ValueString(), plan.Account.ValueString(), plan.Forward.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating mail forward", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s/%s", user, plan.Domain.ValueString(), plan.Account.ValueString(), plan.Forward.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailForwardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailForwardResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {}

func (r *MailForwardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMailForward(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString(), state.Forward.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting mail forward", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// MAIL ALIAS
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &MailAliasResource{}

func NewMailAliasResource() resource.Resource { return &MailAliasResource{} }

type MailAliasResource struct{ client *client.Client }

type MailAliasResourceModel struct {
	ID      types.String `tfsdk:"id"`
	User    types.String `tfsdk:"user"`
	Domain  types.String `tfsdk:"domain"`
	Account types.String `tfsdk:"account"`
	Alias   types.String `tfsdk:"alias"`
}

func (r *MailAliasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_alias"
}

func (r *MailAliasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an email alias for a HestiaCP mail account.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":    schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain":  schema.StringAttribute{Required: true, MarkdownDescription: "Mail domain.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"account": schema.StringAttribute{Required: true, MarkdownDescription: "Mailbox name the alias points to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"alias":   schema.StringAttribute{Required: true, MarkdownDescription: "Alias address (part before the @).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *MailAliasResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateMailAlias(user, plan.Domain.ValueString(), plan.Account.ValueString(), plan.Alias.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating mail alias", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s/%s", user, plan.Domain.ValueString(), plan.Account.ValueString(), plan.Alias.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailAliasResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {}

func (r *MailAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMailAlias(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString(), state.Alias.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting mail alias", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// MAIL AUTOREPLY
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &MailAutoreplyResource{}

func NewMailAutoreplyResource() resource.Resource { return &MailAutoreplyResource{} }

type MailAutoreplyResource struct{ client *client.Client }

type MailAutoreplyResourceModel struct {
	ID      types.String `tfsdk:"id"`
	User    types.String `tfsdk:"user"`
	Domain  types.String `tfsdk:"domain"`
	Account types.String `tfsdk:"account"`
	Message types.String `tfsdk:"message"`
}

func (r *MailAutoreplyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_autoreply"
}

func (r *MailAutoreplyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an autoreply message for a HestiaCP mail account.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":    schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain":  schema.StringAttribute{Required: true, MarkdownDescription: "Mail domain.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"account": schema.StringAttribute{Required: true, MarkdownDescription: "Mailbox name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"message": schema.StringAttribute{Required: true, MarkdownDescription: "Autoreply message body."},
		},
	}
}

func (r *MailAutoreplyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailAutoreplyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailAutoreplyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateMailAutoreply(user, plan.Domain.ValueString(), plan.Account.ValueString(), plan.Message.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating mail autoreply", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", user, plan.Domain.ValueString(), plan.Account.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailAutoreplyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailAutoreplyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.ReadMailAutoreply(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail autoreply", err.Error())
		return
	}
	if a == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Message = types.StringValue(a.Message)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailAutoreplyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MailAutoreplyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete and recreate to update message
	_ = r.client.DeleteMailAutoreply(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString())
	if err := r.client.CreateMailAutoreply(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString(), plan.Message.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating mail autoreply", err.Error())
		return
	}

	plan.ID = state.ID
	plan.User = state.User
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailAutoreplyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailAutoreplyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMailAutoreply(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting mail autoreply", err.Error())
	}
}
