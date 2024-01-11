package provider

import (
	"context"
	"encoding/json"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &CarrierDiagResource{}
	_ resource.ResourceWithConfigure   = &CarrierDiagResource{}
	_ resource.ResourceWithImportState = &CarrierDiagResource{}
)

// NewCarrierDiagResource is a helper function to simplify the provider implementation.
func NewCarrierDiagResource() resource.Resource {
	return &CarrierDiagResource{}
}

type CarrierDiagResource struct {
	client *xrcm_pf.Client
}

type CarrierDiagResourceData struct {
	Id             types.String `tfsdk:"id"`
	N              types.String `tfsdk:"n"`
	DeviceId       types.String `tfsdk:"deviceid"`
	LinePTPId      types.String `tfsdk:"lineptpid"`
	CarrierId      types.String `tfsdk:"carrierid"`
	TermLB         types.String `tfsdk:"termlb"`
	TermLBDuration types.Int64  `tfsdk:"termlbduration"`
}

// Metadata returns the data source type name.
func (r *CarrierDiagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_carrier_diag"
}

// Schema defines the schema for the data source.
func (r *CarrierDiagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Carrier Diagnostic",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Carrier Resource ID.",
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
			"lineptpid": schema.StringAttribute{
				Description: "line ptp id",
				Optional:    true,
			},
			"carrierid": schema.StringAttribute{
				Description: "carrier id",
				Optional:    true,
			},
			"termlb": schema.StringAttribute{
				Description: "Term Loopback",
				Optional:    true,
			},
			"termlbduration": schema.Int64Attribute{
				Description: "Term Loopback Duration",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CarrierDiagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*xrcm_pf.Client)
}

// Create creates the resource and sets the initial Terraform state.

func (r CarrierDiagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CarrierDiagResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "CarrierDiagResource: Create", map[string]interface{}{"CarrierDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r CarrierDiagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CarrierDiagResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "CarrierDiagResource: Read", map[string]interface{}{"CarrierDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.read(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id == types.StringValue("") {
		resp.State = tfsdk.State{}
	} else {
		diags = resp.State.Set(ctx, &data)
	}

	resp.Diagnostics.Append(diags...)
}

func (r CarrierDiagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CarrierDiagResourceData

	diags := req.Plan.Get(ctx, &data)

	tflog.Debug(ctx, "CarrierDiagResource: Update", map[string]interface{}{"CarrierDiagResourceData": data})
	// diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r CarrierDiagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CarrierDiagResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "CarrierDiagResource: Delete - ", map[string]interface{}{"CarrierDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CarrierDiagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CarrierDiagResource) update(plan *CarrierDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {
	if plan.LinePTPId.IsNull() || plan.CarrierId.IsNull() {
		diags.AddError(
			"CarrierDiagResource: update ##: Error Update Carrier",
			"Update: Could not Update Carrier, Port ID, Carrier ID  are not specified",
		)
		return
	}

	tflog.Debug(ctx, "CarrierDiagResource: update ## ", map[string]interface{}{"LinePTPId": plan.LinePTPId.ValueString(), "Carrier": plan.CarrierId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.TermLB.IsNull()) {
		cmd["termLB"] = plan.TermLB.ValueString()
	}

	if !(plan.TermLBDuration.IsNull()) {
		cmd["termLBDuration"] = plan.TermLBDuration.ValueInt64()
	}

	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"CarrierDiagResource: update ##: Error creating Carrier Diagnostic",
			"Update: not create Carrier Diagnostic, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + plan.LinePTPId.ValueString() + "/carriers/" + plan.CarrierId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "CarrierDiagResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources/lineptps/"+plan.LinePTPId.ValueString()+"/carriers/"+plan.CarrierId.ValueString()+"/diagnostic", rb)

	if err != nil {
		diags.AddError(
			"CarrierDiagResource: update ##: Error Update Carrier Diagnostic",
			"Update:Could not Update, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "CarrierDiagResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	_, err = SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"CarrierDiagResource: update ##: Error Update Carrier Diagnostic",
			"Update: Could not SetData Carrier, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *CarrierDiagResource) read(plan *CarrierDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.LinePTPId.IsNull() || plan.CarrierId.IsNull() {
		diags.AddError(
			"CarrierDiagResource: read ##: Error Read Carrier Diagnostic",
			"Read: Could not Read  Carrier, Port ID, Carrier ID are not specified",
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + plan.LinePTPId.ValueString() + "/carriers/" + plan.CarrierId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "CarrierDiagResource: read ## ", map[string]interface{}{"device": plan.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"CarrierDiagResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "CarrierDiagResource: read - not found ## 404", map[string]interface{}{"plan": plan})
		return
	}

	tflog.Debug(ctx, "CarrierDiagResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"CarrierDiagResource: read ##: Error SetData Carrier Diagnostic",
			"Read: Could not SetData Carrier Diagnostic, unexpected error: "+err.Error(),
		)
		return
	}

	if !(plan.TermLB.IsNull()) {
		plan.TermLB = types.StringValue(content["termLB"].(string))
	}

	if !(plan.TermLBDuration.IsNull()) {
		plan.TermLBDuration = types.Int64Value(int64(content["termLBDuration"].(float64)))
	}

	tflog.Debug(ctx, "CarrierDiagResource: read ## ", map[string]interface{}{"plan": plan})
}
