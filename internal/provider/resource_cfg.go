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
	_ resource.Resource                = &CfgResource{}
	_ resource.ResourceWithConfigure   = &CfgResource{}
	_ resource.ResourceWithImportState = &CfgResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewCfgResource() resource.Resource {
	return &CfgResource{}
}

type CfgResource struct {
	client *xrcm_pf.Client
}

type CfgResourceData struct {
	Id                 types.String `tfsdk:"id"`
	N                  types.String `tfsdk:"n"`
	DeviceId           types.String `tfsdk:"deviceid"`
	Aid                types.String `tfsdk:"aid"`
	ConfiguredRole     types.String `tfsdk:"configuredrole"`
	CurrentRole        types.String `tfsdk:"currentrole"`
	RoleStatus         types.String `tfsdk:"rolestatus"`
	SerdesRate         types.String `tfsdk:"serdesrate"`
	TrafficMode        types.String `tfsdk:"trafficmode"`
	TcMode             types.Bool   `tfsdk:"tcmode"`
	RestartAction      types.String `tfsdk:"restartaction"`
	Topology           types.String `tfsdk:"topology"`
	ConfigState        types.String `tfsdk:"configstate"`
	FactoryResetAction types.Bool   `tfsdk:"factoryresetaction"`
	HId                types.String `tfsdk:"hid"`
	HPortId            types.String `tfsdk:"hportid"`
}

// Metadata returns the resource type name.
func (r *CfgResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cfg"
}

// Schema defines the schema for the data source.
func (r *CfgResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cfg resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the XR Cfg.",
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
				Description: "aid",
				Computed:    true,
			},
			"configuredrole": schema.StringAttribute{
				Description: "configured role",
				Optional:    true,
			},
			"currentrole": schema.StringAttribute{
				Description: "current role",
				Computed:    true,
			},
			"rolestatus": schema.StringAttribute{
				Description: "role status",
				Computed:    true,
			},
			"trafficmode": schema.StringAttribute{
				Description: "traffic mode",
				Optional:    true,
			},
			"serdesrate": schema.StringAttribute{
				Description: "serdes rate",
				Computed:    true,
			},
			"tcmode": schema.BoolAttribute{
				Description: "TC Mode",
				Optional:    true,
			},
			"restartaction": schema.StringAttribute{
				Description: "restart action",
				Optional:    true,
			},
			"configstate": schema.StringAttribute{
				Description: "config State",
				Computed:    true,
			},
			"factoryresetaction": schema.BoolAttribute{
				Description: "factory Reset Action",
				Optional:    true,
			},
			"topology": schema.StringAttribute{
				Description: "topology",
				Optional:    true,
			},
			"hid": schema.StringAttribute{
				Description: "Host ID",
				Computed:    true,
			},
			"hportid": schema.StringAttribute{
				Description: "Host Port ID",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CfgResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r CfgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CfgResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "CfgResource: Create", map[string]interface{}{"CfgResourceData": data})
	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r CfgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CfgResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "CfgResource: Read", map[string]interface{}{"CfgResourceData": data})

	r.read(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() {
		resp.State = tfsdk.State{}
	} else {
		diags = resp.State.Set(ctx, &data)
	}

	resp.Diagnostics.Append(diags...)
}

func (r CfgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CfgResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "CfgResource: Update", map[string]interface{}{"CfgResourceData": data})
	// diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r CfgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CfgResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "CfgResource: Delete", map[string]interface{}{"CfgResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CfgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CfgResource) update(plan *CfgResourceData, ctx context.Context, diags *diag.Diagnostics) {

	// convert TF to Cfg json - required do to camel case not supported in TF
	tflog.Debug(ctx, "CfgResource: createUpdate ## ")
	var cmd = make(map[string]interface{})
	if !(plan.ConfiguredRole.IsNull()) {
		cmd["configuredRole"] = plan.ConfiguredRole.ValueString()
	}

	if !(plan.TrafficMode.IsNull()) {
		cmd["trafficMode"] = plan.TrafficMode.ValueString()
	}

	if !(plan.Topology.IsNull()) {
		cmd["topology"] = plan.Topology.ValueString()
	}
	if !(plan.N.IsNull()) {
		cmd["n"] = plan.N.ValueString()
	}

	if !(plan.TcMode.IsNull()) {
		cmd["tcMode"] = plan.TcMode.ValueBool()
	}

	if !(plan.RestartAction.IsNull()) {
		cmd["restartAction"] = plan.RestartAction.ValueString()
	}

	if !(plan.FactoryResetAction.IsNull()) {
		cmd["factoryResetAction"] = plan.FactoryResetAction.ValueBool()
	}

	rb, err := json.Marshal(cmd)

	tflog.Debug(ctx, "CfgResource: createUpdate ## ", map[string]interface{}{"data": string(rb)})

	if err != nil {
		diags.AddError(
			"CfgResource: createUpdate ##: Error creating cfg",
			"Could not create cfg, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/cfg"
	}

	body, deviceid, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources/cfg", rb)

	if err != nil {
		diags.AddError(
			"CfgResource: createUpdate ##: Error creating cfg",
			"Could not create cfg, unexpected error: "+err.Error(),
		)
		return
	}
	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	plan.DeviceId = types.StringValue(deviceid)
	if err != nil {
		diags.AddError(
			"CfgResource: update ##: Error update LC",
			"Update: Could not SetResourceId CfgResource, unexpected error: "+err.Error(),
		)
		return
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}

	if content["currentRole"] != nil {
		plan.CurrentRole = types.StringValue(content["currentRole"].(string))
	}

	if content["roleStatus"] != nil {
		plan.RoleStatus = types.StringValue(content["roleStatus"].(string))
	}

	if content["serdesRate"] != nil {
		plan.SerdesRate = types.StringValue(content["serdesRate"].(string))
	}

	if content["hId"] != nil {
		plan.HId = types.StringValue(content["hId"].(string))
	}

	if content["hportid"] != nil {
		plan.HPortId = types.StringValue(content["hportid"].(string))
	}

	if content["configState"] != nil {
		plan.ConfigState = types.StringValue(content["configState"].(string))
	}

	tflog.Debug(ctx, "CfgResource: createUpdate ## ", map[string]interface{}{"deviceid": deviceid})
}

func (r *CfgResource) read(plan *CfgResourceData, ctx context.Context, diags *diag.Diagnostics) {

	tflog.Debug(ctx, "CfgResource: read ", map[string]interface{}{"deviceID": plan.N.ValueString(), "URL": "resources/cfg"})

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/cfg"
	}

	body, deviceid, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"CfgResource: read ##: Error Get config",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringNull()
		tflog.Debug(ctx, "CfgResource: read - not found ## 404", map[string]interface{}{"plan": plan})
		return
	}

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"CfgResource: read ##: Error creating cfg",
			"Could not create cfg, unexpected error: "+err.Error(),
		)
		return
	}

	plan.DeviceId = types.StringValue(deviceid)

	for k, v := range content {
		switch k {
		case "aid":
			plan.Aid = types.StringValue(v.(string))
		case "currentRole":
			plan.CurrentRole = types.StringValue(v.(string))
		case "roleStatus":
			plan.RoleStatus = types.StringValue(v.(string))
		case "serdesRate":
			plan.SerdesRate = types.StringValue(v.(string))
		case "configuredRole":
			if !(plan.ConfiguredRole.IsNull()) {
				plan.ConfiguredRole = types.StringValue(v.(string))
			}
		case "trafficMode":
			if !(plan.TrafficMode.IsNull()) {
				plan.TrafficMode = types.StringValue(v.(string))
			}
		case "hId":
			plan.HId = types.StringValue(v.(string))
		case "hPortId":
			plan.HPortId = types.StringValue(v.(string))
		case "topology":
			if !(plan.Topology.IsNull()) {
				plan.Topology = types.StringValue(v.(string))
			}
		case "restartAction":
			if !(plan.RestartAction.IsNull()) {
				plan.RestartAction = types.StringValue(v.(string))
			}
		case "configState":
			if !(plan.ConfigState.IsNull()) {
				plan.ConfigState = types.StringValue(v.(string))
			}
		case "tcMode":
			if !(plan.TcMode.IsNull()) {
				plan.TcMode = types.BoolValue(v.(bool))
			}
		case "factoryResetAction":
			if !(plan.FactoryResetAction.IsNull()) {
				plan.FactoryResetAction = types.BoolValue(v.(bool))
			}
		}
	}
	tflog.Debug(ctx, "CfgResource: read ## ", map[string]interface{}{"Plan": plan})
}
