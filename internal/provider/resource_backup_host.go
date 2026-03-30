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

var _ resource.Resource = &BackupHostResource{}

func NewBackupHostResource() resource.Resource { return &BackupHostResource{} }

type BackupHostResource struct{ client *client.Client }

type BackupHostResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Type     types.String `tfsdk:"type"`
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Path     types.String `tfsdk:"path"`
	Port     types.String `tfsdk:"port"`
}

func (r *BackupHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_host"
}

func (r *BackupHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configures a remote backup destination host in HestiaCP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Backup protocol: `ftp`, `sftp`, `b2`, `rclone`, etc.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"host": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Hostname or IP of the backup server.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Username for authenticating to the backup server.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for authenticating to the backup server.",
			},
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Remote path where backups are stored.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"port": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Port number of the backup server.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *BackupHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BackupHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BackupHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateBackupHost(
		plan.Type.ValueString(),
		plan.Host.ValueString(),
		plan.Username.ValueString(),
		plan.Password.ValueString(),
		plan.Path.ValueString(),
		plan.Port.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error creating backup host", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Type.ValueString() + ":" + plan.Host.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BackupHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BackupHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	h, err := r.client.ReadBackupHost(state.Type.ValueString(), state.Host.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading backup host", err.Error())
		return
	}
	if h == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Username = types.StringValue(h.Username)
	state.Path = types.StringValue(h.Path)
	state.Port = types.StringValue(h.Port)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BackupHostResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all significant fields are RequiresReplace — no in-place updates
}

func (r *BackupHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BackupHostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteBackupHost(state.Type.ValueString(), state.Host.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting backup host", err.Error())
	}
}
