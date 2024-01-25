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
	_ resource.Resource                = &EthernetResource{}
	_ resource.ResourceWithConfigure   = &EthernetResource{}
	_ resource.ResourceWithImportState = &EthernetResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewEthernetResource() resource.Resource {
	return &EthernetResource{}
}

type EthernetResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *EthernetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ethernet"
}

type EthernetResourceData struct {
	Id         types.String `tfsdk:"id"`
	N          types.String `tfsdk:"n"`
	DeviceId   types.String `tfsdk:"deviceid"`
	Aid        types.String `tfsdk:"aid"`
	EthernetId types.String `tfsdk:"ethernetid"`
	FecMode    types.String `tfsdk:"fecmode"`
	FecType    types.String `tfsdk:"fectype"`
	PortSpeed  types.Int64  `tfsdk:"portspeed"`
	MaxPktLen  types.Int64  `tfsdk:"maxpktlen"`
	ConfigState    types.String `tfsdk:"configstate"`
}

// Schema defines the schema for the  resource.
func (r *EthernetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Ethernet",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the Ethernet.",
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
			"ethernetid": schema.StringAttribute{
				Description: "ethernet id",
				Optional:    true,
			},
			"aid": schema.StringAttribute{
				Description: "AID",
				Computed:    true,
			},
			"fecmode": schema.StringAttribute{
				Description: "fec mode",
				Optional:    true,
			},
			"fectype": schema.StringAttribute{
				Description: "fec type",
				Computed:    true,
			},
			"portspeed": schema.Int64Attribute{
				Description: "fec mode",
				Computed:    true,
			},
			"maxpktlen": schema.Int64Attribute{
				Description: "maxpktlen",
				Optional:    true,
			},
			"configstate": schema.StringAttribute{
				Description: "configstate",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *EthernetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r EthernetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EthernetResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetResource: Create", map[string]interface{}{"EthernetResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//r.create(&data, ctx, &resp.Diagnostics)
	r.update(&data, ctx, &resp.Diagnostics)

	//	data.Id = types.String{Value: data.N.Value}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r EthernetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EthernetResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetResource: Read", map[string]interface{}{"EthernetResourceData": data})

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

func (r EthernetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EthernetResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "EthernetResource: Update", map[string]interface{}{"EthernetResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r EthernetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EthernetResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "EthernetResource: Delete", map[string]interface{}{"EthernetResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *EthernetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *EthernetResource) update(plan *EthernetResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.EthernetId.IsNull() {
		diags.AddError(
			"EthernetResource: update ##: Error Update Ethernet",
			"Update: Could not Update Ethernet, Ethernet ID  are not specified ",
		)
		return
	}

	tflog.Debug(ctx, "EthernetResource: update ## ", map[string]interface{}{"EthernetId": plan.EthernetId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.FecMode.IsNull()) {
		cmd["fecMode"] = plan.FecMode.ValueString()
	}
	if !(plan.MaxPktLen.IsNull()) {
		cmd["maxPktLen"] = plan.MaxPktLen.ValueInt64()
	}

	if len(cmd) == 0. {
		tflog.Debug(ctx, "EthernetResource: update ## No Settings, Nothing to configure", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources/ethernets/" + plan.EthernetId.ValueString()})
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"EthernetResource: update ##: Error creating Ethernet",
			"Update: Could not create Ethernet, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + plan.EthernetId.ValueString()
	}

	tflog.Debug(ctx, "EthernetResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"EthernetResource: update ##: Error Update Ethernet",
			"Update: Could not create Ethernet, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "EthernetResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"EthernetResource: update ##: Error Read Ethernet",
			"Update: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	if content["fecType"] != nil {
		plan.FecType = types.StringValue(content["fecType"].(string))
	}
	if content["maxPktLen"] != nil {
		plan.MaxPktLen = types.Int64Value(int64(content["maxPktLen"].(float64)))
	}
	tflog.Debug(ctx, "EthernetResource: update ## ", map[string]interface{}{"plan": plan})

}

func (r *EthernetResource) read(plan *EthernetResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.EthernetId.IsNull() {
		diags.AddError(
			"EthernetResource: read ##: Error Read Ethernet",
			"Read: Could not Update Ethernet, Ethernet ID  are not specified ",
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + plan.EthernetId.ValueString()
	}

	tflog.Debug(ctx, "EthernetResource: read ##  ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources/" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"EthernetResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "EthernetResource: read - not found ## 404", map[string]interface{}{"plan": plan})
		return
	}

	tflog.Debug(ctx, "EthernetResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"EthernetResource: read ##: Error Create Ethernet",
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
		case "fecType":
			if len(v.(string)) > 0 {
				plan.FecType = types.StringValue(v.(string))
			}
		case "fecMode":
			if !(plan.FecMode.IsNull()) {
				plan.FecMode = types.StringValue(v.(string))
			}
		case "portSpeed":
			plan.PortSpeed = types.Int64Value(int64(v.(float64)))
		case "maxPktLen":
			if !(plan.MaxPktLen.IsNull()) {
				plan.MaxPktLen = types.Int64Value(int64(v.(float64)))
			}
		case "configState":
			if len(v.(string)) > 0 {
				plan.ConfigState = types.StringValue(v.(string))
			}
		}

	}
	tflog.Debug(ctx, "EthernetResource: read ## ", map[string]interface{}{"plan": plan})
}
