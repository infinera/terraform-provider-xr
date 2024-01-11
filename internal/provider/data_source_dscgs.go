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
	_ datasource.DataSource              = &DSCGsDataSource{}
	_ datasource.DataSourceWithConfigure = &DSCGsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewDSCGsDataSource() datasource.DataSource {
	return &DSCGsDataSource{}
}

// coffeesDataSource is the data source implementation.
type DSCGsDataSource struct {
	client *xrcm_pf.Client
}

type ModuleDSCGsDataSourceData struct {
	N         types.String       `tfsdk:"n"`
	LinePTPId types.String       `tfsdk:"lineptpid"`
	CarrierId types.String       `tfsdk:"carrierid"`
	DSCGIds   []types.String     `tfsdk:"dscgids"`
	DSCGs     []DSCGResourceData `tfsdk:"dscgs"`
}

type DSCGsDataSourceData struct {
	ModuleDSCGs []ModuleDSCGsDataSourceData `tfsdk:"moduledscgs"`
}

// Metadata returns the data source type name.
func (d *DSCGsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dscgs"
}

func (d *DSCGsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Module DSCGs",
		Attributes: map[string]schema.Attribute{
			"moduledscgs": schema.ListNestedAttribute{
				Description: "List of module's ethernets",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"n": schema.StringAttribute{
							Description: "Device Name",
							Required:    true,
						},
						"lineptpid": schema.StringAttribute{
							Description: "line ptp id",
							Required:    true,
						},
						"carrierid": schema.StringAttribute{
							Description: "carrier id",
							Required:    true,
						},
						"dscgids": schema.ListAttribute{
							Description: "List of DSCG IDs",
							Optional:    true,
							ElementType: types.StringType,
						},
						"dscgs": schema.ListNestedAttribute{
							Description: "List of DSCGs",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "DSC TF id",
										Computed:    true,
									},
									"n": schema.StringAttribute{
										Description: "Device Name",
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
									"dscgid": schema.StringAttribute{
										Description: "DSCG id",
										Computed:    true,
									},
									"aid": schema.StringAttribute{
										Description: "DSCG AID",
										Computed:    true,
									},
									"txdscs": schema.ListAttribute{
										Description: "List of TX DSC IDs",
										Computed:    true,
										ElementType: types.Int64Type,
									},
									"rxdscs": schema.ListAttribute{
										Description: "List of RX DSC IDs",
										Computed:    true,
										ElementType: types.Int64Type,
									},
									"idlecdscs": schema.ListAttribute{
										Description: "Idle DSCs",
										Computed:    true,
										ElementType: types.Int64Type,
									},
									"dscgctrl": schema.Int64Attribute{
										Description: "dscg ctrl",
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
func (d *DSCGsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *DSCGsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	queriesData := DSCGsDataSourceData{}

	diags := req.Config.Get(ctx, &queriesData)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "DSCGsDataSource: get DSCG", map[string]interface{}{"queriesData": queriesData})

	var moduleDSCGs []ModuleDSCGsDataSourceData

	for _, queryData := range queriesData.ModuleDSCGs {
		tflog.Debug(ctx, "DSCGsDataSource Read: get DSCGs", map[string]interface{}{"queryData": queryData})

		data, deviceId, err := GetResource(ctx, d.client, queryData.N.ValueString(), "resources/lineptps/"+queryData.LinePTPId.ValueString()+"/carriers/"+queryData.CarrierId.ValueString()+"/dscgs")

		if err != nil {
			resp.Diagnostics.AddError(
				"DSCGsDataSource Read: Error Get DSCGs",
				"Read: Could not GET DSCGs, unexpected error: "+err.Error(),
			)
			return
		}
		tflog.Debug(ctx, "DSCGsDataSource: get DSCGs", map[string]interface{}{"data": data})

		resultData := (data["data"].(map[string]interface{}))["content"].(map[string]interface{})
		links := resultData["links"].([]interface{})
		tflog.Debug(ctx, "DSCGsDataSource: get DSCG links", map[string]interface{}{"links": links})
		var dscgs []DSCGResourceData

		for _, v := range links {
			dscgrec := v.(map[string]interface{})
			href := dscgrec["href"].(string)
			dscgId := href[strings.LastIndex(href, "/")+1:]
			if len(queryData.DSCGIds) > 0 && Find(dscgId, queryData.DSCGIds) == -1 {
				continue
			}

			data2, err := GetResourcebyID(ctx, d.client, deviceId, "resources/"+dscgrec["href"].(string))
			if err != nil {
				continue
			}

			resultData2 := data2["data"].(map[string]interface{})
			dscgDataRec := resultData2["content"].(map[string]interface{})
			dscgData := DSCGResourceData{}
			dscgData.N = types.StringValue(queryData.N.ValueString())
			dscgData.DscgId = types.StringValue(dscgId)
			dscgData.Id = types.StringValue(queryData.N.ValueString() + href)
			dscgData.Aid = types.StringValue(dscgDataRec["aid"].(string))
			idleCDSCList := getBits(int(dscgDataRec["idleCDSCs"].(float64)))
			dscgData.IdleCDSCs, _ = types.ListValue(types.Int64Type, idleCDSCList)
			dscgData.DscgCtrl = types.Int64Value(int64(dscgDataRec["dscgCtrl"].(float64)))
			dscgData.DeviceId = types.StringValue(deviceId)
			rxCDSCList := getBits(int(dscgDataRec["rxCDSCs"].(float64)))
			dscgData.RxCDSCs, _ = types.ListValue(types.Int64Type, rxCDSCList)
			txCDSCList := getBits(int(dscgDataRec["txCDSCs"].(float64)))
			dscgData.TxCDSCs, _ = types.ListValue(types.Int64Type, txCDSCList)
			dscgs = append(dscgs, dscgData)
		}
		tflog.Debug(ctx, "DSCGsDataSource: get carriers", map[string]interface{}{"N": queryData.N.ValueString(), "linePTPId": queryData.LinePTPId.ValueString(), "CarrierId": queryData.CarrierId.ValueString(), "module DSCG IDs": queryData.DSCGIds, "Module DSCGs": dscgs})
		queryData.DSCGs = make([]DSCGResourceData, len(dscgs))
		queryData.DSCGs = dscgs
		moduleDSCGs = append(moduleDSCGs, queryData)
	}
	tflog.Debug(ctx, "DSCGsDataSource: get module DSCGs", map[string]interface{}{"DSCGs": moduleDSCGs})
	queriesData.ModuleDSCGs = make([]ModuleDSCGsDataSourceData, len(moduleDSCGs))
	queriesData.ModuleDSCGs = moduleDSCGs
	diags = resp.State.Set(ctx, &queriesData)
	resp.Diagnostics.Append(diags...)
}
