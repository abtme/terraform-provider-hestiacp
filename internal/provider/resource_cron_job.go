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

var _ resource.Resource = &CronJobResource{}

func NewCronJobResource() resource.Resource { return &CronJobResource{} }

type CronJobResource struct{ client *client.Client }

type CronJobResourceModel struct {
	ID      types.String `tfsdk:"id"`
	User    types.String `tfsdk:"user"`
	Min     types.String `tfsdk:"min"`
	Hour    types.String `tfsdk:"hour"`
	Day     types.String `tfsdk:"day"`
	Month   types.String `tfsdk:"month"`
	Wday    types.String `tfsdk:"wday"`
	Command types.String `tfsdk:"command"`
}

func (r *CronJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cron_job"
}

func (r *CronJobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a cron job for a HestiaCP user.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":    schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"min":     schema.StringAttribute{Required: true, MarkdownDescription: "Minute field (0-59 or `*`)."},
			"hour":    schema.StringAttribute{Required: true, MarkdownDescription: "Hour field (0-23 or `*`)."},
			"day":     schema.StringAttribute{Required: true, MarkdownDescription: "Day of month field (1-31 or `*`)."},
			"month":   schema.StringAttribute{Required: true, MarkdownDescription: "Month field (1-12 or `*`)."},
			"wday":    schema.StringAttribute{Required: true, MarkdownDescription: "Day of week field (0-6, Sunday=0, or `*`)."},
			"command": schema.StringAttribute{Required: true, MarkdownDescription: "Shell command to execute.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *CronJobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CronJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CronJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateCronJob(
		user,
		plan.Min.ValueString(),
		plan.Hour.ValueString(),
		plan.Day.ValueString(),
		plan.Month.ValueString(),
		plan.Wday.ValueString(),
		plan.Command.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error creating cron job", err.Error())
		return
	}

	id, err := r.client.FindCronJobID(user,
		plan.Min.ValueString(), plan.Hour.ValueString(),
		plan.Day.ValueString(), plan.Month.ValueString(),
		plan.Wday.ValueString(), plan.Command.ValueString(),
	)
	if err != nil || id == "" {
		id = fmt.Sprintf("%s/%s", user, plan.Command.ValueString())
	}
	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CronJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CronJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobs, err := r.client.ListCronJobs(state.User.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading cron jobs", err.Error())
		return
	}

	job, ok := jobs[state.ID.ValueString()]
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Min = types.StringValue(job.Min)
	state.Hour = types.StringValue(job.Hour)
	state.Day = types.StringValue(job.Day)
	state.Month = types.StringValue(job.Month)
	state.Wday = types.StringValue(job.Wday)
	state.Command = types.StringValue(job.Command)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CronJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CronJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateCronJob(
		state.User.ValueString(),
		state.ID.ValueString(),
		plan.Min.ValueString(),
		plan.Hour.ValueString(),
		plan.Day.ValueString(),
		plan.Month.ValueString(),
		plan.Wday.ValueString(),
		plan.Command.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error updating cron job", err.Error())
		return
	}

	plan.ID = state.ID
	plan.User = state.User
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CronJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CronJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteCronJob(state.User.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting cron job", err.Error())
	}
}
