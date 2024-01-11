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
	_ datasource.DataSource              = &HostNeighborsDataSource{}
	_ datasource.DataSourceWithConfigure = &HostNeighborsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewHostNeighborsDataSource() datasource.DataSource {
	return &HostNeighborsDataSource{}
}

// coffeesDataSource is the data source implementation.
type HostNeighborsDataSource struct {
	client *xrcm_pf.Client
}

type HostNeighborData struct {
	LocalPortSourceMAC types.String `tfsdk:"localportsourcemac"`
	ChassisIdSubtype   types.String `tfsdk:"chassisidsubtype"`
	ChassisId          types.String `tfsdk:"chassisid"`
	PortIdSubtype      types.String `tfsdk:"portidsubtype"`
	PortId             types.String `tfsdk:"portid"`
	PortDescr          types.String `tfsdk:"portdescr"`
	SysName            types.String `tfsdk:"sysname"`
	SysDescr           types.String `tfsdk:"sysdescr"`
	SysTTL             types.Int64  `tfsdk:"systtl"`
	LldpPdu            types.String `tfsdk:"lldppdu"`
}

type HostNeighborsDataSourceData struct {
	N          types.String       `tfsdk:"n"`
	DeviceId   types.String       `tfsdk:"deviceid"`
	EthernetId types.String       `tfsdk:"ethernetid"`
	Aid        types.String       `tfsdk:"aid"`
	Neighbors  []HostNeighborData `tfsdk:"neighbors"`
}

// Metadata returns the data source type name.
func (d *HostNeighborsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host_neighbors"
}

func (d *HostNeighborsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of host neighbors.",
		Attributes: map[string]schema.Attribute{
			"deviceid": schema.StringAttribute{
				Description: "Device ID",
				Optional:    true,
			},
			"n": schema.StringAttribute{
				Description: "Device Name",
				Required:    true,
			},
			"ethernetid": schema.StringAttribute{
				Description: "Ethernet Id",
				Required:    true,
			},
			"aid": schema.StringAttribute{
				Description: "aid",
				Computed:    true,
			},
			"neighbors": schema.ListNestedAttribute{
				Description: "List of discovered neighbors",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"localportsourcemac": schema.StringAttribute{
							Description: "local port source emac",
							Computed:    true,
						},
						"chassisidsubtype": schema.StringAttribute{
							Description: "chassis id subtype",
							Computed:    true,
						},
						"chassisid": schema.StringAttribute{
							Description: "chassis id",
							Computed:    true,
						},
						"portidsubtype": schema.StringAttribute{
							Description: "port id subtype",
							Computed:    true,
						},
						"portid": schema.StringAttribute{
							Description: "port id",
							Computed:    true,
						},
						"portdescr": schema.StringAttribute{
							Description: "port descr",
							Computed:    true,
						},
						"sysname": schema.StringAttribute{
							Description: "sys name",
							Computed:    true,
						},
						"sysdescr": schema.StringAttribute{
							Description: "sys descr",
							Computed:    true,
						},
						"systtl": schema.Int64Attribute{
							Description: "sys ttl",
							Computed:    true,
						},
						"lldppdu": schema.StringAttribute{
							Description: "lld ppdu",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *HostNeighborsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *HostNeighborsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	queryData := HostNeighborsDataSourceData{}
	diags := req.Config.Get(ctx, &queryData)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "HostNeighborDataSource Read: get HostNeighbor", map[string]interface{}{"req": req})

	body, deviceId, err := d.client.ExecuteDeviceHttpCommand(queryData.N.ValueString(), "GET", "resources/ethernets/"+queryData.EthernetId.ValueString()+"/host-neighbors", nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"HostNeighborDataSource Read: Error Get HostNeighbor",
			"Read: Could not GET HostNeighbor, unexpected error: "+err.Error(),
		)
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"HostNeighborDataSource Read: Error Read HostNeighbor",
			"Read: Could not Unmarshal HostNeighbor, unexpected error: "+err.Error(),
		)
		return
	}

	var resultData = make(map[string]interface{})
	var content = map[string]interface{}{}
	var neighbors = []interface{}{}

	resultData = data["data"].(map[string]interface{})
	tflog.Debug(ctx, "HostNeighborDataSource: get Hostneighbors", map[string]interface{}{"resultData": resultData})
	queryData.Neighbors = make([]HostNeighborData, 0)
	queryData.DeviceId = types.StringValue(deviceId)
	content = resultData["content"].(map[string]interface{})
	queryData.Aid = types.StringValue(content["aid"].(string))
	neighbors = content["neighbors"].([]interface{})
	for _, n := range neighbors {
		neighbor := n.(map[string]interface{})
		hostNeighborData := HostNeighborData{}
		hostNeighborData.LocalPortSourceMAC = types.StringValue(neighbor["localPortSourceMAC"].(string))
		hostNeighborData.ChassisIdSubtype = types.StringValue(neighbor["chassisIdSubtype"].(string))
		hostNeighborData.ChassisId = types.StringValue(neighbor["chassisId"].(string))
		hostNeighborData.PortIdSubtype = types.StringValue(neighbor["portIdSubtype"].(string))
		hostNeighborData.PortId = types.StringValue(neighbor["portId"].(string))
		hostNeighborData.PortDescr = types.StringValue(neighbor["portDescr"].(string))
		hostNeighborData.SysName = types.StringValue(neighbor["sysName"].(string))
		hostNeighborData.SysDescr = types.StringValue(neighbor["sysDescr"].(string))
		hostNeighborData.SysTTL = types.Int64Value(int64(neighbor["sysTTL"].(float64)))
		hostNeighborData.LldpPdu = types.StringValue(neighbor["lldpPDU"].(string))

		queryData.Neighbors = append(queryData.Neighbors, hostNeighborData)
	}
	diags = resp.State.Set(ctx, &queryData)
	resp.Diagnostics.Append(diags...)
}
