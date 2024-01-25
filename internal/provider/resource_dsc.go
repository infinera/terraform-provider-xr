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
	_ resource.Resource                = &DSCResource{}
	_ resource.ResourceWithConfigure   = &DSCResource{}
	_ resource.ResourceWithImportState = &DSCResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewDSCResource() resource.Resource {
	return &DSCResource{}
}

type DSCResource struct {
	client *xrcm_pf.Client
}

type DSCResourceData struct {
	Id          types.String `tfsdk:"id"`
	DeviceId    types.String `tfsdk:"deviceid"`
	N           types.String `tfsdk:"n"`
	LinePTPId   types.String `tfsdk:"lineptpid"`
	CarrierId   types.String `tfsdk:"carrierid"`
	Aid         types.String `tfsdk:"aid"`
	DscId       types.String `tfsdk:"dscid"`
	CDsc        types.Int64  `tfsdk:"cdsc"`
	TxStatus    types.String `tfsdk:"txstatus"`
	RxStatus    types.String `tfsdk:"rxstatus"`
	RelativeDPO types.Int64  `tfsdk:"relativedpo"`
	ConfigState    types.String `tfsdk:"configstate"`
}

// Metadata returns the data source type name.
func (r *DSCResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dsc"
}

// Configure adds the provider configured client to the resource.
func (r *DSCResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

// Schema defines the schema for the DSC resource.
func (r *DSCResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an DSC",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the DSC.",
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
			"aid": schema.StringAttribute{
				Description: "aid",
				Computed:    true,
			},
			"dscid": schema.StringAttribute{
				Description: "DSC ID",
				Optional:    true,
			},
			"cdsc": schema.Int64Attribute{
				Description: "constellation dsc ID",
				Computed:    true,
			},
			"txstatus": schema.StringAttribute{
				Description: "tx status",
				Computed:    true,
			},
			"rxstatus": schema.StringAttribute{
				Description: "Rx status",
				Computed:    true,
			},
			"relativedpo": schema.Int64Attribute{
				Description: "Relative DPO",
				Optional:    true,
			},
			"configstate": schema.StringAttribute{
				Description: "configstate",
				Computed:    true,
			},
		},
	}
}

func (r DSCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DSCResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "DSCResource: Create", map[string]interface{}{"DSCResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r DSCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DSCResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "DSCResource: Read", map[string]interface{}{"DSCResourceData": data})

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

func (r DSCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DSCResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "DSCResource: Update", map[string]interface{}{"DSCResourceData": data})
	// diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r DSCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DSCResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "DSCResource: Delete", map[string]interface{}{"DSCResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *DSCResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DSCResource) update(plan *DSCResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.LinePTPId.IsNull() || plan.CarrierId.IsNull() || plan.DscId.IsNull() {
		diags.AddError(
			"DSCResource: update ##: Error Update DSC",
			"Update: Could not Update DSC, Port ID, Carrier ID, DSC ID  are not specified",
		)
		return
	}
	tflog.Debug(ctx, "DSCResource: update ## ", map[string]interface{}{"LinePTPId": plan.LinePTPId.ValueString(), "Carrier": plan.CarrierId.ValueString(), "DSCID": plan.DscId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.RelativeDPO.IsNull()) {
		cmd["relativeDPO"] = plan.RelativeDPO.ValueInt64()
	}

	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"DSCResource: update ##: Error Update DSC",
			"Update: Could not Update DSC, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + plan.LinePTPId.ValueString() + "/carriers/" + plan.CarrierId.ValueString() + "/dscs/" + plan.DscId.ValueString()
	}

	tflog.Debug(ctx, "DSCResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"DSCResource: update ##: Error Update DSC",
			"Update: Could not Update DSC, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "DSCResource: update ##  ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"DSCResource: update ##: Error Update DSC",
			"Update: Could not SetResourceId DSC, unexpected error: "+err.Error(),
		)
		return
	}

	if content["cDsc"] != nil {
		plan.CDsc = types.Int64Value(int64(content["cDsc"].(float64)))
	}

	if content["txStatus"] != nil {
		plan.TxStatus = types.StringValue(content["txStatus"].(string))
	}

	if content["rxStatus"] != nil {
		plan.RxStatus = types.StringValue(content["rxStatus"].(string))
	}
	if content["configState"] != nil {
		plan.ConfigState = types.StringValue(content["configState"].(string))
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}

	tflog.Debug(ctx, "DSCResource: update ## ", map[string]interface{}{"plan": plan})
}

func (r DSCResource) read(state *DSCResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if state.LinePTPId.IsNull() || state.CarrierId.IsNull() || state.DscId.IsNull() {
		diags.AddError(
			"DSCResource: read ##: DSCResource: read ##: Error Read DSC",
			"Read: Could not Read DSC, Port ID, Carrier ID, and DSC ID are not specified",
		)
		return
	}

	href := after(state.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + state.LinePTPId.ValueString() + "/carriers/" + state.CarrierId.ValueString() + "/dscs/" + state.DscId.ValueString()
	}

	tflog.Debug(ctx, "DSCResource: read ## ", map[string]interface{}{"Device": state.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(state.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"DSCResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		state.Id = types.StringValue("")
		tflog.Debug(ctx, "DSCResource: read - not found ## 404", map[string]interface{}{"state": state})
		return
	}

	tflog.Debug(ctx, "DSCResource: read ## ", map[string]interface{}{"response": string(body)})

	state.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(state.N.ValueString(), &state.Id, body)
	if err != nil {
		diags.AddError(
			"DSCResource: read ##: Error Read DSC",
			"Read: Could not SetResourceId DSC, unexpected error: "+err.Error(),
		)
		return
	}
	if content["aid"] != nil {
		state.Aid = types.StringValue(content["aid"].(string))
	}
	if !state.RelativeDPO.IsNull()  {
		state.RelativeDPO = types.Int64Value(int64(content["relativeDPO"].(float64)))
	}

	if content["cDsc"] != nil {
		state.CDsc = types.Int64Value(int64(content["cDsc"].(float64)))
	}
	if content["txStatus"] != nil {
		state.TxStatus = types.StringValue(content["txStatus"].(string))
	}
	if content["rxStatus"] != nil {
		state.RxStatus = types.StringValue(content["rxStatus"].(string))
	}

	if content["configState"] != nil {
		state.ConfigState = types.StringValue(content["configState"].(string))
	}

	tflog.Debug(ctx, "DSCResource: read ## ", map[string]interface{}{"state": state})
}
