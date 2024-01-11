package provider

import (
	"context"
	"encoding/json"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DetailDevicesDataSource{}
	_ datasource.DataSourceWithConfigure = &DetailDevicesDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewDetailDevicesDataSource() datasource.DataSource {
	return &DetailDevicesDataSource{}
}

// coffeesDataSource is the data source implementation.
type DetailDevicesDataSource struct {
	client *xrcm_pf.Client
}

type DetailDeviceData struct {
	DeviceId            types.String `tfsdk:"deviceid"`
	N                   types.String `tfsdk:"n"`
	ManufacturerName    types.String `tfsdk:"manufacturername"`
	SoftwareVersion     types.String `tfsdk:"sv"`
	PIID                types.String `tfsdk:"piid"`
	Type                types.String `tfsdk:"type"`
	Status              types.String `tfsdk:"status"`
	ConfiguredRole      types.String `tfsdk:"configuredrole"`
	CurrentRole         types.String `tfsdk:"currentrole"`
	SerdesRate          types.String `tfsdk:"serdesrate"`
	TrafficMode         types.String `tfsdk:"trafficmode"`
	TcMode              types.Bool   `tfsdk:"tcmode"`
	RoleStatus          types.String `tfsdk:"rolestatus"`
	RestartAction       types.String  `tfsdk:"restartaction"`
	FactoryResetAction  types.Bool   `tfsdk:"factoryresetaction"`
	Mnfv                types.String `tfsdk:"mnfv"`
	Mnmn                types.String `tfsdk:"mnmn"`
	Mnmo                types.String `tfsdk:"mnmo"`
	Mnhw                types.String `tfsdk:"mnhw"`
	Mndt                types.String `tfsdk:"mndt"`
	Mnsel               types.String `tfsdk:"mnsel"`
	Clei                types.String `tfsdk:"clei"`
	MacAddress          types.String `tfsdk:"macaddress"`
	ConnectorType       types.String `tfsdk:"connectortype"`
	FormFactor          types.String `tfsdk:"formfactor"`
	//Capabilities        types.Map    `tfsdk:"capabilities"`
}

type DetailDevicesDataSourceData struct {
	State   types.String       `tfsdk:"state"`
	Names   []types.String     `tfsdk:"names"`
	Devices []DetailDeviceData `tfsdk:"devices"`
}

// Metadata returns the data source type name.
func (d *DetailDevicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detaildevices"
}

func (d *DetailDevicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Device infos",
		Attributes: map[string]schema.Attribute{
			"state": schema.StringAttribute{
				Description: "Device state",
				Optional:    true,
			},
			"names": schema.ListAttribute{
				Description: "List of device names",
				Required:    true,
				ElementType: types.StringType,
			},
			"devices": schema.ListNestedAttribute{
				Description: "List of devices'infos",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"n": schema.StringAttribute{
							Description: "Device Name",
							Computed:    true,
						},
						"deviceid": schema.StringAttribute{
							Description: " ID",
							Computed:    true,
						},
						"manufacturername": schema.StringAttribute{
							Description: "manufacturer name",
							Computed:    true,
						},
						"sv": schema.StringAttribute{
							Description: "Device software version",
							Computed:    true,
						},
						"piid": schema.StringAttribute{
							Description: "piid",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Device type",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Device status",
							Computed:    true,
						},
						"configuredrole": schema.StringAttribute{
							Description: "Device configured role",
							Computed:    true,
						},
						"currentrole": schema.StringAttribute{
							Description: "Device current role",
							Computed:    true,
						},
						"trafficmode": schema.StringAttribute{
							Description: "Device traffic mode",
							Computed:    true,
						},
						"serdesrate": schema.StringAttribute{
							Description: "serdes rate",
							Computed:    true,
						},
						"tcmode": schema.BoolAttribute{
							Description: "Device tc mode",
							Computed:    true,
						},
						"rolestatus": schema.StringAttribute{
							Description: "role status",
							Computed:    true,
						},
						"restartaction": schema.StringAttribute{
							Description: "restart action",
							Computed:    true,
						},
						"factoryresetaction": schema.BoolAttribute{
							Description: "Factory Reset Action",
							Computed:    true,
						},
						"mnfv": schema.StringAttribute{
							Description: "mnfv",
							Computed:    true,
						},
						"mnmn": schema.StringAttribute{
							Description: "mnmn",
							Computed:    true,
						},
						"mnmo": schema.StringAttribute{
							Description: "mnmo",
							Computed:    true,
						},
						"mnhw": schema.StringAttribute{
							Description: "mnhw",
							Computed:    true,
						},
						"mndt": schema.StringAttribute{
							Description: "mndt",
							Computed:    true,
						},
						"mnsel": schema.StringAttribute{
							Description: "mnsel",
							Computed:    true,
						},
						"clei": schema.StringAttribute{
							Description: "clei",
							Computed:    true,
						},
						"macaddress": schema.StringAttribute{
							Description: "macAddress",
							Computed:    true,
						},
						"connectortype": schema.StringAttribute{
							Description: "connectorType",
							Computed:    true,
						},
						"formfactor": schema.StringAttribute{
							Description: "formFactor",
							Computed:    true,
						},
						/*"capabilities": schema.MapAttribute{
							Description: "capabilities",
							Computed:    true,
							ElementType: types.StringType,
						},*/
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DetailDevicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *DetailDevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	data := DetailDevicesDataSourceData{}
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(data.Names) == 0 {
		diags.AddError(
			"DetailDevicesDataSource: Names is empty",
			"Read: Names must be specified ",
		)
		return
	}

	var devices []DetailDeviceData
	for _, name := range data.Names {
		deviceId, found := d.client.GetDeviceIdFromName(name.ValueString())
		if !found {
			continue
		}
		queryStr := "devices/" + deviceId

		body, err := d.client.ExecuteHttpCommand("GET", queryStr, nil)
		if err != nil {
			if !strings.Contains(err.Error(), "status: 404") {
				diags.AddError(
					"DetailDevicesDataSource: read ##: Error Get Device: "+name.ValueString(),
					"Read: Could not Get , unexpected error: "+err.Error(),
				)
				return
			}
			tflog.Debug(ctx, "DetailDevicesDataSource: read - not found device "+name.ValueString())
			continue
		}

		tflog.Debug(ctx, "DetailDevicesDataSource Read: get device", map[string]interface{}{"queryStr": queryStr, "body": string(body)})
		var result = make(map[string]interface{})
		json.Unmarshal([]byte(body), &result)
		deviceData := DetailDeviceData{}
		connection := (result["metadata"].(map[string]interface{}))["connection"]
		status := connection.(map[string]interface{})["status"].(string)
		if data.State.ValueString() != "" && data.State.ValueString() != status {
			continue
		}
		deviceData.DeviceId = types.StringValue(result["id"].(string))
		manu := (result["manufacturerName"].([]interface{}))[0]
		manuValue := manu.(map[string]interface{})["value"].(string)
		deviceData.ManufacturerName = types.StringValue(manuValue)
		deviceData.N = types.StringValue(result["name"].(string))
		deviceData.Status = types.StringValue(status)
		xrTypes := result["types"].([]interface{})
		deviceData.Type = types.StringValue(xrTypes[0].(string))
		content := (result["data"].(map[string]interface{}))["content"].(map[string]interface{})
		deviceData.SoftwareVersion = types.StringValue(content["sv"].(string))
		deviceData.PIID = types.StringValue(content["piid"].(string))

		//get device config
		config, err := GetResourcebyID(ctx, d.client, deviceId, "resources/cfg")
		if err != nil {
			tflog.Debug(ctx, "DetailDevicesDataSource Read: get resources/cfg FAILED", map[string]interface{}{"device name": name.ValueString()})
			continue
		}
		resultData2 := config["data"].(map[string]interface{})
		cfgRec := resultData2["content"].(map[string]interface{})
		deviceData.ConfiguredRole = types.StringValue(cfgRec["configuredRole"].(string))
		deviceData.CurrentRole = types.StringValue(cfgRec["currentRole"].(string))
		deviceData.SerdesRate = types.StringValue(cfgRec["serdesRate"].(string))
		deviceData.TrafficMode = types.StringValue(cfgRec["trafficMode"].(string))
		deviceData.TcMode = types.BoolValue(cfgRec["tcMode"].(bool))
		deviceData.RoleStatus = types.StringValue(cfgRec["roleStatus"].(string))
		deviceData.RestartAction = types.StringValue(cfgRec["restartAction"].(string))
		deviceData.FactoryResetAction = types.BoolValue(cfgRec["factoryResetAction"].(bool))

		//get device config
		platform, err := GetResourcebyID(ctx, d.client, deviceId, "resources/oic/p")
		if err != nil {
			tflog.Debug(ctx, "DetailDevicesDataSource Read: get resources/oic/p FAILED", map[string]interface{}{"device name": name.ValueString()})
			continue
		}
		tflog.Debug(ctx, "DetailDevicesDataSource Read: get resources/oic/p SUCCESS", map[string]interface{}{"platform": platform})
		resultData2 = platform["data"].(map[string]interface{})
		platformRec := resultData2["content"].(map[string]interface{})
		deviceData.Mnfv = types.StringValue(platformRec["mnfv"].(string))
		deviceData.Mnmn = types.StringValue(platformRec["mnmn"].(string))
		deviceData.Mnhw = types.StringValue(platformRec["mnhw"].(string))
		deviceData.Mndt = types.StringValue(platformRec["mndt"].(string))
		deviceData.Mnsel = types.StringValue(platformRec["mnsel"].(string))
		deviceData.Clei = types.StringValue(platformRec["clei"].(string))
		deviceData.MacAddress = types.StringValue(platformRec["macAddress"].(string))
		deviceData.ConnectorType = types.StringValue(platformRec["connectorType"].(string))
		deviceData.FormFactor = types.StringValue(platformRec["formFactor"].(string))
		/*if platformRec["capabilities"] != nil {
			capMap := make(map[string]attr.Value)
			for k, cap := range platformRec["capabilities"].(map[string]interface{})  {
				capMap[k] = types.StringValue(cap.(string))
			}
			deviceData.Capabilities, _ = types.MapValue(types.StringType, capMap)
		} */
		devices = append(devices, deviceData)
	}

	data.Devices = make([]DetailDeviceData, len(devices))
	data.Devices = devices
	tflog.Debug(ctx, "DetailDevicesDataSource Read: get devices", map[string]interface{}{"# Device": len(devices), "devices": data})
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}
