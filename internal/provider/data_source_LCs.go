package provider

import (
	"context"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LCsDataSource{}
	_ datasource.DataSourceWithConfigure = &LCsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewLCsDataSource() datasource.DataSource {
	return &LCsDataSource{}
}

// coffeesDataSource is the data source implementation.
type LCsDataSource struct {
	client *xrcm_pf.Client
}

type LCsDataSourceData struct {
	N   types.String     `tfsdk:"n"`
	LCs []LCResourceData `tfsdk:"lcs"`
}

// Metadata returns the data source type name.
func (d *LCsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lcs"
}

func (d *LCsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of LCs.",
		Attributes: map[string]schema.Attribute{
			"n": schema.StringAttribute{
				Description: "Device Name",
				Required:    true,
			},
			"lcs": schema.ListNestedAttribute{
				Description: "List of local connections",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Numeric identifier of the Carrier.",
							Computed:    true,
						},
						"n": schema.StringAttribute{
							Description: "XR Device Name",
							Computed:    true,
						},
						"deviceid": schema.StringAttribute{
							Description: "device id",
							Computed:    true,
						},
						"lineptpid": schema.StringAttribute{
							Description: "line ptp id",
							Computed:    true,
						},
						"carrierid": schema.StringAttribute{
							Description: "carrier id",
							Computed:    true,
						},
						"aid": schema.StringAttribute{
							Description: "aid",
							Computed:    true,
						},
						"direction": schema.StringAttribute{
							Description: "direction",
							Computed:    true,
						},
						"lcctrl": schema.Int64Attribute{
							Description: "LC Control",
							Computed:    true,
						},
						"lineaid": schema.StringAttribute{
							Description: "line aid",
							Computed:    true,
						},
						"clientaid": schema.StringAttribute{
							Description: "client aid",
							Computed:    true,
						},
						"dscgaid": schema.StringAttribute{
							Description: "dscg aid",
							Computed:    true,
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
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LCsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *LCsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	queryData := LCsDataSourceData{}
	diags := req.Config.Get(ctx, &queryData)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "LCsDataSource: get LCS ", map[string]interface{}{"queryData": queryData})

	data, deviceId, err := GetResource(ctx, d.client, queryData.N.ValueString(), "resources/lcs")

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Read LCS",
			"LCsDataSource: Could not GET LCS, unexpected error: "+err.Error(),
		)
		return
	}

	resultData := (data["data"].(map[string]interface{}))["content"].(map[string]interface{})
	links := resultData["links"].([]interface{})
	tflog.Debug(ctx, "LCsDataSource: get LC links", map[string]interface{}{"links": links})

	var lcs []LCResourceData

	for _, v := range links {
		lcrec := v.(map[string]interface{})
		data2, err := GetResourcebyID(ctx, d.client, deviceId, "resources/"+lcrec["href"].(string))
		if err != nil {
			continue
		}

		resultData2 := data2["data"].(map[string]interface{})
		lcDataRec := resultData2["content"].(map[string]interface{})
		lcData := LCResourceData{}
		lcData.N = types.StringValue(queryData.N.ValueString())
		lcData.DeviceId = types.StringValue(deviceId)
		lcData.Aid = types.StringValue(lcDataRec["aid"].(string))
		lcData.LcCtrl = types.Int64Value(int64(lcDataRec["lcCtrl"].(float64)))
		lcData.Id = types.StringValue(queryData.N.ValueString() + lcrec["href"].(string))
		lcData.ClientAid = types.StringValue(lcDataRec["clientAid"].(string))
		lcData.LineAid = types.StringValue(lcDataRec["lineAid"].(string))
		lcData.Direction = types.StringValue(lcDataRec["direction"].(string))
		lcData.DscgAid = types.StringValue(lcDataRec["dscgAid"].(string))
		lcData.RemoteModuleId = types.StringValue(lcDataRec["remoteModuleId"].(string))
		lcData.RemoteClientId = types.StringValue(lcDataRec["remoteClientId"].(string))
		lcs = append(lcs, lcData)

	}
	tflog.Debug(ctx, "LCsDataSource: get LCS", map[string]interface{}{"lcs": lcs})
	queryData.LCs = make([]LCResourceData, len(lcs))
	queryData.LCs = lcs
	diags = resp.State.Set(ctx, &queryData)
	resp.Diagnostics.Append(diags...)

}
