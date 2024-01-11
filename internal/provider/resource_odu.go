package provider

import (
	"context"
	//"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	//"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ODUResource{}
	_ resource.ResourceWithConfigure   = &ODUResource{}
	_ resource.ResourceWithImportState = &ODUResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewODUResource() resource.Resource {
	return &ODUResource{}
}

type ODUResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *ODUResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_odu"
}

type ODUResourceData struct {
	Id       types.String `tfsdk:"id"`
	N        types.String `tfsdk:"n"`
	DeviceId types.String `tfsdk:"deviceid"`
	Aid      types.String `tfsdk:"aid"`
	OtuId    types.String `tfsdk:"otuid"`
	OduId    types.String `tfsdk:"oduid"`
	OduType  types.String `tfsdk:"odutype"`
}

// Schema defines the schema for the resource.
func (r *ODUResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Carrier",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the ODU.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"n": schema.StringAttribute{
				Description: "XR Device Name",
				Required:    true,
			},
			"deviceid": schema.StringAttribute{
				Description: "device id",
				Computed:    true,
			},
			"otuid": schema.StringAttribute{
				Description: "otu id",
				Optional:    true,
			},
			"oduid": schema.StringAttribute{
				Description: "odu id",
				Optional:    true,
			},
			"aid": schema.StringAttribute{
				Description: "aid",
				Computed:    true,
			},
			"odutype": schema.StringAttribute{
				Description: "odu type",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ODUResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r ODUResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ODUResourceData

	diags := req.Config.Get(ctx, &data)
	tflog.Debug(ctx, "ODUResource: Create", map[string]interface{}{"ODUResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//r.update(&data, ctx, &resp.Diagnostics)
	r.read(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ODUResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ODUResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "ODUResource: Read", map[string]interface{}{"ODUResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.read(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &data)

	resp.Diagnostics.Append(diags...)
}

func (r ODUResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ODUResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "ODUResource: Update", map[string]interface{}{"ODUResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ODUResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ODUResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "ODUResource: Delete", map[string]interface{}{"ODUResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ODUResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ODUResource) update(plan *ODUResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.OtuId.IsNull() || plan.OduId.IsNull() {
		diags.AddError(
			"ODUResource: update ##: Error Update ODU",
			"Create: Could not Create ODU, OTU ID, ODU ID  are not specified",
		)
		return
	}

	tflog.Debug(ctx, "ODUResource: update Nothing to configure ## ", map[string]interface{}{"OtuId": plan.OtuId.ValueString(), "OduId": plan.OduId.ValueString()})

}

func (r *ODUResource) read(plan *ODUResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.OtuId.IsNull() || plan.OduId.IsNull() {
		diags.AddError(
			"ODUResource: read ##: Error Read ODU",
			"Read: Could not Read ODU, OTU ID and ODU ID are not specified",
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/otus/" + plan.OtuId.ValueString() + "/odus/" + plan.OduId.ValueString()
	}

	tflog.Debug(ctx, "ODUResource: read ## ", map[string]interface{}{"device": plan.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		/*if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"ODUResource: read ##: Error Get ODU",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "ODUResource: read - not found ## 404", map[string]interface{}{"href": href})
		return*/
		diags.AddError(
			"ODUResource: read ##: Error Get ODU",
			"Read: Could not Get , unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "ODUResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"ODUResource: read ##: Error Read ODU",
			"Read: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}
	if content["oduType"] != nil {
		plan.OduType = types.StringValue(content["oduType"].(string))
	}

	tflog.Debug(ctx, "ODUResource: read ## ", map[string]interface{}{"plan": plan})
}
