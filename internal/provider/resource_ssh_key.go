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

var _ resource.Resource = &SSHKeyResource{}

func NewSSHKeyResource() resource.Resource { return &SSHKeyResource{} }

type SSHKeyResource struct{ client *client.Client }

type SSHKeyResourceModel struct {
	ID   types.String `tfsdk:"id"`
	User types.String `tfsdk:"user"`
	Key  types.String `tfsdk:"key"`
}

func (r *SSHKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *SSHKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an SSH public key for a HestiaCP user. The key comment (third field) is used as the key identifier in HestiaCP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Key identifier — the comment field of the public key.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"user": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "SSH public key in `TYPE BASE64 COMMENT` format, e.g. `ssh-ed25519 AAAA... my-key`. The comment is required and must be unique per user.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *SSHKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSHKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSHKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	keyStr := plan.Key.ValueString()
	comment := client.SSHKeyComment(keyStr)
	if comment == "" {
		resp.Diagnostics.AddError("Invalid SSH key", "Key must be in 'TYPE BASE64 COMMENT' format — the comment field is required and used as the key identifier.")
		return
	}

	if err := r.client.CreateSSHKey(user, keyStr); err != nil {
		resp.Diagnostics.AddError("Error adding SSH key", err.Error())
		return
	}

	plan.ID = types.StringValue(comment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSHKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SSHKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keys, err := r.client.ListSSHKeys(state.User.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSH keys", err.Error())
		return
	}

	// HestiaCP uses the key comment as the map key.
	if _, ok := keys[state.ID.ValueString()]; !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	// Key content is not returned by HestiaCP (only fingerprint) — preserve state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SSHKeyResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// All fields have RequiresReplace — Update is never called.
}

func (r *SSHKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SSHKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteSSHKey(state.User.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SSH key", err.Error())
	}
}
