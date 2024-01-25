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
	_ datasource.DataSource              = &ACsDataSource{}
	_ datasource.DataSourceWithConfigure = &ACsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewACsDataSource() datasource.DataSource {
	return &ACsDataSource{}
}

// coffeesDataSource is the data source implementation.
type ACsDataSource struct {
	client *xrcm_pf.Client
}

type ModuleACsDataSourceData struct {
	N          types.String     `tfsdk:"n"`
	EthernetId types.String     `tfsdk:"ethernetid"`
	ACIds      []types.String   `tfsdk:"acids"`
	ACs        []ACResourceData `tfsdk:"acs"`
}

type ACsDataSourceData struct {
	ModuleACs []ModuleACsDataSourceData `tfsdk:"moduleacs"`
}

// Metadata returns the data source type name.
func (d *ACsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acs"
}

func (d *ACsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Module ACs",
		Attributes: map[string]schema.Attribute{
			"moduleacs": schema.ListNestedAttribute{
				Description: "List of module's ACs",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"n": schema.StringAttribute{
							Description: "Device Name",
							Required:    true,
						},
						"ethernetid": schema.StringAttribute{
							Description: "ethernet id",
							Required:    true,
						},
						"acids": schema.ListAttribute{
							Description: "List of AC IDs",
							Optional:    true,
							ElementType: types.StringType,
						},
						"acs": schema.ListNestedAttribute{
							Description: "List of acs",
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
									"ethernetid": schema.StringAttribute{
										Description: "ethernet id",
										Computed:    true,
									},
									"aid": schema.StringAttribute{
										Description: "aid",
										Computed:    true,
									},
									"acid": schema.StringAttribute{
										Description: "ac id",
										Computed:    true,
									},
									"capacity": schema.Int64Attribute{
										Description: "capacity",
										Computed:    true,
									},
									"imc": schema.StringAttribute{
										Description: "imc",
										Computed:    true,
									},
									"imc_outer_vid": schema.StringAttribute{
										Description: "imc outer vid",
										Computed:    true,
									},
									"emc": schema.StringAttribute{
										Description: "emc",
										Computed:    true,
									},
									"emc_outer_vid": schema.StringAttribute{
										Description: "emc outer vid",
										Computed:    true,
									},
									"maxpktlen": schema.Int64Attribute{
										Description: "maxpktlen",
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
func (d *ACsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *ACsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	queriesData := ACsDataSourceData{}
	diags := req.Config.Get(ctx, &queriesData)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "ACsDataSource: get AC", map[string]interface{}{"queriesData": queriesData})

	var moduleACs []ModuleACsDataSourceData

	for _, queryData := range queriesData.ModuleACs {

		tflog.Debug(ctx, "ACsDataSource: get ACS, request", map[string]interface{}{"queryData": queryData})

		data, deviceId, err := GetResource(ctx, d.client, queryData.N.ValueString(), "resources/ethernets/"+queryData.EthernetId.ValueString()+"/acs")

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Read ACS",
				"ACsDataSource: Could not GET ACS, unexpected error: "+err.Error(),
			)
			return
		}

		tflog.Debug(ctx, "ACsDataSource: get ACS", map[string]interface{}{"data": data})

		resultData := data["data"].(map[string]interface{})
		content := resultData["content"].(map[string]interface{})
		links := content["links"].([]interface{})
		tflog.Debug(ctx, "ACsDataSource: get ACS links", map[string]interface{}{"links": links})
		var acs []ACResourceData

		for _, v := range links {
			rec := v.(map[string]interface{})
			href := rec["href"].(string)
			acId := href[strings.LastIndex(href, "/")+1:]

			if len(queryData.ACIds) > 0 && Find(acId, queryData.ACIds) == -1 {
				continue
			}
			data2, err := GetResourcebyID(ctx, d.client, deviceId, "resources/"+href)
			if err != nil {
				continue
			}

			resultData2 := data2["data"].(map[string]interface{})
			acDataRec := resultData2["content"].(map[string]interface{})
			acData := ACResourceData{}
			acData.N = types.StringValue(queryData.N.ValueString())
			acData.DeviceId = types.StringValue(deviceId)
			acData.Aid = types.StringValue(acDataRec["aid"].(string))
			acData.Id = types.StringValue(queryData.N.ValueString() + href)
			acData.AcId = types.StringValue(acId)
			var x int64 = int64(acDataRec["capacity"].(float64))
			acData.Capacity = types.Int64Value(x)
			acData.Imc = types.StringValue(acDataRec["imc"].(string))
			acData.ImcOuterVID = types.StringValue(acDataRec["imcOuterVID"].(string))
			acData.Emc = types.StringValue(acDataRec["emc"].(string))
			acData.EmcOuterVID = types.StringValue(acDataRec["emcOuterVID"].(string))
			acData.MaxPktLen = types.Int64Value(int64(acDataRec["maxPktLen"].(float64)))
			acData.ConfigState = types.StringValue(acDataRec["configState"].(string))
			acs = append(acs, acData)
		}
		tflog.Debug(ctx, "ACsDataSource: get carriers", map[string]interface{}{"device Name": queryData.N.ValueString(), "ethernetId": queryData.EthernetId.ValueString(), "ACIDs": queryData.ACIds, "module ACs": acs})
		queryData.ACs = make([]ACResourceData, len(acs))
		queryData.ACs = acs
		moduleACs = append(moduleACs, queryData)
	}
	tflog.Debug(ctx, "ACsDataSource: get module ACs", map[string]interface{}{"Module ACs": moduleACs})
	queriesData.ModuleACs = make([]ModuleACsDataSourceData, len(moduleACs))
	queriesData.ModuleACs = moduleACs
	diags = resp.State.Set(ctx, &queriesData)
	resp.Diagnostics.Append(diags...)

}
