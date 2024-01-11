package provider

import (
	"context"
	"encoding/json"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &CarrierResource{}
	_ resource.ResourceWithConfigure   = &CarrierResource{}
	_ resource.ResourceWithImportState = &CarrierResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewCarrierResource() resource.Resource {
	return &CarrierResource{}
}

type CarrierResource struct {
	client *xrcm_pf.Client
}

type CarrierResourceData struct {
	Id                      types.String `tfsdk:"id"`
	N                       types.String `tfsdk:"n"`
	DeviceId                types.String `tfsdk:"deviceid"`
	LinePTPId               types.String `tfsdk:"lineptpid"`
	CarrierId               types.String `tfsdk:"carrierid"`
	Aid                     types.String `tfsdk:"aid"`
	Modulation              types.String `tfsdk:"modulation"`
	HModulation             types.String `tfsdk:"hmodulation"`
	OModulation             types.String `tfsdk:"omodulation"`
	ClientPortMode          types.String `tfsdk:"clientportmode"`
	FecIterations           types.String `tfsdk:"feciterations"`
	HFecIterations          types.String `tfsdk:"hfeciterations"`
	OFecIterations          types.String `tfsdk:"ofeciterations"`
	ConstellationFrequency  types.Int64  `tfsdk:"constellationfrequency"`
	HFrequency              types.Int64  `tfsdk:"hfrequency"`
	AConstellationFrequency types.Int64  `tfsdk:"aconstellationfrequency"`
	OperatingFrequency      types.Int64  `tfsdk:"operatingfrequency"`
	BaudRate                types.Int64  `tfsdk:"baudrate"`
	TxCLPtarget             types.Int64  `tfsdk:"txclptarget"`
	HTxCLPtarget            types.Int64  `tfsdk:"htxclptarget"`
	ATxCLPtarget            types.Int64  `tfsdk:"atxclptarget"`
	MaxDSCs                 types.Int64  `tfsdk:"maxdscs"`
	HMaxDSCs                types.Int64  `tfsdk:"hmaxdscs"`
	OMaxDSCs                types.Int64  `tfsdk:"omaxdscs"`
	MaxTxDSCs               types.Int64  `tfsdk:"maxtxdscs"`
	HMaxTxDSCs              types.Int64  `tfsdk:"hmaxtxdscs"`
	OMaxTxDSCs              types.Int64  `tfsdk:"omaxtxdscs"`
	AdvLineCtrl             types.String `tfsdk:"advlinectrl"`
	SpectralBandwidth       types.Int64  `tfsdk:"spectralbandwidth"`
	AllowedTxCDSCs          types.Int64  `tfsdk:"allowedtxcdscs"`
	HAllowedTxCDSCs         types.Int64  `tfsdk:"hallowedtxcdscs"`
	AAllowedTxCDSCs         types.Int64  `tfsdk:"aallowedtxcdscs"`
	AllowedRxCDSCs          types.Int64  `tfsdk:"allowedrxcdscs"`
	HAllowedRxCDSCs         types.Int64  `tfsdk:"hallowedrxcdscs"`
	AAllowedRxCDSCs         types.Int64  `tfsdk:"aallowedrxcdscs"`
	Capabilities            types.Map    `tfsdk:"capabilities"`
}

// Metadata returns the data source type name.
func (r *CarrierResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_carrier"
}

// Schema defines the schema for the resource.
func (r *CarrierResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Carrier",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the Carrier.",
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
			"constellationfrequency": schema.Int64Attribute{
				Description: "Constellation Frequency",
				Optional:    true,
			},
			"operatingfrequency": schema.Int64Attribute{
				Description: "operating frequency",
				Computed:    true,
			},
			"hfrequency": schema.Int64Attribute{
				Description: "Host frequency",
				Computed:    true,
			},
			"aconstellationfrequency": schema.Int64Attribute{
				Description: "Actual Constellation Frequency",
				Computed:    true,
			},
			"maxdscs": schema.Int64Attribute{
				Description: "Max Allowed DSCs",
				Optional:    true,
			},
			"hmaxdscs": schema.Int64Attribute{
				Description: "Host Max Allowed DSCs",
				Computed:    true,
			},
			"omaxdscs": schema.Int64Attribute{
				Description: "Operational Max Allowed DSCs",
				Computed:    true,
			},
			"maxtxdscs": schema.Int64Attribute{
				Description: "Max Tx DSCs",
				Optional:    true,
			},
			"hmaxtxdscs": schema.Int64Attribute{
				Description: "Host Max Tx DSCs",
				Computed:    true,
			},
			"omaxtxdscs": schema.Int64Attribute{
				Description: "Operational Max Tx DSCs",
				Computed:    true,
			},
			"allowedrxcdscs": schema.Int64Attribute{
				Description: "Allowed Rx DSCs",
				Optional:    true,
			},
			"hallowedrxcdscs": schema.Int64Attribute{
				Description: "Host Allowed Rx DSCs",
				Computed:    true,
			},
			"aallowedrxcdscs": schema.Int64Attribute{
				Description: "Actual Allowed Rx DSCs",
				Computed:    true,
			},
			"allowedtxcdscs": schema.Int64Attribute{
				Description: "Allowed Tx DSCs",
				Optional:    true,
			},
			"hallowedtxcdscs": schema.Int64Attribute{
				Description: "Host Allowed Tx DSCs",
				Computed:    true,
			},
			"aallowedtxcdscs": schema.Int64Attribute{
				Description: "Actual Allowed Tx DSCs",
				Computed:    true,
			},
			"txclptarget": schema.Int64Attribute{
				Description: "Tx CLP Target",
				Optional:    true,
			},
			"htxclptarget": schema.Int64Attribute{
				Description: "Host Tx CLP Target",
				Computed:    true,
			},
			"atxclptarget": schema.Int64Attribute{
				Description: "Actual Tx CLP Target",
				Computed:    true,
			},
			"advlinectrl": schema.StringAttribute{
				Description: "adv Line Ctrl",
				Computed:    true,
			},
			"spectralbandwidth": schema.Int64Attribute{
				Description: "Spectral bandwidth",
				Computed:    true,
			},
			"modulation": schema.StringAttribute{
				Description: "modulation",
				Optional:    true,
			},
			"omodulation": schema.StringAttribute{
				Description: "Operational modulationControl",
				Computed:    true,
			},
			"hmodulation": schema.StringAttribute{
				Description: "Host modulationControl",
				Computed:    true,
			},
			"clientportmode": schema.StringAttribute{
				Description: "client port mode",
				Optional:    true,
			},
			"feciterations": schema.StringAttribute{
				Description: "fec Iterations",
				Optional:    true,
			},
			"ofeciterations": schema.StringAttribute{
				Description: "Operational fec Iterations",
				Computed:    true,
			},
			"hfeciterations": schema.StringAttribute{
				Description: "Host fec Iterations",
				Computed:    true,
			},
			"baudrate": schema.Int64Attribute{
				Description: "baud rate",
				Optional:    true,
			},
			"capabilities": schema.MapAttribute{
				Description: "capabilities",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CarrierResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*xrcm_pf.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *CarrierResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan CarrierResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&plan, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r CarrierResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CarrierResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "CarrierResource: Read", map[string]interface{}{"CarrierResourceData": data})

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

func (r CarrierResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CarrierResourceData
	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "CarrierResource: Update", map[string]interface{}{"CarrierResourceData": data})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r CarrierResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CarrierResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "CarrierResource: Delete - ", map[string]interface{}{"CarrierResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CarrierResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CarrierResource) update(plan *CarrierResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.LinePTPId.IsNull() || plan.CarrierId.IsNull() {
		diags.AddError(
			"CarrierResource: update ##: Error Update Carrier",
			"Update: Could not Update Carrier, Port ID, Carrier ID  are not specified",
		)
		return
	}

	tflog.Debug(ctx, "CarrierResource: update ## ", map[string]interface{}{"LinePTPId": plan.LinePTPId.ValueString(), "Carrier": plan.CarrierId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.AllowedTxCDSCs.IsNull()) {
		cmd["allowedTxCDSCs"] = plan.AllowedTxCDSCs.ValueInt64()
	}

	if !(plan.AllowedRxCDSCs.IsNull()) {
		cmd["allowedRxCDSCs"] = plan.AllowedRxCDSCs.ValueInt64()
	}

	if !(plan.Modulation.IsNull()) {
		cmd["modulation"] = plan.Modulation.ValueString()
	}

	if !(plan.ClientPortMode.IsNull()) {
		cmd["ClientPortMode"] = plan.ClientPortMode.ValueString()
	}

	if !(plan.ConstellationFrequency.IsNull()) {
		cmd["constellationFrequency"] = plan.ConstellationFrequency.ValueInt64()
	}
	if !(plan.BaudRate.IsNull()) {
		cmd["baudRate"] = plan.BaudRate.ValueInt64()
	}

	if !(plan.MaxDSCs.IsNull()) {
		cmd["maxDSCs"] = plan.MaxDSCs.ValueInt64()
	}

	if !(plan.MaxTxDSCs.IsNull()) {
		cmd["maxTxDSCs"] = plan.MaxTxDSCs.ValueInt64()
	}

	if !(plan.TxCLPtarget.IsNull()) {
		cmd["txCLPtarget"] = plan.TxCLPtarget.ValueInt64()
	}
	if !(plan.FecIterations.IsNull()) {
		cmd["fecIterations"] = plan.FecIterations.ValueString()
	}
	if !(plan.AdvLineCtrl.IsNull()) {
		cmd["advLineCtrl"] = plan.AdvLineCtrl.ValueString()
	}

	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"CarrierResource: update ##: Error creating Carrier",
			"Update: not create Carrier, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + plan.LinePTPId.ValueString() + "/carriers/" + plan.CarrierId.ValueString()
	}

	tflog.Debug(ctx, "CarrierResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"CarrierResource: update ##: Error Update Carrier",
			"Update:Could not Update, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "CarrierResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	_, err = SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"CarrierResource: update ##: Error update LC",
			"Update: Could not SetResourceId CarrierResource, unexpected error: "+err.Error(),
		)
		return
	}
	/*for k, v := range content {
		switch k {
		case "aid":
			plan.Aid = types.StringValue(v.(string))
		case "hModulation":
			plan.HModulation = types.StringValue(v.(string))
		case "oModulation":
			plan.OModulation = types.StringValue(v.(string))
		case "hFecIterations":
			plan.HFecIterations = types.StringValue(v.(string))
		case "oFecIterations":
			plan.OFecIterations = types.StringValue(v.(string))
		case "hFrequency":
			plan.HFrequency = types.Int64Value(int64(v.(float64)))
		case "aConstellationFrequency":
			plan.AConstellationFrequency = types.Int64Value(int64(v.(float64)))
		case "operatingFrequency":
			plan.OperatingFrequency = types.Int64Value(int64(v.(float64)))
		case "spectralBandwidth":
			plan.SpectralBandwidth = types.Int64Value(int64(v.(float64)))
		case "hTxCLPtarget":
			plan.HTxCLPtarget = types.Int64Value(int64(v.(float64)))
		case "aTxCLPtarget":
			plan.ATxCLPtarget = types.Int64Value(int64(v.(float64)))
		case "hMaxDSCs":
			plan.HMaxDSCs = types.Int64Value(int64(v.(float64)))
		case "oMaxDSCs":
			plan.OMaxDSCs = types.Int64Value(int64(v.(float64)))
		case "hMaxTxDSCs":
			plan.HMaxTxDSCs = types.Int64Value(int64(v.(float64)))
		case "oMaxTxDSCs":
			plan.OMaxTxDSCs = types.Int64Value(int64(v.(float64)))
		case "hAllowedTxCDSCs":
			plan.HAllowedTxCDSCs = types.Int64Value(int64(v.(float64)))
		case "aAllowedTxCDSCs":
			plan.AAllowedTxCDSCs = types.Int64Value(int64(v.(float64)))
		case "hAllowedRxCDSCs":
			plan.HAllowedRxCDSCs = types.Int64Value(int64(v.(float64)))
		case "aAllowedRxCDSCs":
			plan.AAllowedRxCDSCs = types.Int64Value(int64(v.(float64)))
		case "capabilities":
			plan.Capabilities, _ = types.MapValue(types.StringType, v.(map[string]attr.Value))
		}
	}*/
	r.read(plan, ctx, diags)
	/*if err != nil {
		diags.AddError(
			"CarrierResource: update ##: Error Update Carrier",
			"Update: Could not SetData Carrier, unexpected error: "+err.Error(),
		)
		return
	}*/

	tflog.Debug(ctx, "CarrierResource: update ## ", map[string]interface{}{"plan": plan})

}

func (r *CarrierResource) read(state *CarrierResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if state.LinePTPId.IsNull() || state.CarrierId.IsNull() {
		diags.AddError(
			"CarrierResource: read ##: Error Read Carrier",
			"Read: Could not Read  Carrier, Port ID, Carrier ID are not specified",
		)
		return
	}

	href := after(state.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + state.LinePTPId.ValueString() + "/carriers/" + state.CarrierId.ValueString()
	}

	tflog.Debug(ctx, "CarrierResource: read ## ", map[string]interface{}{"device": state.N.ValueString, "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(state.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"CarrierResource: read ##: Error Read Carrier",
				"Read: Could not get Carrier, unexpected error: "+err.Error(),
			)
			return
		}
		state.Id = types.StringValue("")
		tflog.Debug(ctx, "CarrierResource: Read - Not found  ## 404", map[string]interface{}{"state": state})
		return
	}

	tflog.Debug(ctx, "CarrierResource: read ## ", map[string]interface{}{"response": string(body)})

	state.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(state.N.ValueString(), &state.Id, body)

	if err != nil {
		diags.AddError(
			"CarrierResource: read ##: Error SetData Carrier",
			"Read: Could not SetData Carrier, unexpected error: "+err.Error(),
		)
		return
	}

	for k, v := range content {
		switch k {
		case "aid":
			state.Aid = types.StringValue(v.(string))
		case "FecIterations":
			if !(state.FecIterations.IsNull()) {
				state.FecIterations = types.StringValue(v.(string))
			}
		case "advLineCtrl":
			if !(state.AdvLineCtrl.IsNull()) {
				state.AdvLineCtrl = types.StringValue(v.(string))
			}
		case "modulation":
			if !(state.Modulation.IsNull()) {
				state.Modulation = types.StringValue(v.(string))
			}
		case "clientPortMode":
			if !(state.ClientPortMode.IsNull()) {
				state.ClientPortMode = types.StringValue(v.(string))
			}
		case "constellationFrequency":
			if !(state.ConstellationFrequency.IsNull()) {
				state.ConstellationFrequency = types.Int64Value(int64(v.(float64)))
			}
		case "baudRate":
			if !(state.BaudRate.IsNull()) {
				state.BaudRate = types.Int64Value(int64(v.(float64)))
			}
		case "maxDSCs":
			if !(state.MaxDSCs.IsNull()) {
				state.MaxDSCs = types.Int64Value(int64(v.(float64)))
			}
		case "maxTxDSCs":
			if !(state.MaxTxDSCs.IsNull()) {
				state.MaxTxDSCs = types.Int64Value(int64(v.(float64)))
			}
		case "spectralBandwidth":
			state.SpectralBandwidth = types.Int64Value(int64(v.(float64)))
		case "txCLPtarget":
			if !(state.TxCLPtarget.IsNull()) {
				state.TxCLPtarget = types.Int64Value(int64(v.(float64)))
			}
		case "allowedTxCDSCs":
			if !(state.AllowedTxCDSCs.IsNull()) {
				state.AllowedTxCDSCs = types.Int64Value(int64(v.(float64)))
			}
		case "allowedRxCDSCs":
			if !(state.AllowedRxCDSCs.IsNull()) {
				state.AllowedRxCDSCs = types.Int64Value(int64(v.(float64)))
			}
		case "hModulation":
			state.HModulation = types.StringValue(v.(string))
		case "oModulation":
			state.OModulation = types.StringValue(v.(string))
		case "hFecIterations":
			state.HFecIterations = types.StringValue(v.(string))
		case "oFecIterations":
			state.OFecIterations = types.StringValue(v.(string))
		case "hFrequency":
			state.HFrequency = types.Int64Value(int64(v.(float64)))
		case "aConstellationFrequency":
			state.AConstellationFrequency = types.Int64Value(int64(v.(float64)))
		case "operatingFrequency":
			state.OperatingFrequency = types.Int64Value(int64(v.(float64)))
		case "hTxCLPtarget":
			state.HTxCLPtarget = types.Int64Value(int64(v.(float64)))
		case "aTxCLPtarget":
			state.ATxCLPtarget = types.Int64Value(int64(v.(float64)))
		case "hMaxDSCs":
			state.HMaxDSCs = types.Int64Value(int64(v.(float64)))
		case "oMaxDSCs":
			state.OMaxDSCs = types.Int64Value(int64(v.(float64)))
		case "hMaxTxDSCs":
			state.HMaxTxDSCs = types.Int64Value(int64(v.(float64)))
		case "oMaxTxDSCs":
			state.OMaxTxDSCs = types.Int64Value(int64(v.(float64)))
		case "hAllowedTxCDSCs":
			state.HAllowedTxCDSCs = types.Int64Value(int64(v.(float64)))
		case "aAllowedTxCDSCs":
			state.AAllowedTxCDSCs = types.Int64Value(int64(v.(float64)))
		case "hAllowedRxCDSCs":
			state.HAllowedRxCDSCs = types.Int64Value(int64(v.(float64)))
		case "aAllowedRxCDSCs":
			state.AAllowedRxCDSCs = types.Int64Value(int64(v.(float64)))
		case "capabilities":
			capMap := make(map[string]attr.Value)
			for k,v2 := range v.(map[string]interface{})  {
				capMap[k] = types.StringValue(v2.(string))
			}
			state.Capabilities, _ = types.MapValue(types.StringType, capMap)
		}
	}
	tflog.Debug(ctx, "CarrierResource: read ## ", map[string]interface{}{"state": state})
}
