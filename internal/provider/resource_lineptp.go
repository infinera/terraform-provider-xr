package provider

import (
	"context"
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
	_ resource.Resource                = &LinePTPResource{}
	_ resource.ResourceWithConfigure   = &LinePTPResource{}
	_ resource.ResourceWithImportState = &LinePTPResource{}
)

// NewACResource is a helper function to simplify the provider implementation.
func NewLinePTPResource() resource.Resource {
	return &LinePTPResource{}
}

type LinePTPResource struct {
	client *xrcm_pf.Client
}

// Metadata returns the data source type name.
func (r *LinePTPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_LinePTP"
}

type LinePTPResourceData struct {
	Id        types.String `tfsdk:"id"`
	DeviceId  types.String `tfsdk:"deviceid"`
	N         types.String `tfsdk:"n"`
	Aid       types.String `tfsdk:"aid"`
	LinePTPId types.String `tfsdk:"lineptpid"`
}

// Schema defines the schema for the resource.
func (r *LinePTPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Line PTP",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the Line PTP.",
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
				Description: "Line PTP id",
				Optional:    true,
			},
			"aid": schema.StringAttribute{
				Description: "aid",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *LinePTPResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*xrcm_pf.Client)
}

func (r LinePTPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LinePTPResourceData

	diags := req.Config.Get(ctx, &data)
	tflog.Debug(ctx, "LinePTPResource: Create", map[string]interface{}{"LinePTPResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.read(&data, ctx, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r LinePTPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LinePTPResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "LinePTPResource: Read", map[string]interface{}{"LinePTPResourceData": data})

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

func (r LinePTPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LinePTPResourceData

	diags := req.Plan.Get(ctx, &data)
	tflog.Debug(ctx, "LinePTPResource: Update", map[string]interface{}{"LinePTPResourceData": data})
	// diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.read(&data, ctx, &resp.Diagnostics)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r LinePTPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LinePTPResourceData

	diags := req.State.Get(ctx, &data)
	tflog.Debug(ctx, "LinePTPResource: Delete", map[string]interface{}{"LinePTPResourceData": data})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *LinePTPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve imLinePTP ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *LinePTPResource) read(plan *LinePTPResourceData, ctx context.Context, diags *diag.Diagnostics) {

	if plan.LinePTPId.IsNull() {
		diags.AddError(
			"LinePTPResource: read ##: Error Read LinePTP",
			"Read: Could not get LinePTP, LinePTP ID  are not specified ",
		)
		return
	}

	href := after(plan.Id.ValueString(), "/")
	if len(href) == 0 {
		href = "/lineptps/" + plan.LinePTPId.ValueString()
	}

	tflog.Debug(ctx, "LinePTPResource: read ## ", map[string]interface{}{"device": plan.N.ValueString(), "URL": "resources" + href})

	body, deviceId, err := r.client.ExecuteDeviceHttpCommand(plan.N.ValueString(), "GET", "resources"+href, nil)

	if err != nil {
		if !strings.Contains(err.Error(), "status: 404") {
			diags.AddError(
				"LinePTPResource: read ##: Error Get AC",
				"Read: Could not Get , unexpected error: "+err.Error(),
			)
			return
		}
		plan.Id = types.StringValue("")
		tflog.Debug(ctx, "LinePTPResource: read - not found ## 404", map[string]interface{}{"plan": plan})
		return
	}

	tflog.Debug(ctx, "LinePTPResource: read ## ", map[string]interface{}{"response": string(body)})

	plan.DeviceId = types.StringValue(deviceId)

	content, err := SetResourceId(plan.N.ValueString(), &plan.Id, body)

	if err != nil {
		diags.AddError(
			"LinePTPResource: read ##: Error Read LinePTP",
			"Create: Could not SetResourceId , unexpected error: "+err.Error(),
		)
		return
	}
	if content["aid"] != nil {
		plan.Aid = types.StringValue(content["aid"].(string))
	}
	tflog.Debug(ctx, "LinePTPResource: read ## ", map[string]interface{}{"plan": plan})
}
