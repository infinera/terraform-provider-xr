package provider

import (
	"context"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &EthernetsDataSource{}
	_ datasource.DataSourceWithConfigure = &EthernetsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewEthernetsDataSource() datasource.DataSource {
	return &EthernetsDataSource{}
}

// coffeesDataSource is the data source implementation.
type EthernetsDataSource struct {
	client *xrcm_pf.Client
}

type EthernetsDataSourceData struct {
	N           types.String           `tfsdk:"n"`
	EthernetIds []types.String         `tfsdk:"ethernetids"`
	Ethernets   []EthernetResourceData `tfsdk:"ethernets"`
}

type ModuleEthernetsDataSourceData struct {
	ModuleEthernets []EthernetsDataSourceData `tfsdk:"moduleethernets"`
}

// Metadata returns the data source type name.
func (d *EthernetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ethernets"
}

func (d *EthernetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of ModuleEthernets",
		Attributes: map[string]schema.Attribute{
			"moduleethernets": schema.ListNestedAttribute{
				Description: "List of module's ethernets",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"n": schema.StringAttribute{
							Description: "Device Name",
							Required:    true,
						},
						"ethernetids": schema.ListAttribute{
							Description: "List of ethernet IDs",
							Optional:    true,
							ElementType: types.StringType,
						},
						"ethernets": schema.ListNestedAttribute{
							Description: "List of local connections",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Numeric identifier of the Ethernet.",
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
									"ethernetid": schema.StringAttribute{
										Description: "ethernet id",
										Computed:    true,
									},
									"aid": schema.StringAttribute{
										Description: "AID",
										Computed:    true,
									},
									"fecmode": schema.StringAttribute{
										Description: "fec mode",
										Computed:    true,
									},
									"fectype": schema.StringAttribute{
										Description: "fec type",
										Computed:    true,
									},
									"portspeed": schema.Int64Attribute{
										Description: "fec mode",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *EthernetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *EthernetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	querysData := ModuleEthernetsDataSourceData{}

	diags := req.Config.Get(ctx, &querysData)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "EthernetsDataSource: get EThernet", map[string]interface{}{"querysData": querysData})

	var moduleEthernets []EthernetsDataSourceData

	for _, queryData := range querysData.ModuleEthernets {

		tflog.Debug(ctx, "EthernetsDataSource: get Ehernets", map[string]interface{}{"queryData": queryData})

		data, deviceId, err := GetResource(ctx, d.client, queryData.N.ValueString(), "resource-links")

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Read EthernetsDataSource",
				"EthernetsDataSource: Could not GET Ethernet, unexpected error: "+err.Error(),
			)
			return
		}

		tflog.Debug(ctx, "EthernetsDataSource: get Ethernets", map[string]interface{}{"data": data})

		resources := data["resources"].([]interface{})

		var ethernets []EthernetResourceData

		for _, v := range resources {
			rec := v.(map[string]interface{})
			rsType := rec["resourceTypes"].([]interface{})
			href := rec["href"].(string)

			if rsType[0].(string) != "xr.ethernet" {
				continue
			}

			ethernetId := href[strings.LastIndex(href, "/")+1:]

			if len(queryData.EthernetIds) > 0 && Find(ethernetId, queryData.EthernetIds) == -1 {
				continue
			}

			data2, err := GetResourcebyID(ctx, d.client, deviceId, "resources/"+rec["href"].(string))
			if err != nil {
				continue
			}

			resultData2 := data2["data"].(map[string]interface{})
			ethernet := resultData2["content"].(map[string]interface{})
			ethernetData := EthernetResourceData{}
			ethernetData.N = types.StringValue(queryData.N.ValueString())
			ethernetData.DeviceId = types.StringValue(deviceId)
			ethernetData.EthernetId = types.StringValue(ethernetId)
			ethernetData.Aid = types.StringValue(ethernet["aid"].(string))
			ethernetData.Id = types.StringValue(queryData.N.ValueString() + href)
			ethernetData.FecType = types.StringValue(ethernet["fecType"].(string))
			ethernetData.FecMode = types.StringValue(ethernet["fecMode"].(string))
			ethernetData.PortSpeed = types.Int64Value(int64(ethernet["portSpeed"].(float64)))
			ethernets = append(ethernets, ethernetData)
		}
		tflog.Debug(ctx, "ethernetsDataSource: get ethernets", map[string]interface{}{"ethernets": ethernets})
		queryData.Ethernets = make([]EthernetResourceData, len(ethernets))
		queryData.Ethernets = ethernets
		moduleEthernets = append(moduleEthernets, queryData)
	}
	tflog.Debug(ctx, "ethernetsDataSource: get ethernets", map[string]interface{}{"Module ethernets": moduleEthernets})
	querysData.ModuleEthernets = make([]EthernetsDataSourceData, len(moduleEthernets))
	querysData.ModuleEthernets = moduleEthernets
	diags = resp.State.Set(ctx, &querysData)
	resp.Diagnostics.Append(diags...)

}
