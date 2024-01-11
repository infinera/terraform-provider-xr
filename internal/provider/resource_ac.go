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
	_ resource.Resource                = &ACResource{}
	_ resource.ResourceWithConfigure   = &ACResource{}
	_ resource.ResourceWithImportState = &ACResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewACResource() resource.Resource {
	return &ACResource{}
}

type ACResource struct {
	client *xrcm_pf.Client
}

type ACResourceData struct {
	Id          types.String `tfsdk:"id"`
	N           types.String `tfsdk:"n"`
	DeviceId    types.String `tfsdk:"deviceid"`
	EthernetId  types.String `tfsdk:"ethernetid"`
	AcId        types.String `tfsdk:"acid"`
	Aid         types.String `tfsdk:"aid"`
	Capacity    types.Int64  `tfsdk:"capacity"`
	Imc         types.String `tfsdk:"imc"`
	ImcOuterVID types.String `tfsdk:"imc_outer_vid"`
	Emc         types.String `tfsdk:"emc"`
	EmcOuterVID types.String `tfsdk:"emc_outer_vid"`
	AcCtrl      types.Int64  `tfsdk:"acctrl"`
}

// Metadata returns the data source type name.
func (r *ACResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ac"
}

// Schema defines the schema for the data source.
func (r *ACResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an AC",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the AC.",
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
			"acid": schema.StringAttribute{
				Description: "AC id",
				Required:    true,
			},
			"aid": schema.StringAttribute{
				Description: "Aid",
				Computed:    true,
			},
			"ethernetid": schema.StringAttribute{
				Description: "ethernet id",
				Optional:    true,
			},
			"acctrl": schema.Int64Attribute{
				Description: "AC Control",
				Optional:    true,
			},
			"capacity": schema.Int64Attribute{
				Description: "capacity",
				Optional:    true,
			},
			"imc": schema.StringAttribute{
				Description: "imc",
				Optional:    true,
			},
			"imc_outer_vid": schema.StringAttribute{
				Description: "imc outer vid",
				Optional:    true,
			},
			"emc": schema.StringAttribute{
				Description: "emc",
				Optional:    true,
			},
			"emc_outer_vid": schema.StringAttribute{
				Description: "emc_outer_vid",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ACResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r ACResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ACResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "ACResource: Create - ", map[string]interface{}{"ACResourceData": data})

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.create(&data, ctx, &resp.Diagnostics)

	if data.Id.IsNull() {
		resp.State = tfsdk.State{}
	} else {
		diags = resp.State.Set(ctx, &data)
	}

	resp.Diagnostics.Append(diags...)
}

func (r ACResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ACResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "ACResource: Create - ", map[string]interface{}{"ACResourceData": data})

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.read(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ACResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ACResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "CfgResource: Update", map[string]interface{}{"ACResourceData": data})

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ACResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ACResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "CfgResource: Update", map[string]interface{}{"ACResourceData": data})

	resp.Diagnostics.Append(diags...)

	r.delete(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ACResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ACResource) create(plan *ACResourceData, ctx context.Context, diags *diag.Diagnostics) {
	if plan.AcId.IsNull() || plan.EthernetId.IsNull() {
		diags.AddError(
			"Error Create AC",
			"Create: Could not create AC, AC ID or Ethernet Id is not specified",
		)
		return
	}
	tflog.Debug(ctx, "ACResource: create ## ", map[string]interface{}{"Acid": plan.AcId.ValueString(), "EthernetId": plan.EthernetId.ValueString()})
	var rep = make(map[string]interface{})

	if !(plan.Capacity.IsNull()) {
		rep["capacity"] = plan.Capacity.ValueInt64()
	}

	if !(plan.AcCtrl.IsNull()) {
		rep["acCtrl"] = plan.AcCtrl.ValueInt64()
	}

	if !(plan.Imc.IsNull()) {
		rep["imc"] = plan.Imc.ValueString()
	}

	if !(plan.ImcOuterVID.IsNull()) {
		rep["imcOuterVID"] = plan.ImcOuterVID.ValueString()
	}

	if !(plan.Emc.IsNull()) {
		rep["emc"] = plan.Emc.ValueString()
	}

	if !(plan.EmcOuterVID.IsNull()) {
		rep["emcOuterVID"] = plan.EmcOuterVID.ValueString()
	}

	var cmd = make(map[string]interface{})

	cmd["rep"] = rep

	var ifs []string
	ifs = append(ifs, "oic.if.baseline", "oic.if.rw")
	cmd["if"] = ifs
	var rt []string
	rt = append(rt, "xr.ethernet.ac")
	cmd["rt"] = rt
	var p = make(map[string]int)
	p["bm"] = 1
	cmd["p"] = p
	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"ACResource: create ##: Error Create AC",
			"Create: Could not Marshal AC, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "ACResource: create ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resource-links/ethernets/" + plan.EthernetId.ValueString() + "/acs", "cmd": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "POST", "resource-links/ethernets/"+plan.EthernetId.ValueString()+"/acs", rb)

	if err != nil {
		diags.AddError(
			"ACResource: create ##: Error creating AC",
			"Create: Could not create AC, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "ACResource: create ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"ACResource: create ##: Error Create LC",
			"Create: Could not AC SetResourceId, unexpected error: "+err.Error(),
		)
		return
	}
	rep = content["rep"].(map[string]interface{})
	aid := rep["aid"]
	if aid != nil && len(aid.(string)) > 0 {
		plan.Aid = types.StringValue(aid.(string))
	} else {
		plan.Aid = types.StringValue("")
	}
	plan.DeviceId = types.StringValue(deviceId)

	tflog.Debug(ctx, "ACResource: create ##", map[string]interface{}{"plan": plan})
}

func (r *ACResource) read(plan *ACResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.AcId.IsNull() || plan.EthernetId.IsNull() {
		diags.AddError(
			"Error get AC",
			"Create: Could not get AC, AC ID or Ethernet Id is not specified",
		)
		return
	}

	tflog.Debug(ctx, "ACResource: read ## ", map[string]interface{}{"Acid": plan.AcId.ValueString(), "EthernetId": plan.EthernetId.ValueString()})

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "ACResource: read - href is empty", map[string]interface{}{"plan": plan})
		return
	}

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"ACResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "ACResource: read - not found ## 404", map[string]interface{}{"plan": plan})
		return
	}

	tflog.Debug(ctx, "ACResource: read ## ", map[string]interface{}{"response": string(body)})

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"LCResource: read ##: Error Read LC",
			"Create: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	plan.DeviceId = types.StringValue(deviceId)

	for k, v := range content {
		switch k {
		case "aid":
			{
				plan.Aid = types.StringValue(v.(string))
			}
		case "capacity":
			if !(plan.Capacity.IsNull()) {
				x := int64(v.(float64))
				plan.Capacity = types.Int64Value(x)
			}
		case "acCtrl":
			if !(plan.AcCtrl.IsNull()) {
				x := int64(v.(float64))
				plan.Capacity = types.Int64Value(x)
			}
		case "imc":
			if !(plan.Imc.IsNull()) {
				plan.Imc = types.StringValue(v.(string))
			}
		case "emc":
			if !(plan.Emc.IsNull()) {
				plan.Emc = types.StringValue(v.(string))
			}
		case "imcOuterVID":
			if !(plan.ImcOuterVID.IsNull()) {
				plan.ImcOuterVID = types.StringValue(v.(string))
			}
		case "emcOuterVID":
			if !(plan.EmcOuterVID.IsNull()) {
				plan.EmcOuterVID = types.StringValue(v.(string))
			}
		}
	}
	tflog.Debug(ctx, "ACResource: read ## ", map[string]interface{}{"plan": plan})
}

func (r *ACResource) update(plan *ACResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.AcId.IsNull() || plan.EthernetId.IsNull() {
		diags.AddError(
			"Error Update AC",
			"Create: Could not Update AC, AC ID or Ethernet Id is not specified",
		)
		return
	}

	tflog.Debug(ctx, "ACResource: create ## ", map[string]interface{}{"Acid": plan.AcId.ValueString(), "EthernetId": plan.EthernetId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.Capacity.IsNull()) {
		cmd["capacity"] = plan.Capacity.ValueInt64()
	}

	if !(plan.Imc.IsNull()) {
		cmd["imc"] = plan.Imc.ValueString()
	}

	if !(plan.ImcOuterVID.IsNull()) {
		cmd["imcOuterVID"] = plan.ImcOuterVID.ValueString()
	}

	if !(plan.Emc.IsNull()) {
		cmd["emc"] = plan.Emc.ValueString()
	}

	if !(plan.EmcOuterVID.IsNull()) {
		cmd["emcOuterVID"] = plan.EmcOuterVID.ValueString()
	}

	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"ACResource: update ##: Error creating ACResource",
			"Update: not create ACResource, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + plan.EthernetId.ValueString() + "/acs/" + plan.AcId.ValueString()
	}

	tflog.Debug(ctx, "ACResource: Update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resource-links" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"ACResource: update ##: Error Update Carrier",
			"Update:Could not Update, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "ACResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"ACResource: update ##: Error update LC",
			"Update: Could not SetResourceId ACResource, unexpected error: "+err.Error(),
		)
		return
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}

	tflog.Debug(ctx, "ACResource: update ## ", map[string]interface{}{"plan": plan})
}

func (r *ACResource) delete(plan *ACResourceData, ctx context.Context, diags *diag.Diagnostics) {

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		return
	}

	tflog.Debug(ctx, "ACResource: delete ## ", map[string]interface{}{"href": href})

	body, _, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "DELETE", "resource-links"+href, nil)

	if err != nil && !strings.Contains(err.Error(), "status: 404") {
		diags.AddError(
			"Error Delete LC",
			"Delete: Could not Delete LC, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "ACResource: delete ##  ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

}
