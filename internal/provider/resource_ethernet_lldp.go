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
	_ resource.Resource                = &EthernetLLDPResource{}
	_ resource.ResourceWithConfigure   = &EthernetLLDPResource{}
	_ resource.ResourceWithImportState = &EthernetLLDPResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewEthernetLLDPResource() resource.Resource {
	return &EthernetLLDPResource{}
}

type EthernetLLDPResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *EthernetLLDPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ethernet_lldp"
}

type EthernetLLDPResourceData struct {
	Id               types.String `tfsdk:"id"`
	N                types.String `tfsdk:"n"`
	DeviceId         types.String `tfsdk:"deviceid"`
	Aid              types.String `tfsdk:"aid"`
	EthernetId       types.String `tfsdk:"ethernetid"`
	AdminStatus      types.String `tfsdk:"adminstatus"`
	GccFwd           types.Bool   `tfsdk:"gccfwd"`
	HostRxDrop       types.Bool   `tfsdk:"hostrxdrop"`
	TTLUsage         types.Bool   `tfsdk:"ttlusage"`
	ClrStats         types.Bool   `tfsdk:"clrstats"`
	FlushHostDb      types.Bool   `tfsdk:"flushhostdb"`
	TooManyNeighbors types.Bool   `tfsdk:"toomanyneighbors"`
}

// Schema defines the schema for the  resource.
func (r *EthernetLLDPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an EthernetLLDP",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the EthernetLLDP.",
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
			"adminstatus": schema.StringAttribute{
				Description: "admin status",
				Optional:    true,
			},
			"gccfwd": schema.BoolAttribute{
				Description: "gcc fwd",
				Optional:    true,
			},
			"hostrxdrop": schema.BoolAttribute{
				Description: "host rx drop",
				Optional:    true,
			},
			"ttlusage": schema.BoolAttribute{
				Description: "ttl usage",
				Optional:    true,
			},
			"clrstats": schema.BoolAttribute{
				Description: "clr stats",
				Optional:    true,
			},
			"flushhostdb": schema.BoolAttribute{
				Description: "flush host db",
				Optional:    true,
			},
			"toomanyneighbors": schema.BoolAttribute{
				Description: "too many neighbors",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the  resource.
func (r *EthernetLLDPResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r EthernetLLDPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EthernetLLDPResourceData
	diags := req.Config.Get(ctx, &data)
	tflog.Debug(ctx, "EthernetLLDPResource: Create", map[string]interface{}{"EthernetLLDPResourceData": data})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r EthernetLLDPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EthernetLLDPResourceData
	diags := req.State.Get(ctx, &data)
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

func (r EthernetLLDPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EthernetLLDPResourceData

	diags := req.Plan.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetLLDPResource: Update", map[string]interface{}{"EthernetLLDPResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r EthernetLLDPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EthernetLLDPResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetLLDPResource: Delete", map[string]interface{}{"EthernetLLDPResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *EthernetLLDPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *EthernetLLDPResource) read(state *EthernetLLDPResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if state.EthernetId.IsNull() {
		diags.AddError(
			"EthernetLLDPResource: read ##: Error Read Ethernet LLDP",
			"Read: Could not Read Ethernet LLDP, Ethernet ID is not specified ",
		)
		return
	}

	href := after(state.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + state.EthernetId.ValueString() + "/lldp-cfg"
	}

	tflog.Debug(ctx, "EthernetLLDPResource: read ## ", map[string]interface{}{"Device": state.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(state.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"EthernetLLDPResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		state.Id = types.StringValue("")
		tflog.Debug(ctx, "EthernetLLDPResource: read - not found ## 404", map[string]interface{}{"state": state})
		return
	}

	tflog.Debug(ctx, "EthernetLLDPResource: read ## ", map[string]interface{}{"response": string(body)})

	if err != nil {
		diags.AddError(
			"EthernetLLDPResource: read ##: Error Read Ethernet LLDP",
			"Could not Unmarshal Ethernet LLDP, unexpected error: "+err.Error(),
		)
		return
	}

	state.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(state.N.ValueString(), &state.Id, body)
	if err != nil {
		diags.AddError(
			"EthernetLLDPResource: read ##: Error Read Ethernet",
			"Read: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	for k, v := range content {
		switch k {
		case "aid":
			if len(v.(string)) > 0 {
				state.Aid = types.StringValue(v.(string))
			}
		case "tooManyNeighbors":
			if v != nil {
				state.TooManyNeighbors = types.BoolValue(v.(bool))
			}
		case "adminStatus":
			if !(state.AdminStatus.IsNull()) {
				state.AdminStatus = types.StringValue(v.(string))
			}
		case "gccFwd":
			if !(state.GccFwd.IsNull()) {
				state.GccFwd = types.BoolValue(v.(bool))
			}
		case "hostRxDrop":
			if !(state.HostRxDrop.IsNull()) {
				state.HostRxDrop = types.BoolValue(v.(bool))
			}
		case "TTLUsage":
			if !(state.TTLUsage.IsNull()) {
				state.TTLUsage = types.BoolValue(v.(bool))
			}
		case "clrStats":
			if !(state.ClrStats.IsNull()) {
				state.ClrStats = types.BoolValue(v.(bool))
			}
		case "flushHostDb":
			if !(state.FlushHostDb.IsNull()) {
				state.FlushHostDb = types.BoolValue(v.(bool))
			}
		}
	}

	tflog.Debug(ctx, "EthernetLLDPResource: read ## ", map[string]interface{}{"state": state})
}

func (r *EthernetLLDPResource) update(plan *EthernetLLDPResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.EthernetId.IsNull() {
		diags.AddError(
			"EthernetLLDPResource: update ##: Error Update Ethernet LLDP",
			"Update: Could not Update Ethernet LLDP, Ethernet ID is not specified ",
		)
		return
	}

	tflog.Debug(ctx, "EthernetLLDPResource: update ## ", map[string]interface{}{"EthernetId": plan.EthernetId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.AdminStatus.IsNull()) {
		cmd["adminStatus"] = plan.AdminStatus.ValueString()
	}

	if !(plan.GccFwd.IsNull()) {
		cmd["gccFwd"] = plan.GccFwd.ValueBool()
	}

	if !(plan.HostRxDrop.IsNull()) {
		cmd["hostRxDrop"] = plan.HostRxDrop.ValueBool()
	}

	if !(plan.TTLUsage.IsNull()) {
		cmd["TTLUsage"] = plan.TTLUsage.ValueBool()
	}

	if !(plan.ClrStats.IsNull()) {
		cmd["clrStats"] = plan.ClrStats.ValueBool()
	}

	if !(plan.FlushHostDb.IsNull()) {
		cmd["flushHostDb"] = plan.FlushHostDb.ValueBool()
	}
	if len(cmd) == 0. {
		tflog.Debug(ctx, "EthernetLLDPResource: update ## No Settings, Nothing to configure", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources/ethernets/" + plan.EthernetId.ValueString() + "/lldp-cfg"})
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"EthernetLLDPResource: update ##: Error Update EthernetLLDP",
			"Update: Could not Update EthernetLLDP, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + plan.EthernetId.ValueString() + "/lldp-cfg"
	}

	tflog.Debug(ctx, "EthernetLLDPResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"EthernetLLDPResource: update ##: Error Update EthernetLLDP",
			"Update: Could not Update EthernetLLDP, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "EthernetLLDPResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"EthernetLLDPResource: update ##: Error Update Ethernet",
			"Update: Could not GetSetData EthernetLLDP, unexpected error: "+err.Error(),
		)
		return
	}

	plan.DeviceId = types.StringValue(deviceId)
	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}

	if content["tooManyNeighbors"] != nil {
		plan.TooManyNeighbors = types.BoolValue(content["tooManyNeighbors"].(bool))
	}

	tflog.Debug(ctx, "EthernetLLDPResource: update ## ", map[string]interface{}{"plan": plan})

}
