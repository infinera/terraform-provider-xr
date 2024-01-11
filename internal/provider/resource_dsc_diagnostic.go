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
	_ resource.Resource                = &DSCDiagResource{}
	_ resource.ResourceWithConfigure   = &DSCDiagResource{}
	_ resource.ResourceWithImportState = &DSCDiagResource{}
)

// NewCarrierDiagResource is a helper function to simplify the provider implementation.
func NewDSCDiagResource() resource.Resource {
	return &DSCDiagResource{}
}

type DSCDiagResource struct {
	client *xrcm_pf.Client
}

type DSCDiagResourceData struct {
	Id         types.String `tfsdk:"id"`
	DeviceId   types.String `tfsdk:"deviceid"`
	N          types.String `tfsdk:"n"`
	LinePTPId  types.String `tfsdk:"lineptpid"`
	CarrierId  types.String `tfsdk:"carrierid"`
	DscId      types.String `tfsdk:"dscid"`
	FacPRBSGen types.Bool   `tfsdk:"facprbsgen"`
	FacPRBSMon types.Bool   `tfsdk:"facprbsmon"`
}

// Metadata returns the data source type name.
func (r *DSCDiagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dsc_diag"
}

func (r *DSCDiagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an DSC Diagnostic",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "DSC Diagnostic ID.",
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
			"dscid": schema.StringAttribute{
				Description: "DSC id",
				Optional:    true,
			},
			"facprbsgen": schema.BoolAttribute{
				Description: "fac PRBS gen",
				Optional:    true,
			},
			"facprbsmon": schema.BoolAttribute{
				Description: "fac PRBS mon",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *DSCDiagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r DSCDiagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DSCDiagResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "DSCDiagResource: Create", map[string]interface{}{"DSCDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r DSCDiagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DSCDiagResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "DSCDiagResource: Read", map[string]interface{}{"DSCDiagResourceData": data})

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

func (r DSCDiagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DSCDiagResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "DSCDiagResource: Update", map[string]interface{}{"DSCDiagResourceData": data})
	// diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r DSCDiagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DSCDiagResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "DSCDiagResource: Delete", map[string]interface{}{"DSCDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *DSCDiagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DSCDiagResource) update(plan *DSCDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.LinePTPId.IsNull() || plan.CarrierId.IsNull() || plan.DscId.IsNull() {
		diags.AddError(
			"DSCDiagResource: update ##: Error Update DSC Diagnostic",
			"Update: Could not Update DSC, Port ID, Carrier ID, DSC ID  are not specified",
		)
		return
	}
	tflog.Debug(ctx, "DSCDiagResource: update ## ", map[string]interface{}{"LinePTPId": plan.LinePTPId.ValueString(), "Carrier": plan.CarrierId.ValueString(), "DscId": plan.DscId.ValueString()})

	var cmd = make(map[string]interface{})
	if !(plan.FacPRBSGen.IsNull()) {
		cmd["facPRBSGen"] = plan.FacPRBSGen.ValueBool()
	}

	if !(plan.FacPRBSMon.IsNull()) {
		cmd["facPRBSMon"] = plan.FacPRBSMon.ValueBool()
	}
	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"DSCDiagResource: update ##: Error Update DSC Diagnostic",
			"Update: Could not Update DSC, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + plan.LinePTPId.ValueString() + "/carriers/" + plan.CarrierId.ValueString() + "/dscs/" + plan.DscId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "DSCDiagResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources/lineptps/"+plan.LinePTPId.ValueString()+"/carriers/"+plan.CarrierId.ValueString()+"/dscs/"+plan.DscId.ValueString()+"/diagnostic", rb)

	if err != nil {
		diags.AddError(
			"DSCDiagResource: update ##: Error Update DSC Diagnostic",
			"Update: Could not Update DSC, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "DSCDiagResource: update ##  ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	_, err = SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"DSCDiagResource: update ##: Error Update DSC Diagnostic",
			"Update: Could not SetResourceId DSC, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "DSCDiagResource: update ## ", map[string]interface{}{"plan": plan})
}

func (r DSCDiagResource) read(state *DSCDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if state.LinePTPId.IsNull() || state.CarrierId.IsNull() || state.DscId.IsNull() {
		diags.AddError(
			"DSCDiagResource: read ##: DSCDiagResource: read ##: Error Read DSCG",
			"Read: Could not Read DSC, Port ID, Carrier ID, and DSC ID are not specified",
		)
		return
	}

	href := after(state.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + state.LinePTPId.ValueString() + "/carriers/" + state.CarrierId.ValueString() + "/dscs/" + state.DscId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "DSCDiagResource: read ## ", map[string]interface{}{"Device": state.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(state.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"DSCDiagResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		state.Id = types.StringValue("")
		tflog.Debug(ctx, "DSCDiagResource: read - not found ## 404", map[string]interface{}{"state": state})
		return
	}

	tflog.Debug(ctx, "DSCDiagResource: read ## ", map[string]interface{}{"response": string(body)})

	state.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(state.N.ValueString(), &state.Id, body)
	if err != nil {
		diags.AddError(
			"DSCDiagResource: read ##: Error Read DSC Diagnostic",
			"Read: Could not SetResourceId DSC, unexpected error: "+err.Error(),
		)
		return
	}

	if !(state.FacPRBSGen.IsNull()) {
		state.FacPRBSGen = types.BoolValue(content["facPRBSGen"].(bool))
	}
	if !(state.FacPRBSMon.IsNull()) {
		state.FacPRBSMon = types.BoolValue(content["facPRBSMon"].(bool))
	}

	tflog.Debug(ctx, "DSCDiagResource: read ## ", map[string]interface{}{"state": state})
}
