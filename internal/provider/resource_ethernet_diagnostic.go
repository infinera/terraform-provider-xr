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
	_ resource.Resource                = &EthernetDiagResource{}
	_ resource.ResourceWithConfigure   = &EthernetDiagResource{}
	_ resource.ResourceWithImportState = &EthernetDiagResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewEthernetDiagResource() resource.Resource {
	return &EthernetDiagResource{}
}

type EthernetDiagResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *EthernetDiagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ethernet_diag"
}

type EthernetDiagResourceData struct {
	Id             types.String `tfsdk:"id"`
	N              types.String `tfsdk:"n"`
	DeviceId       types.String `tfsdk:"deviceid"`
	Aid            types.String `tfsdk:"aid"`
	EthernetId     types.String `tfsdk:"ethernetid"`
	TermLB         types.String `tfsdk:"termlb"`
	TermLBDuration types.Int64  `tfsdk:"termlbduration"`
	FacLB          types.String `tfsdk:"faclb"`
	FacLBDuration  types.Int64  `tfsdk:"faclbduration"`
}

// Schema defines the schema for the EthernetDiag resource.
func (r *EthernetDiagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an EthernetDiag",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the EthernetDiag.",
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
			"aid": schema.StringAttribute{
				Description: "AID",
				Computed:    true,
			},
			"ethernetid": schema.StringAttribute{
				Description: "ethernet id",
				Optional:    true,
			},
			"termlb": schema.StringAttribute{
				Description: "term Loopback",
				Optional:    true,
			},
			"termlbduration": schema.Int64Attribute{
				Description: "term Loopback Duration",
				Optional:    true,
			},
			"faclb": schema.StringAttribute{
				Description: "loopback type",
				Optional:    true,
			},
			"faclbduration": schema.Int64Attribute{
				Description: "term Loopback Duration",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *EthernetDiagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r EthernetDiagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EthernetDiagResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetDiagResource: Create", map[string]interface{}{"EthernetDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r EthernetDiagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EthernetDiagResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetDiagResource: Read", map[string]interface{}{"EthernetDiagResourceData": data})

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

func (r EthernetDiagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EthernetDiagResourceData

	diags := req.Plan.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetDiagResource: Update", map[string]interface{}{"EthernetDiagResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r EthernetDiagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EthernetDiagResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetDiagResource: Delete", map[string]interface{}{"EthernetDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *EthernetDiagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *EthernetDiagResource) update(plan *EthernetDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.EthernetId.IsNull() {
		diags.AddError(
			"EthernetDiagResource: update ##: Error Update Ethernet Diagnostic",
			"Update: Could not Update Ethernet Diag, Ethernet ID is not specified ",
		)
		return
	}

	tflog.Debug(ctx, "EthernetDiagResource: update ## ", map[string]interface{}{"EthernetId": plan.EthernetId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.TermLB.IsNull()) {
		cmd["termLB"] = plan.TermLB.ValueString()
	}

	if !(plan.TermLBDuration.IsNull()) {
		cmd["termLBDuration"] = plan.TermLBDuration.ValueInt64()
	}

	if !(plan.FacLB.IsNull()) {
		cmd["facLB"] = plan.FacLB.ValueString()
	}

	if !(plan.FacLBDuration.IsNull()) {
		cmd["facLBDuration"] = plan.FacLBDuration.ValueInt64()
	}

	if len(cmd) == 0. {
		tflog.Debug(ctx, "EthernetDiagResource: update ## No Settings, Nothing to Run Diagnostic", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources/ethernets/" + plan.EthernetId.ValueString() + "/diagnostic"})
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"EthernetDiagResource: update ##: Error Update Ethernet Diagnostic",
			"Update: Could not Update EthernetDiag, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + plan.EthernetId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "EthernetDiagResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"EthernetDiagResource: update ##: Error Update Ethernet Diagnostic",
			"Update: Could not Update EthernetDiag, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "EthernetDiagResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"EthernetDiagResource: update ##: Error Update Ethernet Diagnostic",
			"Update: Could not GetSetData EthernetDiag, unexpected error: "+err.Error(),
		)
		return
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}

	plan.DeviceId = types.StringValue(deviceId)

	tflog.Debug(ctx, "EthernetDiagResource: update ## ", map[string]interface{}{"plan": plan})

}
func (r EthernetDiagResource) read(state *EthernetDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if state.EthernetId.IsNull() {
		diags.AddError(
			"EthernetDiagResource: read ##: Error Read Ethernet Diagnostic",
			"Read: Could not Read Ethernet Diag, Ethernet ID is not specified ",
		)
		return
	}

	href := after(state.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + state.EthernetId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "EthernetDiagResource: read ## ", map[string]interface{}{"Device": state.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(state.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"EthernetDiagResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		state.Id = types.StringValue("")
		tflog.Debug(ctx, "EthernetDiagResource: read - not found ## 404", map[string]interface{}{"state": state})
		return
	}

	tflog.Debug(ctx, "EthernetDiagResource: read ## ", map[string]interface{}{"response": string(body)})
	state.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(state.N.ValueString(), &state.Id, body)

	if err != nil {
		diags.AddError(
			"EthernetDiagResource: read ##: Error Read Ethernet Diagnostic",
			"Read: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	aid := content["aid"]
	if aid != nil {
		state.Aid = types.StringValue(aid.(string))
	} 

	if !state.TermLB.IsNull() && content["termLB"] != nil {
		state.TermLB = types.StringValue(content["termLB"].(string))
	}

	if !state.TermLBDuration.IsNull() && content["termLBDuration"] != nil {
		state.TermLBDuration = types.Int64Value(int64(content["termLBDuration"].(float64)))
	}

	if !state.FacLB.IsNull() && content["facLB"] != nil {
		state.FacLB = types.StringValue(content["facLB"].(string))
	}

	if !state.FacLBDuration.IsNull() && content["facLBDuration"] != nil {
		state.FacLBDuration = types.Int64Value(int64(content["facLBDuration"].(float64)))
	}

	tflog.Debug(ctx, "EthernetDiagResource: read ## ", map[string]interface{}{"state": state})
}
