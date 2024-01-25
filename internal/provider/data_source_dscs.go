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
	_ datasource.DataSource              = &DSCsDataSource{}
	_ datasource.DataSourceWithConfigure = &DSCsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewDSCsDataSource() datasource.DataSource {
	return &DSCsDataSource{}
}

// coffeesDataSource is the data source implementation.
type DSCsDataSource struct {
	client *xrcm_pf.Client
}

type ModuleDSCsDataSourceData struct {
	N         types.String      `tfsdk:"n"`
	LinePTPId types.String      `tfsdk:"lineptpid"`
	CarrierId types.String      `tfsdk:"carrierid"`
	DSCIds    []types.String    `tfsdk:"dscids"`
	DSCs      []DSCResourceData `tfsdk:"dscs"`
}

type DSCsDataSourceData struct {
	ModuleDSCs []ModuleDSCsDataSourceData `tfsdk:"moduledscs"`
}

// Metadata returns the data source type name.
func (d *DSCsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dscs"
}

func (d *DSCsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Module DSCs",
		Attributes: map[string]schema.Attribute{
			"moduledscs": schema.ListNestedAttribute{
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
						"dscids": schema.ListAttribute{
							Description: "List of DSC IDs",
							Optional:    true,
							ElementType: types.StringType,
						},
						"dscs": schema.ListNestedAttribute{
							Description: "List of DSCs",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Numeric identifier of the DSC.",
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
									"dscid": schema.StringAttribute{
										Description: "DSC ID",
										Computed:    true,
									},
									"cdsc": schema.Int64Attribute{
										Description: "constellation dsc ID",
										Computed:    true,
									},
									"txstatus": schema.StringAttribute{
										Description: "tx status",
										Computed:    true,
									},
									"rxstatus": schema.StringAttribute{
										Description: "Rx status",
										Computed:    true,
									},
									"relativedpo": schema.Int64Attribute{
										Description: "Relative DPO",
										Computed:    true,
									},
									"configstate": schema.StringAttribute{
										Description: "configstate",
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
func (d *DSCsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *DSCsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	queriesData := DSCsDataSourceData{}

	diags := req.Config.Get(ctx, &queriesData)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DSCsDataSource: get DSCs", map[string]interface{}{"queriesData": queriesData})

	var moduleDSCs []ModuleDSCsDataSourceData

	for _, queryData := range queriesData.ModuleDSCs {

		tflog.Debug(ctx, "DSCsDataSource: get dscs", map[string]interface{}{"queryData": queryData})

		data, deviceId, err := GetResource(ctx, d.client, queryData.N.ValueString(), "resource-links")

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Read DSCsDataSource",
				"DSCsDataSource: Could not GET DSC, unexpected error: "+err.Error(),
			)
			return
		}

		tflog.Debug(ctx, "DSCsDataSource: get DSCs", map[string]interface{}{"data": data})

		resources := data["resources"].([]interface{})

		var dscs []DSCResourceData

		for _, v := range resources {
			rec := v.(map[string]interface{})
			rsType := rec["resourceTypes"].([]interface{})
			href := rec["href"].(string)

			if rsType[0].(string) != "xr.carrier.dsc" || !strings.Contains(href, "lineptps/"+queryData.LinePTPId.ValueString()+"/carriers/"+queryData.CarrierId.ValueString()+"/") {
				continue
			}

			dscId := href[strings.LastIndex(href, "/")+1:]

			if len(queryData.DSCIds) > 0 && Find(dscId, queryData.DSCIds) == -1 {
				continue
			}

			data2, err := GetResourcebyID(ctx, d.client, deviceId, "resources/"+rec["href"].(string))
			if err != nil {
				continue
			}

			resultData2 := data2["data"].(map[string]interface{})
			dsc := resultData2["content"].(map[string]interface{})
			dscData := DSCResourceData{}
			dscData.N = types.StringValue(queryData.N.ValueString())
			dscData.DeviceId = types.StringValue(deviceId)
			dscData.LinePTPId = types.StringValue(queryData.LinePTPId.ValueString())
			dscData.CarrierId = types.StringValue(queryData.CarrierId.ValueString())
			dscData.DscId = types.StringValue(dscId)
			dscData.Aid = types.StringValue(dsc["aid"].(string))
			dscData.Id = types.StringValue(queryData.N.ValueString() + href)
			dscData.TxStatus = types.StringValue(dsc["txStatus"].(string))
			dscData.RxStatus = types.StringValue(dsc["rxStatus"].(string))
			dscData.RelativeDPO = types.Int64Value(int64(dsc["relativeDPO"].(float64)))
			dscData.CDsc = types.Int64Value(int64(dsc["cDsc"].(float64)))
			dscData.ConfigState = types.StringValue(dsc["configState"].(string))
			dscs = append(dscs, dscData)
		}
		tflog.Debug(ctx, "dscsDataSource: get dscs", map[string]interface{}{"dscs": dscs})
		queryData.DSCs = make([]DSCResourceData, len(dscs))
		queryData.DSCs = dscs
		moduleDSCs = append(moduleDSCs, queryData)
	}
	tflog.Debug(ctx, "DSCsDataSource: get module DSCs", map[string]interface{}{"DSCs": moduleDSCs})
	queriesData.ModuleDSCs = make([]ModuleDSCsDataSourceData, len(moduleDSCs))
	queriesData.ModuleDSCs = moduleDSCs
	diags = resp.State.Set(ctx, &queriesData)
	resp.Diagnostics.Append(diags...)
}
