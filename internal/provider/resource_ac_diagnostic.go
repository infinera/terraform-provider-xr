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
	_ resource.Resource                = &ACDiagResource{}
	_ resource.ResourceWithConfigure   = &ACDiagResource{}
	_ resource.ResourceWithImportState = &ACDiagResource{}
)

// NewACDiagResource is a helper function to simplify the provider implementation.
func NewACDiagResource() resource.Resource {
	return &ACDiagResource{}
}

type ACDiagResource struct {
	client *xrcm_pf.Client
}

type ACDiagResourceData struct {
	Id         types.String `tfsdk:"id"`
	N          types.String `tfsdk:"n"`
	DeviceId   types.String `tfsdk:"deviceid"`
	EthernetId types.String `tfsdk:"ethernetid"`
	AcId       types.String `tfsdk:"acid"`
	Aid        types.String `tfsdk:"aid"`
	TermLB     types.String `tfsdk:"termlb"`
}

// Metadata returns the data source type name.
func (r *ACDiagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ac"
}

// Schema defines the schema for the data source.
func (r *ACDiagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"termlb": schema.StringAttribute{
				Description: "term Loopback",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ACDiagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r ACDiagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ACDiagResourceData

	diags := req.Config.Get(ctx, &data)

	tflog.Debug(ctx, "ACDiagResource: Create - ", map[string]interface{}{"ACDiagResourceData": data})

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

func (r ACDiagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ACDiagResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "ACDiagResource: Create - ", map[string]interface{}{"ACDiagResourceData": data})

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.read(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ACDiagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ACDiagResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "CfgResource: Update", map[string]interface{}{"ACDiagResourceData": data})

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r ACDiagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ACDiagResourceData

	diags := req.State.Get(ctx, &data)

	tflog.Debug(ctx, "CfgResource: Update", map[string]interface{}{"ACDiagResourceData": data})

	resp.Diagnostics.Append(diags...)

	r.delete(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ACDiagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ACDiagResource) read(plan *ACDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.AcId.IsNull() || plan.EthernetId.IsNull() {
		diags.AddError(
			"Error get AC",
			"Create: Could not get AC, AC ID or Ethernet Id is not specified",
		)
		return
	}

	tflog.Debug(ctx, "ACDiagResource: read ## ", map[string]interface{}{"Acid": plan.AcId.ValueString(), "EthernetId": plan.EthernetId.ValueString()})

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "ACDiagResource: read - href is empty", map[string]interface{}{"plan": plan})
		return
	}

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"ACDiagResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "ACDiagResource: read - not found ## 404", map[string]interface{}{"plan": plan})
		return
	}

	tflog.Debug(ctx, "ACDiagResource: read ## ", map[string]interface{}{"response": string(body)})

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"LCResource: read ##: Error Read LC",
			"Create: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	plan.DeviceId = types.StringValue(deviceId)

	if content["termLB"] != nil {

	}

	for k, v := range content {
		switch k {
		case "aid":
			{
				plan.Aid = types.StringValue(v.(string))
			}
		case "termLB":
			if !(plan.TermLB.IsNull()) {
				plan.TermLB = types.StringValue(v.(string))
			}
		}
	}
	tflog.Debug(ctx, "ACDiagResource: read ## ", map[string]interface{}{"plan": plan})
}

func (r *ACDiagResource) update(plan *ACDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.AcId.IsNull() || plan.EthernetId.IsNull() {
		diags.AddError(
			"Error Update AC",
			"Create: Could not Update AC, AC ID or Ethernet Id is not specified",
		)
		return
	}

	tflog.Debug(ctx, "ACDiagResource: create ## ", map[string]interface{}{"Acid": plan.AcId.ValueString(), "EthernetId": plan.EthernetId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.TermLB.IsNull()) {
		cmd["termLB"] = plan.TermLB.ValueString()
	}

	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"ACDiagResource: update ##: Error creating ACDiagResource",
			"Update: not create ACDiagResource, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/ethernets/" + plan.EthernetId.ValueString() + "/acs/" + plan.AcId.ValueString()
	}

	tflog.Debug(ctx, "ACDiagResource: Update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resource-links" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"ACDiagResource: update ##: Error Update Carrier",
			"Update:Could not Update, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "ACDiagResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"ACDiagResource: update ##: Error update LC",
			"Update: Could not SetResourceId ACDiagResource, unexpected error: "+err.Error(),
		)
		return
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}

	tflog.Debug(ctx, "ACDiagResource: update ## ", map[string]interface{}{"plan": plan})
}

func (r *ACDiagResource) delete(plan *ACDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		return
	}

	tflog.Debug(ctx, "ACDiagResource: delete ## ", map[string]interface{}{"href": href})

	body, _, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "DELETE", "resource-links"+href, nil)

	if err != nil && !strings.Contains(err.Error(), "status: 404") {
		diags.AddError(
			"Error Delete LC",
			"Delete: Could not Delete LC, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "ACDiagResource: delete ##  ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

}
