package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/your-org/terraform-provider-hestiacp/internal/client"
)

// ═══════════════════════════════════════════════════════════════════════════
// DATABASE
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &DatabaseResource{}

func NewDatabaseResource() resource.Resource { return &DatabaseResource{} }

type DatabaseResource struct{ client *client.Client }

type DatabaseResourceModel struct {
	ID         types.String `tfsdk:"id"`
	User       types.String `tfsdk:"user"`
	DbName     types.String `tfsdk:"db_name"`
	DbUser     types.String `tfsdk:"db_user"`
	DbPassword types.String `tfsdk:"db_password"`
	DbType     types.String `tfsdk:"db_type"`
}

func (r *DatabaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HestiaCP database and database user.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":        schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username that owns this database. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"db_name":     schema.StringAttribute{Required: true, MarkdownDescription: "Database name (without user prefix).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"db_user":     schema.StringAttribute{Required: true, MarkdownDescription: "Database username (without user prefix).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"db_password": schema.StringAttribute{Required: true, Sensitive: true},
			"db_type":     schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("mysql"), MarkdownDescription: "Database type: `mysql` or `pgsql`. Defaults to `mysql`."},
		},
	}
}

func (r *DatabaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DatabaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateDatabase(user, plan.DbName.ValueString(), plan.DbUser.ValueString(), plan.DbPassword.ValueString(), plan.DbType.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating database", err.Error())
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", user, plan.DbName.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DatabaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	db, err := r.client.ReadDatabase(state.User.ValueString(), state.DbName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading database", err.Error())
		return
	}
	if db == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	state.DbType = types.StringValue(db.Type)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DatabaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DatabaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteDatabase(state.User.ValueString(), state.DbName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting database", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// EMAIL DOMAIN
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &EmailDomainResource{}
var _ resource.ResourceWithImportState = &EmailDomainResource{}

func NewEmailDomainResource() resource.Resource { return &EmailDomainResource{} }

type EmailDomainResource struct{ client *client.Client }

type EmailDomainResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
}

func (r *EmailDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_domain"
}

func (r *EmailDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a mail domain in HestiaCP.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":   schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username that owns this mail domain. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *EmailDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmailDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EmailDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateMailDomain(user, plan.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating mail domain", err.Error())
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", user, plan.Domain.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EmailDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EmailDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	md, err := r.client.ReadMailDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail domain", err.Error())
		return
	}
	if md == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EmailDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EmailDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EmailDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EmailDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMailDomain(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting mail domain", err.Error())
	}
}

// ImportState accepts "user/domain".
func (r *EmailDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: user/domain")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[1])...)
}

// ═══════════════════════════════════════════════════════════════════════════
// EMAIL ACCOUNT
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &EmailAccountResource{}
var _ resource.ResourceWithImportState = &EmailAccountResource{}

func NewEmailAccountResource() resource.Resource { return &EmailAccountResource{} }

type EmailAccountResource struct{ client *client.Client }

type EmailAccountResourceModel struct {
	ID       types.String `tfsdk:"id"`
	User     types.String `tfsdk:"user"`
	Domain   types.String `tfsdk:"domain"`
	Account  types.String `tfsdk:"account"`
	Password types.String `tfsdk:"password"`
	Quota    types.Int64  `tfsdk:"quota"`
}

func (r *EmailAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_account"
}

func (r *EmailAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a mail account (mailbox) within a HestiaCP mail domain.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":     schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username that owns this mail account. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain":   schema.StringAttribute{Required: true, MarkdownDescription: "Mail domain this account belongs to.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"account":  schema.StringAttribute{Required: true, MarkdownDescription: "Mailbox name (the part before the @).", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"password": schema.StringAttribute{Required: true, Sensitive: true},
			"quota": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "Mailbox quota in MB. `0` means unlimited.",
			},
		},
	}
}

func (r *EmailAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmailAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EmailAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateMailAccount(user, plan.Domain.ValueString(), plan.Account.ValueString(), plan.Password.ValueString(), int(plan.Quota.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error creating mail account", err.Error())
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s@%s", plan.Account.ValueString(), plan.Domain.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EmailAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EmailAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ma, err := r.client.ReadMailAccount(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail account", err.Error())
		return
	}
	if ma == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EmailAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EmailAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// ImportState accepts "user/domain/account".
func (r *EmailAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: user/domain/account")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), fmt.Sprintf("%s@%s", parts[2], parts[1]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account"), parts[2])...)
}

func (r *EmailAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EmailAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMailAccount(state.User.ValueString(), state.Domain.ValueString(), state.Account.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting mail account", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// SSL (Let's Encrypt)
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &SSLResource{}

func NewSSLResource() resource.Resource { return &SSLResource{} }

type SSLResource struct{ client *client.Client }

type SSLResourceModel struct {
	ID      types.String `tfsdk:"id"`
	User    types.String `tfsdk:"user"`
	Domain  types.String `tfsdk:"domain"`
	Aliases types.String `tfsdk:"aliases"`
}

func (r *SSLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssl"
}

func (r *SSLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Issues a Let's Encrypt SSL certificate for a HestiaCP web domain.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":    schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username that owns this domain. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain":  schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"aliases": schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(""), MarkdownDescription: "Comma-separated list of additional SANs, e.g. `www.example.com,mail.example.com`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *SSLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateSSL(user, plan.Domain.ValueString(), plan.Aliases.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error issuing SSL certificate", err.Error())
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", user, plan.Domain.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SSLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	d, err := r.client.ReadWebDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSL state", err.Error())
		return
	}
	if d == nil || d.SSL != "yes" {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SSLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SSLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SSLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteSSL(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SSL certificate", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// MAIL SSL (Let's Encrypt)
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &MailSSLResource{}

func NewMailSSLResource() resource.Resource { return &MailSSLResource{} }

type MailSSLResource struct{ client *client.Client }

type MailSSLResourceModel struct {
	ID     types.String `tfsdk:"id"`
	User   types.String `tfsdk:"user"`
	Domain types.String `tfsdk:"domain"`
}

func (r *MailSSLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail_ssl"
}

func (r *MailSSLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Issues a Let's Encrypt SSL certificate for a HestiaCP mail domain via v-add-letsencrypt-mail-ssl.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user":   schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"domain": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *MailSSLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MailSSLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MailSSLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateMailSSL(user, plan.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error issuing mail SSL certificate", err.Error())
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", user, plan.Domain.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailSSLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MailSSLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	d, err := r.client.ReadMailDomain(state.User.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading mail SSL state", err.Error())
		return
	}
	if d == nil || d.SSL != "yes" {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailSSLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MailSSLResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailSSLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MailSSLResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteMailSSL(state.User.ValueString(), state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting mail SSL certificate", err.Error())
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// BACKUP
// ═══════════════════════════════════════════════════════════════════════════

var _ resource.Resource = &BackupResource{}

func NewBackupResource() resource.Resource { return &BackupResource{} }

type BackupResource struct{ client *client.Client }

type BackupResourceModel struct {
	ID   types.String `tfsdk:"id"`
	User types.String `tfsdk:"user"`
}

func (r *BackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

func (r *BackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Triggers a HestiaCP user backup. Each apply creates a new backup snapshot.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"user": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "HestiaCP username to back up. Defaults to the provider `username`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		},
	}
}

func (r *BackupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user := plan.User.ValueString()
	if user == "" {
		user = r.client.DefaultUser()
	}
	plan.User = types.StringValue(user)

	if err := r.client.CreateBackup(user); err != nil {
		resp.Diagnostics.AddError("Error creating backup", err.Error())
		return
	}
	plan.ID = types.StringValue(user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BackupResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Backups are not deleted by Terraform — they are immutable snapshots.
}
