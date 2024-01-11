package provider

import (
	"context"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &CarriersDataSource{}
	_ datasource.DataSourceWithConfigure = &CarriersDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewCarriersDataSource() datasource.DataSource {
	return &CarriersDataSource{}
}

// coffeesDataSource is the data source implementation.
type CarriersDataSource struct {
	client *xrcm_pf.Client
}

type CarriersResourceData struct {
	N          types.String          `tfsdk:"n"`
	LinePTPId  types.String          `tfsdk:"lineptpid"`
	CarrierIds []types.String        `tfsdk:"carrierids"`
	Carriers   []CarrierResourceData `tfsdk:"carriers"`
}

type CarriersDataSourceData struct {
	ModuleCarriers []CarriersResourceData `tfsdk:"modulecarriers"`
}

// Metadata returns the data source type name.
func (d *CarriersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_carriers"
}

func (d *CarriersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Modules' carries information",
		Attributes: map[string]schema.Attribute{
			"modulecarriers": schema.ListNestedAttribute{
				Description: "List of module's modules' carriers",
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
						"carrierids": schema.ListAttribute{
							Description: "List of carrier ID",
							Optional:    true,
							ElementType: types.StringType,
						},
						"carriers": schema.ListNestedAttribute{
							Description: "List of carriers",
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
									"constellationfrequency": schema.Int64Attribute{
										Description: "Constellation Frequency",
										Optional:    true,
									},
									"operatingfrequency": schema.Int64Attribute{
										Description: "operating frequency",
										Computed:    true,
									},
									"hfrequency": schema.Int64Attribute{
										Description: "Host frequency",
										Computed:    true,
									},
									"aconstellationfrequency": schema.Int64Attribute{
										Description: "Actual Constellation Frequency",
										Computed:    true,
									},
									"maxdscs": schema.Int64Attribute{
										Description: "Max Allowed DSCs",
										Optional:    true,
									},
									"hmaxdscs": schema.Int64Attribute{
										Description: "Host Max Allowed DSCs",
										Computed:    true,
									},
									"omaxdscs": schema.Int64Attribute{
										Description: "Operational Max Allowed DSCs",
										Computed:    true,
									},
									"maxtxdscs": schema.Int64Attribute{
										Description: "Max Tx DSCs",
										Optional:    true,
									},
									"hmaxtxdscs": schema.Int64Attribute{
										Description: "Host Max Tx DSCs",
										Computed:    true,
									},
									"omaxtxdscs": schema.Int64Attribute{
										Description: "Operational Max Tx DSCs",
										Computed:    true,
									},
									"allowedrxcdscs": schema.Int64Attribute{
										Description: "Allowed Rx DSCs",
										Optional:    true,
									},
									"hallowedrxcdscs": schema.Int64Attribute{
										Description: "Host Allowed Rx DSCs",
										Computed:    true,
									},
									"aallowedrxcdscs": schema.Int64Attribute{
										Description: "Actual Allowed Rx DSCs",
										Computed:    true,
									},
									"allowedtxcdscs": schema.Int64Attribute{
										Description: "Allowed Tx DSCs",
										Optional:    true,
									},
									"hallowedtxcdscs": schema.Int64Attribute{
										Description: "Host Allowed Tx DSCs",
										Computed:    true,
									},
									"aallowedtxcdscs": schema.Int64Attribute{
										Description: "Actual Allowed Tx DSCs",
										Computed:    true,
									},
									"txclptarget": schema.Int64Attribute{
										Description: "Tx CLP Target",
										Optional:    true,
									},
									"htxclptarget": schema.Int64Attribute{
										Description: "Host Tx CLP Target",
										Computed:    true,
									},
									"atxclptarget": schema.Int64Attribute{
										Description: "Actual Tx CLP Target",
										Computed:    true,
									},
									"advlinectrl": schema.StringAttribute{
										Description: "adv Line Ctrl",
										Computed:    true,
									},
									"spectralbandwidth": schema.Int64Attribute{
										Description: "Spectral bandwidth",
										Computed:    true,
									},
									"modulation": schema.StringAttribute{
										Description: "modulation",
										Optional:    true,
									},
									"omodulation": schema.StringAttribute{
										Description: "Operational modulationControl",
										Computed:    true,
									},
									"hmodulation": schema.StringAttribute{
										Description: "Host modulationControl",
										Computed:    true,
									},
									"clientportmode": schema.StringAttribute{
										Description: "client port mode",
										Optional:    true,
									},
									"feciterations": schema.StringAttribute{
										Description: "fec Iterations",
										Optional:    true,
									},
									"ofeciterations": schema.StringAttribute{
										Description: "Operational fec Iterations",
										Computed:    true,
									},
									"hfeciterations": schema.StringAttribute{
										Description: "Host fec Iterations",
										Computed:    true,
									},
									"baudrate": schema.Int64Attribute{
										Description: "baud rate",
										Optional:    true,
									},
									"capabilities": schema.MapAttribute{
										Description: "capabilities",
										Optional:    true,
										ElementType: types.StringType,
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
func (d *CarriersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *CarriersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	queriesData := CarriersDataSourceData{}

	diags := req.Config.Get(ctx, &queriesData)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "ModuleCarriersDataSource: get Carriers", map[string]interface{}{"queriesData": queriesData})

	var modulecarriers []CarriersResourceData

	for _, queryData := range queriesData.ModuleCarriers {

		tflog.Debug(ctx, "ModuleCarriersDataSource: get Carriers", map[string]interface{}{"queryData": queryData})

		data, deviceId, err := GetResource(ctx, d.client, queryData.N.ValueString(), "resource-links")

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Read ModuleCarriersDataSource",
				"ModuleCarriersDataSource: Could not GET Carriers, unexpected error: "+err.Error(),
			)
			return
		}

		tflog.Debug(ctx, "ModuleCarriersDataSource: get Carriers", map[string]interface{}{"data": data})

		resources := data["resources"].([]interface{})

		var carriers []CarrierResourceData

		for _, v := range resources {
			rec := v.(map[string]interface{})
			rsType := rec["resourceTypes"].([]interface{})
			href := rec["href"].(string)

			if rsType[0].(string) != "xr.carrier" || !strings.Contains(href, "lineptps/"+queryData.LinePTPId.ValueString()) {
				continue
			}
			carrierId := href[strings.LastIndex(href, "/")+1:]

			if len(queryData.CarrierIds) > 0 && Find(carrierId, queryData.CarrierIds) == -1 {
				continue
			}

			data2, err := GetResourcebyID(ctx, d.client, deviceId, "resources"+rec["href"].(string))
			if err != nil {
				continue
			}

			resultData2 := data2["data"].(map[string]interface{})
			carrier := resultData2["content"].(map[string]interface{})
			carrierData := CarrierResourceData{}
			carrierData.N = queryData.N
			carrierData.DeviceId = types.StringValue(deviceId)
			carrierData.LinePTPId = queryData.LinePTPId
			carrierData.CarrierId = types.StringValue(carrierId)
			carrierData.Aid = types.StringValue(carrier["aid"].(string))
			carrierData.Id = types.StringValue(queryData.N.ValueString() + href)
			carrierData.FecIterations = types.StringValue(carrier["fecIterations"].(string))
			carrierData.AdvLineCtrl = types.StringValue(carrier["advLineCtrl"].(string))
			carrierData.Modulation = types.StringValue(carrier["modulation"].(string))
			carrierData.ClientPortMode = types.StringValue(carrier["clientPortMode"].(string))
			carrierData.ConstellationFrequency = types.Int64Value(int64(carrier["constellationFrequency"].(float64)))
			carrierData.BaudRate = types.Int64Value(int64(carrier["baudRate"].(float64)))
			carrierData.MaxDSCs = types.Int64Value(int64(carrier["maxDSCs"].(float64)))
			carrierData.MaxTxDSCs = types.Int64Value(int64(carrier["maxTxDSCs"].(float64)))
			carrierData.SpectralBandwidth = types.Int64Value(int64(carrier["spectralBandwidth"].(float64)))
			carrierData.TxCLPtarget = types.Int64Value(int64(carrier["txCLPtarget"].(float64)))
			carrierData.AllowedTxCDSCs = types.Int64Value(int64(carrier["allowedTxCDSCs"].(float64)))
			carrierData.AllowedRxCDSCs = types.Int64Value(int64(carrier["allowedRxCDSCs"].(float64)))
			carrierData.HModulation = types.StringValue(carrier["hModulation"].(string))
			carrierData.OModulation = types.StringValue(carrier["oModulation"].(string))
			carrierData.HFecIterations = types.StringValue(carrier["hFecIterations"].(string))
			carrierData.OFecIterations = types.StringValue(carrier["oFecIterations"].(string))
			carrierData.HFrequency = types.Int64Value(int64(carrier["hFrequency"].(float64)))
			carrierData.AConstellationFrequency = types.Int64Value(int64(carrier["aConstellationFrequency"].(float64)))
			carrierData.OperatingFrequency = types.Int64Value(int64(carrier["operatingFrequency"].(float64)))
			carrierData.HTxCLPtarget = types.Int64Value(int64(carrier["hTxCLPtarget"].(float64)))
			carrierData.ATxCLPtarget = types.Int64Value(int64(carrier["aTxCLPtarget"].(float64)))
			carrierData.OMaxDSCs = types.Int64Value(int64(carrier["oMaxDSCs"].(float64)))
			carrierData.HMaxDSCs = types.Int64Value(int64(carrier["hMaxDSCs"].(float64)))
			carrierData.HMaxTxDSCs = types.Int64Value(int64(carrier["hMaxTxDSCs"].(float64)))
			carrierData.OMaxTxDSCs = types.Int64Value(int64(carrier["oMaxTxDSCs"].(float64)))
			carrierData.HAllowedTxCDSCs = types.Int64Value(int64(carrier["hAllowedTxCDSCs"].(float64)))
			carrierData.AAllowedTxCDSCs = types.Int64Value(int64(carrier["aAllowedTxCDSCs"].(float64)))
			carrierData.HAllowedRxCDSCs = types.Int64Value(int64(carrier["hAllowedRxCDSCs"].(float64)))
			carrierData.AAllowedRxCDSCs = types.Int64Value(int64(carrier["aAllowedRxCDSCs"].(float64)))
			capMap := make(map[string]attr.Value)
			for k,v2 := range carrier["capabilities"].(map[string]interface{})  {
				capMap[k] = types.StringValue(v2.(string))
			}
			carrierData.Capabilities, _ = types.MapValue(types.StringType, capMap)
			carriers = append(carriers, carrierData)
		}
		tflog.Debug(ctx, "ModuleCarriersDataSource: get carriers", map[string]interface{}{"N": queryData.N.ValueString, "LinePTP": queryData.LinePTPId.ValueString, "carriers": queryData.Carriers, "module carriers": carriers})
		queryData.Carriers = make([]CarrierResourceData, len(carriers))
		queryData.Carriers = carriers
		modulecarriers = append(modulecarriers, queryData)
	}

	tflog.Debug(ctx, "ModuleCarriersDataSource: get module carriers", map[string]interface{}{"carriers": modulecarriers})
	queriesData.ModuleCarriers = make([]CarriersResourceData, len(modulecarriers))
	queriesData.ModuleCarriers = modulecarriers
	diags = resp.State.Set(ctx, &queriesData)
	resp.Diagnostics.Append(diags...)

}
