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
	_ resource.Resource                = &LCResource{}
	_ resource.ResourceWithConfigure   = &LCResource{}
	_ resource.ResourceWithImportState = &LCResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewLCResource() resource.Resource {
	return &LCResource{}
}

type LCResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *LCResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lc"
}

type LCResourceData struct {
	Id             types.String `tfsdk:"id"`
	N              types.String `tfsdk:"n"`
	DeviceId       types.String `tfsdk:"deviceid"`
	Aid            types.String `tfsdk:"aid"`
	LcCtrl         types.Int64  `tfsdk:"lcctrl"`
	LinePTPId      types.String `tfsdk:"lineptpid"`
	Direction      types.String `tfsdk:"direction"`
	ClientAid      types.String `tfsdk:"clientaid"`
	LineAid        types.String `tfsdk:"lineaid"`
	DscgAid        types.String `tfsdk:"dscgaid"`
	RemoteModuleId types.String `tfsdk:"remotemoduleid"`
	RemoteClientId types.String `tfsdk:"remoteclientid"`
}

// Schema defines the schema for the resource.
func (r *LCResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"aid": schema.StringAttribute{
				Description: "aid",
				Computed:    true,
			},
			"lcctrl": schema.Int64Attribute{
				Description: "LC control",
				Optional:    true,
			},
			"direction": schema.StringAttribute{
				Description: "direction",
				Optional:    true,
			},
			"lineaid": schema.StringAttribute{
				Description: "line aid",
				Computed:    true,
			},
			"clientaid": schema.StringAttribute{
				Description: "client aid",
				Optional:    true,
			},
			"dscgaid": schema.StringAttribute{
				Description: "dscg aid",
				Optional:    true,
			},
			"remotemoduleid": schema.StringAttribute{
				Description: "remote module id",
				Computed:    true,
			},
			"remoteclientid": schema.StringAttribute{
				Description: "remote client id",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *LCResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r LCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LCResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "LCResource: Create", map[string]interface{}{"LCResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.create(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r LCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LCResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "LCResource: Read", map[string]interface{}{"LCResourceData": data})

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

func (r LCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LCResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "LCResource: Update", map[string]interface{}{"LCResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r LCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LCResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "LCResource: Delete", map[string]interface{}{"LCResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.delete(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *LCResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *LCResource) create(plan *LCResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.ClientAid.IsNull() || plan.DscgAid.IsNull() || plan.LinePTPId.IsNull() {
		diags.AddError(
			"LCResource: create ##: Error Create LC",
			"Create: ClientAid, DscgAid, LinePTPId and CarrierId must specify to create LC:",
		)
		return
	}

	tflog.Debug(ctx, "LCResource: create ## ", map[string]interface{}{"ClientAid": plan.ClientAid.ValueString(), "DscgAid": plan.DscgAid.ValueString(), "LinePTPId": plan.LinePTPId.ValueString()})

	//create LC
	var rep = make(map[string]interface{})

	rep["clientAid"] = plan.ClientAid.ValueString()

	rep["dscgAid"] = plan.DscgAid.ValueString()

	if !(plan.LcCtrl.IsNull()) {
		rep["lcCtrl"] = plan.LcCtrl.ValueInt64()
	}

	if !(plan.Direction.IsNull()) {
		rep["direction"] = plan.Direction.ValueString()
	}

	var cmd = make(map[string]interface{})

	cmd["rep"] = rep

	var ifs []string
	ifs = append(ifs, "oic.if.baseline", "oic.if.a")
	cmd["if"] = ifs
	var rt []string
	rt = append(rt, "xr.lc")
	cmd["rt"] = rt
	var p = make(map[string]int)
	p["bm"] = 3
	cmd["p"] = p

	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"LCResource: create ##: Error Create LC",
			"Create: Could not Marshal LC, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "LCResource: create ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resource-links/lcs", "rb": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "POST", "resource-links/lcs", rb)

	if err != nil {
		diags.AddError(
			"LCResource: create ##: Error creating LC",
			"Create: Could not create LC, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "LCResource: create ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"LCResource: create ##: Error Create LC",
			"Create: Could not create LC , unexpected error: "+err.Error(),
		)
		return
	}

	rep = content["rep"].(map[string]interface{})
	aid := rep["aid"]
	if aid != nil {
		plan.Aid = types.StringValue(aid.(string))
	}

	lineAid := rep["lineAid"]
	if lineAid != nil {
		plan.LineAid = types.StringValue(lineAid.(string))
	}

	remoteModuleId := rep["remoteModuleId"]
	if remoteModuleId != nil {
		plan.RemoteModuleId = types.StringValue(remoteModuleId.(string))
	}

	remoteClientId := rep["remoteClientId"]
	if remoteModuleId != nil {
		plan.RemoteClientId = types.StringValue(remoteClientId.(string))
	}

	plan.DeviceId = types.StringValue(deviceId)

	tflog.Debug(ctx, "LCResource: create ##", map[string]interface{}{"plan": plan})
}

func (r *LCResource) update(plan *LCResourceData, ctx context.Context, diags *diag.Diagnostics) {

	tflog.Debug(ctx, "LCResource: update ## ", map[string]interface{}{"ClientAid": plan.ClientAid.ValueString(), "DscgAid": plan.DscgAid.ValueString(), "LinePTPId": plan.LinePTPId.ValueString()})

	/*diags.AddError(
		"LCResource: update ##: Error Update LC",
		"Update: Update LC is not allowed",
	)*/
}

func (r *LCResource) read(plan *LCResourceData, ctx context.Context, diags *diag.Diagnostics) {

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "LCResource: read - href is empty", map[string]interface{}{"plan": plan})
		return
	}

	tflog.Debug(ctx, "LCResource: read ## ", map[string]interface{}{"ClientAid": plan.ClientAid.ValueString(), "DscgAid": plan.DscgAid.ValueString(), "LinePTPId": plan.LinePTPId.ValueString(), "href": href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"LCResource: read ##: Error Read LC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "LCResource: read - Not Found ## 404", map[string]interface{}{"plan": plan})
		return
	}

	tflog.Debug(ctx, "LCResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"LCResource: read ##: Error Read LC",
			"Create: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	for k, v := range content {
		switch k {
		case "aid":
			if len(v.(string)) > 0 {
				plan.Aid = types.StringValue(v.(string))
			}
		case "lcCtrl":
			if !(plan.LcCtrl.IsNull()) {
				plan.LcCtrl = types.Int64Value(int64(v.(float64)))
			}
		case "direction":
			if !(plan.Direction.IsNull()) {
				plan.Direction = types.StringValue(v.(string))
			}
		case "clientAid":
			if !(plan.ClientAid.IsNull()) {
				plan.ClientAid = types.StringValue(v.(string))
			}
		case "lineAid":
			if !(plan.ClientAid.IsNull()) {
				plan.LineAid = types.StringValue(v.(string))
			}
		case "dscgAid":
			if !(plan.DscgAid.IsNull()) {
				plan.DscgAid = types.StringValue(v.(string))
			}
		case "remoteClientId":
			if len(v.(string)) > 0 {
				plan.RemoteClientId = types.StringValue(v.(string))
			}
		case "remoteModuleId":
			if len(v.(string)) > 0 {
				plan.RemoteModuleId = types.StringValue(v.(string))
			}
		}
	}
	tflog.Debug(ctx, "LCResource: read ## ", map[string]interface{}{"plan": plan})
}

func (r *LCResource) delete(plan *LCResourceData, ctx context.Context, diags *diag.Diagnostics) {

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		diags.AddError(
			"Error Delete LC",
			"Delete: Could not Delete LC, href is not specified",
		)
		return
	}

	tflog.Debug(ctx, "LCResource: delete ## ", map[string]interface{}{"href": href})

	_, _, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "DELETE", "resource-links"+href, nil)

	if err != nil && !strings.Contains(err.Error(), "status: 404") {
		diags.AddError(
			"Error Delete LC",
			"Delete: Could not Delete LC, unexpected error: "+err.Error(),
		)
		return
	}
}
