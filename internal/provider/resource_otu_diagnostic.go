package provider

import (
	"context"
	"encoding/json"

	//"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	//"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &OTUDiagResource{}
	_ resource.ResourceWithConfigure   = &OTUDiagResource{}
	_ resource.ResourceWithImportState = &OTUDiagResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewOTUDiagResource() resource.Resource {
	return &OTUDiagResource{}
}

type OTUDiagResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *OTUDiagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_OTUDiag"
}

type OTUDiagResourceData struct {
	Id         types.String `tfsdk:"id"`
	N          types.String `tfsdk:"n"`
	DeviceId   types.String `tfsdk:"deviceid"`
	Aid        types.String `tfsdk:"aid"`
	OtuId      types.String `tfsdk:"otuid"`
	OTUId      types.String `tfsdk:"OTUid"`
	TermLB     types.String `tfsdk:"termlb"`
	facPRBSGen types.Bool   `tfsdk:"facprbsgen"`
	facPRBSMon types.Bool   `tfsdk:"facprbsmon"`
}

// Schema defines the schema for the resource.
func (r *OTUDiagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Carrier",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the OTUDiag.",
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
			"OTUid": schema.StringAttribute{
				Description: "OTU id",
				Optional:    true,
			},
			"aid": schema.StringAttribute{
				Description: "aid",
				Computed:    true,
			},
			"termlb": schema.StringAttribute{
				Description: "Term Loopback",
				Optional:    true,
			},
			"facprbsGen": schema.BoolAttribute{
				Description: "fac prbs Gen",
				Optional:    true,
			},
			"facprbsmon": schema.BoolAttribute{
				Description: "fac prbs MOn",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *OTUDiagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r OTUDiagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OTUDiagResourceData

	diags := req.Config.Get(ctx, &data)
	tflog.Debug(ctx, "OTUDiagResource: Create", map[string]interface{}{"OTUDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//r.update(&data, ctx, &resp.Diagnostics)
	r.update(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r OTUDiagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OTUDiagResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "OTUDiagResource: Read", map[string]interface{}{"OTUDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.read(&data, ctx, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &data)

	resp.Diagnostics.Append(diags...)
}

func (r OTUDiagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OTUDiagResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "OTUDiagResource: Update", map[string]interface{}{"OTUDiagResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.update(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r OTUDiagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OTUDiagResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "OTUDiagResource: Delete", map[string]interface{}{"OTUDiagResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *OTUDiagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *OTUDiagResource) update(plan *OTUDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.OtuId.IsNull() {
		diags.AddError(
			"OTUDiagResource: updated ##: Error Read OTUDiag",
			"Read: Could not updated OTUDiag, OTU ID and OTU ID are not specified",
		)
		return
	}

	var cmd = make(map[string]interface{})
	if !(plan.facPRBSGen.IsNull()) {
		cmd["facPRBSGen"] = plan.facPRBSGen.ValueBool()
	}

	if !(plan.facPRBSGen.IsNull()) {
		cmd["facPRBSMon"] = plan.facPRBSGen.ValueBool()
	}

	if !(plan.TermLB.IsNull()) {
		cmd["termLB"] = plan.TermLB.ValueString()
	}

	if len(cmd) == 0. {
		return
	}

	rb, err := json.Marshal(cmd)
	if err != nil {
		diags.AddError(
			"DSCDiagResource: update ##: Error Update DSC Diagnostic",
			"Update: Could not Update DSC, unexpected error: "+err.Error(),
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/otus/" + plan.OtuId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "OTUDiagResource: updated ## ", map[string]interface{}{"device": plan.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "PUT", "resources"+href, rb)

	if err != nil {
		diags.AddError(
			"OTUDiagResource: read ##: Error Get OTUDiag",
			"Read: Could not Get , unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "OTUDiagResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	_, err = SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"OTUDiagResource: update ##: Error Update DSC Diagnostic",
			"Update: Could not SetResourceId OTU, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "OTUDiagResource: update ## ", map[string]interface{}{"plan": plan})

}

func (r *OTUDiagResource) read(plan *OTUDiagResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.OtuId.IsNull() {
		diags.AddError(
			"OTUDiagResource: read ##: Error Read OTUDiag",
			"Read: Could not Read OTUDiag, OTU ID and OTUDiag ID are not specified",
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/otus/" + plan.OtuId.ValueString() + "/diagnostic"
	}

	tflog.Debug(ctx, "OTUDiagResource: read ## ", map[string]interface{}{"device": plan.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		diags.AddError(
			"OTUDiagResource: read ##: Error Get OTUDiag",
			"Read: Could not Get , unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "OTUDiagResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)
	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)
	if err != nil {
		diags.AddError(
			"OTUDiagResource: read ##: Error Read OTUDiag",
			"Read: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}

	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}
	if content["termLB"] != nil {
		plan.TermLB = types.StringValue(content["termLB"].(string))
	}
	if content["FacPRBSGen"] != nil {
		plan.facPRBSGen = types.BoolValue(content["facPRBSGen"].(bool))
	}
	if content["FacPRBSMon"] != nil {
		plan.facPRBSMon = types.BoolValue(content["facPRBSMon"].(bool))
	}

	tflog.Debug(ctx, "OTUDiagResource: read ## ", map[string]interface{}{"plan": plan})
}
