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
	_ resource.Resource                = &DSCGResource{}
	_ resource.ResourceWithConfigure   = &DSCGResource{}
	_ resource.ResourceWithImportState = &DSCGResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewDSCGResource() resource.Resource {
	return &DSCGResource{}
}

type DSCGResource struct {
	client *xrcm_pf.Client
}

type DSCGResourceData struct {
	Id        types.String `tfsdk:"id"`
	DeviceId  types.String `tfsdk:"deviceid"`
	Aid       types.String `tfsdk:"aid"`
	N         types.String `tfsdk:"n"`
	LinePTPId types.String `tfsdk:"lineptpid"`
	CarrierId types.String `tfsdk:"carrierid"`
	DscgId    types.String `tfsdk:"dscgid"`
	TxCDSCs   types.List   `tfsdk:"txcdscs"`
	RxCDSCs   types.List   `tfsdk:"rxcdscs"`
	IdleCDSCs types.List   `tfsdk:"idlecdscs"`
	DscgCtrl  types.Int64  `tfsdk:"dscgctrl"`
}

// Metadata returns the data source type name.
func (r *DSCGResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dscg"
}

// Schema defines the schema for the DSCG resource.
func (r *DSCGResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an DSCG",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the DSCG.",
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
			"dscgid": schema.StringAttribute{
				Description: "DSCG id",
				Optional:    true,
			},
			"txcdscs": schema.ListAttribute{ElementType: types.Int64Type,
				Optional: true, Description: "Transmitting Constellation DSC IDs",
			},
			"rxcdscs": schema.ListAttribute{ElementType: types.Int64Type,
				Optional: true, Description: "Receiving Constellation DSC IDs",
			},
			"idlecdscs": schema.ListAttribute{ElementType: types.Int64Type,
				Optional: true, Description: "Idle Constellation DSC IDs",
			},
			"dscgctrl": schema.Int64Attribute{
				Description: "dscg ctrl",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DSCGResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r DSCGResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DSCGResourceData

	diags := req.Config.Get(ctx, &data)
	tflog.Debug(ctx, "DSCGResource: Create", map[string]interface{}{"DSCGResourceData": data})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.create(&data, ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r DSCGResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DSCGResourceData
	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "DSCGResource: Read", map[string]interface{}{"DSCGResourceData": data})
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

func (r DSCGResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data DSCGResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DSCGResource: Update", map[string]interface{}{"DSCGResourceData": data})

	r.delete(&data, ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	r.create(&data, ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}

func (r DSCGResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DSCGResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "DSCGResource: Delete", map[string]interface{}{"DSCGResourceData": data})
	resp.Diagnostics.Append(diags...)

	r.delete(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *DSCGResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r DSCGResource) create(plan *DSCGResourceData, ctx context.Context, diags *diag.Diagnostics) {
	if plan.LinePTPId.IsNull() || plan.CarrierId.IsNull() {
		diags.AddError(
			"DSCGResource: create ##: Error Create DSCG",
			"Create: Could not Create DSCG, Port ID or Carrier ID are not specified",
		)
		return
	}

	tflog.Debug(ctx, "DSCGResource: create ## ", map[string]interface{}{"LinePTPId": plan.LinePTPId.ValueString(), "Carrier": plan.CarrierId.ValueString()})

	var rep = make(map[string]interface{})

	if !(plan.RxCDSCs.IsNull()) {
		var rxCDSCList []int
		diag := plan.RxCDSCs.ElementsAs(ctx, &rxCDSCList, true)
		if diag != nil && diag.HasError() {
			diags.AddError(
				"DSCGResource: create ##: Error Create DSCG",
				"Create: Could not Create DSCG, RxCDSCs is invalid "+plan.RxCDSCs.String(),
			)
			return
		}
		rxCDSCs := setBits(rxCDSCList)
		rep["rxCDSCs"] = rxCDSCs
	}

	if !(plan.TxCDSCs.IsNull()) {
		var txCDSCList []int
		diag := plan.TxCDSCs.ElementsAs(ctx, &txCDSCList, true)
		if diag != nil && diag.HasError() {
			diags.AddError(
				"DSCGResource: create ##: Error Create DSCG",
				"Create: Could not Create DSCG, TxCDSCs is invalid "+plan.TxCDSCs.String(),
			)
			return
		}
		txCDSCs := setBits(txCDSCList)
		rep["txCDSCs"] = txCDSCs
	}

	if !(plan.IdleCDSCs.IsNull()) {
		var idleCDSCList []int
		diag := plan.IdleCDSCs.ElementsAs(ctx, &idleCDSCList, true)
		if diag != nil && diag.HasError() {
			diags.AddError(
				"DSCGResource: create ##: Error Create DSCG",
				"Create: Could not Create DSCG, TxCDSCs is invalid "+plan.TxCDSCs.String(),
			)
			return
		}
		idleCDSCs := setBits(idleCDSCList)
		rep["idleCDSCs"] = idleCDSCs
	}

	if !(plan.DscgCtrl.IsNull()) {
		rep["dscgCtrl"] = plan.DscgCtrl.ValueInt64()
	}

	var cmd = make(map[string]interface{})

	cmd["rep"] = rep

	var ifs []string
	ifs = append(ifs, "oic.if.baseline", "oic.if.rw", "oic.if.delete")
	cmd["if"] = ifs
	var rt []string
	rt = append(rt, "xr.carrier.dscg")
	cmd["rt"] = rt
	var p = make(map[string]int)
	p["bm"] = 3
	cmd["p"] = p

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"DSCGResource: create ##: Error Create DSCG",
			"Create: Could not Marshal DSCG, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "DSCGResource: create ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resource-links/lineptps/" + plan.LinePTPId.ValueString() + "/carriers/" + plan.CarrierId.ValueString() + "/dscgs", "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "POST", "resource-links/lineptps/"+plan.LinePTPId.ValueString()+"/carriers/"+plan.CarrierId.ValueString()+"/dscgs", rb)

	tflog.Debug(ctx, "DSCGResource: create ##  ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	if err != nil {
		diags.AddError(
			"DSCGResource: create ##: Error Create DSCG",
			"Create: Could not POST DSCG, unexpected error: "+err.Error(),
		)
		return
	}

	plan.DeviceId = types.StringValue(deviceId)
	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"DSCGResource: create ##: Error Create DSCG",
			"Create: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	rep1 := content["rep"].(map[string]interface{})
	aid := rep1["aid"]
	if aid != nil {
		plan.Aid = types.StringValue(aid.(string))
	} 

	tflog.Debug(ctx, "DSCGResource: create ## ", map[string]interface{}{"plan": plan})
}

func (r DSCGResource) read(state *DSCGResourceData, ctx context.Context, diags *diag.Diagnostics) {
	if state.DscgId.IsNull() {
		diags.AddError(
			"DSCGResource: read ##: Error Read DSCG",
			"Read: Could not Read DSCG, DSCG ID must specify",
		)
		return
	}

	href := after(state.Id.ValueString(), "/")
	if len(href) == 0 {
		state.Id = types.StringValue("")
		tflog.Debug(ctx, "DSCGResource: read - href is empty", map[string]interface{}{"state": state})
		return
	}

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(state.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"DSCGResource: read ##: Error Read DSCG",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		state.Id = types.StringValue("")
		tflog.Debug(ctx, "DSCGResource: read - not found ## 404", map[string]interface{}{"plan": state})
		return
	}

	tflog.Debug(ctx, "DSCGResource: read ## ", map[string]interface{}{"response": string(body)})

	state.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(state.N.ValueString(), &state.Id, body)
	if err != nil {
		diags.AddError(
			"DSCGResource: read ##: Error Read DSCG",
			"Read: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	for k, v := range content {
		switch k {
		case "aid":
			if v != nil {
				state.Aid = types.StringValue(v.(string))
			}
		case "rxCDSCs":
			if !(state.RxCDSCs.IsNull()) {
				rxCDSCList := getBits(int(v.(float64)))
				state.RxCDSCs, _ = types.ListValue(types.Int64Type, rxCDSCList)
			}
		case "txCDSCs":
			if !(state.TxCDSCs.IsNull()) {
				txCDSCList := getBits(int(v.(float64)))
				state.TxCDSCs, _ = types.ListValue(types.Int64Type, txCDSCList)
			}
		case "idleCDSCs":
			if !(state.IdleCDSCs.IsNull()) {
				idleCDSCList := getBits(int(v.(float64)))
				state.IdleCDSCs, _ = types.ListValue(types.Int64Type, idleCDSCList)
			} 
		case "dscgCtrl":
			if !(state.DscgCtrl.IsNull()) {
				state.DscgCtrl = types.Int64Value(int64(v.(float64)))
				}
		}
	}
	tflog.Debug(ctx, "DSCGResource: read ## ", map[string]interface{}{"plan": state})
}

func (r *DSCGResource) delete(plan *DSCGResourceData, ctx context.Context, diags *diag.Diagnostics) {

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		diags.AddError(
			"Error Delete DSCG",
			"Delete DSCG failed!. HREF is empty",
		)
		return
	}

	tflog.Debug(ctx, "DSCGResource: delete ## ", map[string]interface{}{"href": href})

	body, _, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "DELETE", "resource-links"+href, nil)

	if err != nil && !strings.Contains(err.Error(), "status: 404") {
		diags.AddError(
			"DSCGResource: delete ##: Error Delete LC",
			"Delete: Could not Delete LC, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "DSCGResource: delete ##  ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})
}

func (r *DSCGResource) update(plan *DSCGResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.LinePTPId.IsNull() || plan.CarrierId.IsNull() || plan.DscgId.IsNull() || len(plan.DscgId.ValueString()) == 0 {
		diags.AddError(
			"DSCGResource: update ##: Error Update DSCG",
			"Update: Could not Update DSCG, Port ID, Carrier ID, and DSCG ID are not specified",
		)
		return
	}
	tflog.Debug(ctx, "DSCGResource: update ## ", map[string]interface{}{"LinePTPId": plan.LinePTPId.ValueString(), "Carrier": plan.CarrierId.ValueString()})

	var cmd = make(map[string]interface{})

	var rxCDSCList []int
	diag := plan.RxCDSCs.ElementsAs(ctx, &rxCDSCList, true)
	if diag != nil && diag.HasError() {
		diags.AddError(
			"DSCGResource: create ##: Error Create DSCG",
			"Create: Could not Create DSCG, RxCDSCs is invalid"+plan.RxCDSCs.String(),
		)
		return
	}
	rxCDSCs := setBits(rxCDSCList)
	cmd["rxCDSCs"] = rxCDSCs

	var txCDSCList []int
	diag = plan.RxCDSCs.ElementsAs(ctx, &txCDSCList, true)
	if diag != nil && diag.HasError() {
		diags.AddError(
			"DSCGResource: create ##: Error Create DSCG",
			"Create: Could not Create DSCG, TxCDSCs is invalid"+plan.TxCDSCs.String(),
		)
		return
	}
	txCDSCs := setBits(txCDSCList)
	cmd["txCDSCs"] = txCDSCs

	if !(plan.DscgCtrl.IsNull()) {
		cmd["dscgCtrl"] = plan.DscgCtrl.ValueInt64()
	}

	if len(cmd) == 0. { // nothing to update
		return
	}

	rb, err := json.Marshal(cmd)
	if err != nil {
		diags.AddError(
			"DSCGResource: update ##: Error Update DSCG",
			"Update: Could not Marshal DSCG, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "DSCGResource: read - href is empty", map[string]interface{}{"plan": plan})
		return
	}

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	tflog.Debug(ctx, "DSCG Update", map[string]interface{}{"rb=": string(body)})
	if err != nil {
		diags.AddError(
			"DSCGResource: update ##: Error Update DSCG",
			"Update: Could not PUT DSCG, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "DSCGResource: update ##  ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"DSCGResource: update ##: Error Create LC",
			"Create: Could not LC SetResourceId, unexpected error: "+err.Error(),
		)
		return
	}

	aid := content["aid"]
	if aid != nil && len(aid.(string)) > 0 {
		plan.Aid = types.StringValue(aid.(string))
	} else {
		plan.Aid = types.StringValue("")
	}

	if content["idleCDSCs"] != nil {
		idleCDSCList := getBits(int(content["idleCDSCs"].(float64)))
		plan.IdleCDSCs, _ = types.ListValue(types.Int64Type, idleCDSCList)
	}

	tflog.Debug(ctx, "DSCResource: update ## ", map[string]interface{}{"plan": plan})
}
