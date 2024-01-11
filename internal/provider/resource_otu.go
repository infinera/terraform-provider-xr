package provider

import (
	"context"
	"encoding/json"

	//"strings"

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
	_ resource.Resource                = &OTUResource{}
	_ resource.ResourceWithConfigure   = &OTUResource{}
	_ resource.ResourceWithImportState = &OTUResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewOTUResource() resource.Resource {
	return &OTUResource{}
}

type OTUResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *OTUResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_otu"
}

type OTUResourceData struct {
	Id          types.String `tfsdk:"id"`
	N           types.String `tfsdk:"n"`
	DeviceId    types.String `tfsdk:"deviceid"`
	Aid         types.String `tfsdk:"aid"`
	Otutype     types.String `tfsdk:"otutype"`
	Rate        types.Int64  `tfsdk:"rate"`
	OtuId       types.String `tfsdk:"otuid"`
	RxTTI       types.String `tfsdk:"rxtti"`
	TxTTI       types.String `tfsdk:"txtti"`
	ExpectedTTI types.String `tfsdk:"expectedtti"`
}

// Schema defines the schema for the resource.
func (r *OTUResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Carrier",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the OTU.",
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
			"otuid": schema.StringAttribute{
				Description: "otu id",
				Optional:    true,
			},
			"aid": schema.StringAttribute{
				Description: "aid",
				Computed:    true,
			},
			"otutype": schema.StringAttribute{
				Description: "OTu type",
				Computed:    true,
			},
			"rate": schema.Int64Attribute{
				Description: "rate",
				Optional:    true,
			},
			"rxtti": schema.StringAttribute{
				Description: "rx tti",
				Computed:    true,
			},
			"txtti": schema.StringAttribute{
				Description: "tx tti",
				Optional:    true,
			},
			"expectedtti": schema.StringAttribute{
				Description: "expected tti",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *OTUResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r OTUResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OTUResourceData

	diags := req.Config.Get(ctx, &data)
	tflog.Debug(ctx, "OTUResource: Create", map[string]interface{}{"OTUResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r OTUResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OTUResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "OTUResource: Read", map[string]interface{}{"OTUResourceData": data})

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

func (r OTUResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OTUResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "OTUResource: Update", map[string]interface{}{"OTUResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r OTUResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OTUResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "OTUResource: Delete", map[string]interface{}{"OTUResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *OTUResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *OTUResource) update(plan *OTUResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.OtuId.IsNull() {
		diags.AddError(
			"OTUResource: update ##: Error Update OTU",
			"Update: Could not Update OTU, ID is not specified ",
		)
		return
	}

	tflog.Debug(ctx, "OTUResource: update ## ", map[string]interface{}{"OtuId": plan.OtuId.ValueString()})

	var cmd = make(map[string]interface{})

	if !(plan.TxTTI.IsNull()) {
		cmd["txTTI"] = plan.TxTTI.ValueString()
	}

	if !(plan.Rate.IsNull()) {
		cmd["rate"] = plan.Rate.ValueInt64()
	}

	if !(plan.ExpectedTTI.IsNull()) {
		cmd["expectedTTI"] = plan.ExpectedTTI.ValueString()
	}

	if len(cmd) == 0. {
		tflog.Debug(ctx, "OTUResource: update ## No Settings, Nothing to configure", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources/otus/" + plan.OtuId.ValueString()})
		return
	}

	rb, err := json.Marshal(cmd)

	if err != nil {
		diags.AddError(
			"OTUResource: update ##: Error Update OTU",
			"Update: Could not Update OTU, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/otus/" + plan.OtuId.ValueString()
	}

	tflog.Debug(ctx, "OTUResource: update ## ", map[string]interface{}{"Device": plan.N.ValueString(), "URL": "resources" + href, "Input data": string(rb)})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"OTUResource: update ##: Error Update OTU",
			"Update: Could not Update OTU, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "OTUResource: update ## ExecuteDeviceHttpCommand ..", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"OTUResource: update ##: Error Read OTU",
			"Create: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}

	if content["rate"] != nil {
		plan.Rate = types.Int64Value(content["rate"].(int64))
	}

	if content["otutype"] != nil {
		plan.Otutype = types.StringValue(content["otutype"].(string))
	}

	if content["rxTTI"] != nil {
		plan.RxTTI = types.StringValue(content["rxTTI"].(string))
	}

	tflog.Debug(ctx, "OTUResource: update ## ", map[string]interface{}{"plan": plan})
}

func (r *OTUResource) read(plan *OTUResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.OtuId.IsNull() {
		diags.AddError(
			"OTUResource: read ##: Error Read OTU",
			"Read: Could not Create OTU, ID  are not specified ",
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/otus/ " + plan.OtuId.ValueString()
	}

	tflog.Debug(ctx, "OTUResource: read ## ", map[string]interface{}{"device": plan.N.ValueString(), "URL": href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		/*if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"OTUResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "OTUResource: read - not found ## 404", map[string]interface{}{"href": href})
		return*/
		diags.AddError(
			"ODUResource: read ##: Error Get ODU",
			"Read: Could not Get , unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "OTUResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"OTUResource: read ##: Error Read OTU",
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
		case "otutype":
			if len(v.(string)) > 0 {
				plan.Otutype = types.StringValue(v.(string))
			}
		case "rxTTI":
			if !(plan.RxTTI.IsNull()) {
				plan.RxTTI = types.StringValue(v.(string))
			}
		case "txTTI":
			if !(plan.TxTTI.IsNull()) {
				plan.TxTTI = types.StringValue(v.(string))
			}
		case "expectedTTI":
			if !(plan.ExpectedTTI.IsNull()) {
				plan.ExpectedTTI = types.StringValue(v.(string))
			}
		case "rate":
			if !(plan.Rate.IsNull()) {
				plan.Rate = types.Int64Value(v.(int64))
			}
		}
	}
	tflog.Debug(ctx, "OTUResource: read ## ", map[string]interface{}{"plan": plan})
}
