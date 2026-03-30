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

var _ resource.Resource = &WebDomainHTTPAuthResource{}

func NewWebDomainHTTPAuthResource() resource.Resource { return &WebDomainHTTPAuthResource{} }

type WebDomainHTTPAuthResource struct{ client *client.Client }

type WebDomainHTTPAuthResourceModel struct {
	ID           types.String `tfsdk:"id"`
	User         types.String `tfsdk:"user"`
	Domain       types.String `tfsdk:"domain"`
	AuthUser     types.String `tfsdk:"auth_user"`
	AuthPassword types.String `tfsdk:"auth_password"`
}

func (r *WebDomainHTTPAuthResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_domain_httpauth"
}

func (r *WebDomainHTTPAuthResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Adds HTTP basic authentication to a HestiaCP web domain.",
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
				MarkdownDescription: "Web domain to protect.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"auth_user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HTTP basic auth username.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"auth_password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "HTTP basic auth password.",
			},
		},
	}
}

func (r *WebDomainHTTPAuthResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebDomainHTTPAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebDomainHTTPAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}

	if err := r.client.CreateWebDomainHTTPAuth(user, plan.Domain.ValueString(), plan.AuthUser.ValueString(), plan.AuthPassword.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating web domain HTTP auth", err.Error())
		return
	}

	plan.User = types.StringValue(user)
	plan.ID = types.StringValue(user + ":" + plan.Domain.ValueString() + ":" + plan.AuthUser.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainHTTPAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebDomainHTTPAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// No HestiaCP list command for httpauth — pass through existing state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebDomainHTTPAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state WebDomainHTTPAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := state.User.ValueString()
	if err := r.client.UpdateWebDomainHTTPAuth(user, state.Domain.ValueString(), state.AuthUser.ValueString(), plan.AuthPassword.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating web domain HTTP auth", err.Error())
		return
	}

	plan.User = state.User
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebDomainHTTPAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebDomainHTTPAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteWebDomainHTTPAuth(state.User.ValueString(), state.Domain.ValueString(), state.AuthUser.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting web domain HTTP auth", err.Error())
	}
}
