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

var _ resource.Resource = &UserResource{}

func NewUserResource() resource.Resource { return &UserResource{} }

type UserResource struct{ client *client.Client }

type UserResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	Package   types.String `tfsdk:"package"`
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HestiaCP user account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier (same as username).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HestiaCP username.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "User password.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User contact email address.",
			},
			"first_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "User's first name.",
			},
			"package": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "HestiaCP hosting package to assign. Defaults to `default`.",
			},
		},
	}
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pkg := plan.Package.ValueString()
	if pkg == "" {
		pkg = "default"
	}
	firstName := plan.FirstName.ValueString()

	if err := r.client.CreateUser(
		plan.Username.ValueString(),
		plan.Password.ValueString(),
		plan.Email.ValueString(),
		pkg,
		firstName,
	); err != nil {
		resp.Diagnostics.AddError("Error creating HestiaCP user", err.Error())
		return
	}

	plan.ID = plan.Username
	plan.Package = types.StringValue(pkg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := r.client.ReadUser(state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading HestiaCP user", err.Error())
		return
	}
	if u == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Email = types.StringValue(u.Email)
	state.Package = types.StringValue(u.Package)
	state.FirstName = types.StringValue(u.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := state.Username.ValueString()

	// Update password if changed.
	if !plan.Password.Equal(state.Password) {
		if err := r.client.UpdateUserPassword(username, plan.Password.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error updating user password", err.Error())
			return
		}
	}

	// Update email if changed.
	if !plan.Email.Equal(state.Email) {
		if err := r.client.UpdateUserEmail(username, plan.Email.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error updating user email", err.Error())
			return
		}
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteUser(state.Username.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting HestiaCP user", err.Error())
	}
}
