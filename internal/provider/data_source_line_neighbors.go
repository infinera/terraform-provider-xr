package provider

import (
	"context"
	"encoding/json"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LineNeighborDataSource{}
	_ datasource.DataSourceWithConfigure = &LineNeighborDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewLineNeighborDataSource() datasource.DataSource {
	return &LineNeighborDataSource{}
}

// coffeesDataSource is the data source implementation.
type LineNeighborDataSource struct {
	client *xrcm_pf.Client
}

type DiscoveredNeighborData struct {
	MacAddress             types.String `tfsdk:"macaddress"`
	CurrentRole            types.String `tfsdk:"currentrole"`
	DiscoveredTime         types.String `tfsdk:"discoveredtime"`
	ConstellationFrequency types.String `tfsdk:"constellationfrequency"`
}

type ControlPlaneNeighborData struct {
	MacAddress             types.String `tfsdk:"macaddress"`
	CurrentRole            types.String `tfsdk:"currentrole"`
	ConstellationFrequency types.String `tfsdk:"constellationfrequency"`
	ConState               types.String `tfsdk:"constate"`
	LastConStateChange     types.String `tfsdk:"lastconstatechange"`
}

type LineNeighborDataSourceData struct {
	N                     types.String               `tfsdk:"n"`
	DeviceId              types.String               `tfsdk:"deviceid"`
	LinePTPId             types.String               `tfsdk:"lineptpid"`
	DiscoveredNeighbors   []DiscoveredNeighborData   `tfsdk:"discoveredneighbors"`
	ControlPlaneNeighbors []ControlPlaneNeighborData `tfsdk:"controlplaneneighbors"`
}

// Metadata returns the data source type name.
func (d *LineNeighborDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_line_neighbors"
}

func (d *LineNeighborDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of neighbors.",
		Attributes: map[string]schema.Attribute{
			"deviceid": schema.StringAttribute{
				Description: "Device ID",
				Optional:    true,
			},
			"n": schema.StringAttribute{
				Description: "Device Name",
				Required:    true,
			},
			"lineptpid": schema.StringAttribute{
				Description: "line ptp id",
				Required:    true,
			},
			"discoveredneighbors": schema.ListNestedAttribute{
				Description: "List of discovered neighbors",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"macaddress": schema.StringAttribute{
							Description: "mac Address",
							Computed:    true,
						},
						"currentrole": schema.StringAttribute{
							Description: "current role",
							Computed:    true,
						},
						"constellationfrequency": schema.StringAttribute{
							Description: "constellation frequency",
							Computed:    true,
						},
						"discoveredtime": schema.StringAttribute{
							Description: "discovered time",
							Computed:    true,
						},
					},
				},
			},
			"controlplaneneighbors": schema.ListNestedAttribute{
				Description: "List of Control Plane neighbors",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"macaddress": schema.StringAttribute{
							Description: "mac Address",
							Computed:    true,
						},
						"currentrole": schema.StringAttribute{
							Description: "current role",
							Computed:    true,
						},
						"constellationfrequency": schema.StringAttribute{
							Description: "constellation frequency",
							Computed:    true,
						},
						"constate": schema.StringAttribute{
							Description: "connection state",
							Computed:    true,
						},
						"lastconstatechange": schema.StringAttribute{
							Description: "Last connection state",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LineNeighborDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *LineNeighborDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	queryData := LineNeighborDataSourceData{}
	diags := req.Config.Get(ctx, &queryData)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "LineNeighborDataSource Read: get LineNeighbor", map[string]interface{}{"req": req})

	body, deviceId, err := d.client.ExecuteDeviceHttpCommand(queryData.N.ValueString(), "GET", "resources/lineptps/"+queryData.LinePTPId.ValueString()+"/neighbors", nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"LineNeighborDataSource Read: Error Get LineNeighbor",
			"Read: Could not GET LineNeighbor, unexpected error: "+err.Error(),
		)
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"LineNeighborDataSource Read: Error Read LineNeighbor",
			"Read: Could not Unmarshal LineNeighbor, unexpected error: "+err.Error(),
		)
		return
	}

	resultData := data["data"].(map[string]interface{})
	content := resultData["content"].(map[string]interface{})
	tflog.Debug(ctx, "Read: get DSCGs", map[string]interface{}{"content": content})
	var discoveredneighbors []interface{}
	if content["discoveredneighbors"] != nil {
		discoveredneighbors = content["discoveredneighbors"].([]interface{})
	}
	var controlplaneneighbors []interface{}
	if content["controlplaneneighbors"] != nil {
		controlplaneneighbors = content["controlplaneneighbors"].([]interface{})
	}

	queryData.DiscoveredNeighbors = make([]DiscoveredNeighborData, 0)
	queryData.ControlPlaneNeighbors = make([]ControlPlaneNeighborData, 0)
	queryData.DeviceId = types.StringValue(deviceId)

	for _, v := range discoveredneighbors {
		neighbor := v.(map[string]interface{})
		neighborD := DiscoveredNeighborData{}
		neighborD.MacAddress = types.StringValue(neighbor["macAddress"].(string))
		neighborD.CurrentRole = types.StringValue(neighbor["currentRole"].(string))
		neighborD.DiscoveredTime = types.StringValue(neighbor["discoveredTime"].(string))
		neighborD.ConstellationFrequency = types.StringValue(neighbor["constellationFrequency"].(string))
		queryData.DiscoveredNeighbors = append(queryData.DiscoveredNeighbors, neighborD)
	}

	for _, v := range controlplaneneighbors {
		neighbor := v.(map[string]interface{})
		neighborC := ControlPlaneNeighborData{}
		neighborC.MacAddress = types.StringValue(neighbor["macAddress"].(string))
		neighborC.CurrentRole = types.StringValue(neighbor["currentRole"].(string))
		neighborC.ConstellationFrequency = types.StringValue(neighbor["constellationFrequency"].(string))
		neighborC.ConState = types.StringValue(neighbor["conState"].(string))
		neighborC.LastConStateChange = types.StringValue(neighbor["lastConStateChange"].(string))
		queryData.ControlPlaneNeighbors = append(queryData.ControlPlaneNeighbors, neighborC)
	}

	tflog.Debug(ctx, "LineNeighborDataSource Read:", map[string]interface{}{"discoveredneighbors": discoveredneighbors, "controlplaneneighbors": controlplaneneighbors})
	diags = resp.State.Set(ctx, &queryData)
	resp.Diagnostics.Append(diags...)
}
