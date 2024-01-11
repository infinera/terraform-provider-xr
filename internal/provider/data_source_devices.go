package provider

import (
	"context"
	"encoding/json"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DevicesDataSource{}
	_ datasource.DataSourceWithConfigure = &DevicesDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewDevicesDataSource() datasource.DataSource {
	return &DevicesDataSource{}
}

// coffeesDataSource is the data source implementation.
type DevicesDataSource struct {
	client *xrcm_pf.Client
}

type DeviceData struct {
	DeviceId         types.String `tfsdk:"deviceid"`
	N                types.String `tfsdk:"n"`
	ManufacturerName types.String `tfsdk:"manufacturername"`
	SoftwareVersion  types.String `tfsdk:"sv"`
	PIID             types.String `tfsdk:"piid"`
	Type             types.String `tfsdk:"type"`
	Status           types.String `tfsdk:"status"`
}

type DevicesDataSourceData struct {
	State   types.String   `tfsdk:"state"`
	Names   []types.String `tfsdk:"names"`
	Devices []DeviceData   `tfsdk:"devices"`
}

// Metadata returns the data source type name.
func (d *DevicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (d *DevicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DevicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *DevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	data := DevicesDataSourceData{}

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var devices []DeviceData
	queryStr := "devices"
	if data.Names == nil || len(data.Names) == 0 {

		tflog.Debug(ctx, "DevicesDataSource Read: get devices", map[string]interface{}{"request": req})

		body, _ := d.client.ExecuteHttpCommand("GET", queryStr, nil)

		tflog.Debug(ctx, "DevicesDataSource Read: get devices", map[string]interface{}{"queryStr": queryStr, "body": string(body)})

		tflog.Debug(ctx, "####DevicesDataSource Read: get devices", map[string]interface{}{"queryStr": queryStr, "body": string(body)})
		devicesStr := string(body[:])
		devicesStrArray := strings.Split(devicesStr, "\n\n")

		for _, v := range devicesStrArray {
			if len(v) == 0 {
				continue
			}
			var entity = make(map[string]interface{})
			json.Unmarshal([]byte(v), &entity)
			result := entity["result"].(map[string]interface{})

			connection := (result["metadata"].(map[string]interface{}))["connection"]
			status := connection.(map[string]interface{})["status"].(string)
			if data.State.ValueString() != "" && data.State.ValueString() != status {
				continue
			}
			deviceData := DeviceData{}
			deviceData.DeviceId = types.StringValue(result["id"].(string))
			manu := (result["manufacturerName"].([]interface{}))[0]
			manuValue := manu.(map[string]interface{})["value"].(string)
			deviceData.ManufacturerName = types.StringValue(manuValue)
			deviceData.N = types.StringValue(result["name"].(string))
			deviceData.Status = types.StringValue(status)
			xrTypes := result["types"].([]interface{})
			xrtype := xrTypes[0].(string)
			deviceData.Type = types.StringValue(xrtype)
			content := (result["data"].(map[string]interface{}))["content"]
			svValue := (content.(map[string]interface{}))["sv"].(string)
			deviceData.SoftwareVersion = types.StringValue(svValue)
			piid := (content.(map[string]interface{}))["piid"].(string)
			deviceData.PIID = types.StringValue(piid)
			devices = append(devices, deviceData)
		}
	} else {
		queryStr = "devices"
		for _, name := range data.Names {
			queryStr = "devices"
			deviceId, found := d.client.GetDeviceIdFromName(name.ValueString())
			if !found {
				continue
			}
			queryStr += "/" + deviceId

			body, _ := d.client.ExecuteHttpCommand("GET", queryStr, nil)

			tflog.Debug(ctx, "$$$DevicesDataSource Read: get devices", map[string]interface{}{"queryStr": queryStr, "body": string(body)})
			var result = make(map[string]interface{})
			json.Unmarshal([]byte(body), &result)
			deviceData := DeviceData{}
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
			xrtype := xrTypes[0].(string)
			deviceData.Type = types.StringValue(xrtype)
			content := (result["data"].(map[string]interface{}))["content"]
			svValue := (content.(map[string]interface{}))["sv"].(string)
			deviceData.SoftwareVersion = types.StringValue(svValue)
			piid := (content.(map[string]interface{}))["piid"].(string)
			deviceData.PIID = types.StringValue(piid)
			devices = append(devices, deviceData)
		}
	}

	data.Devices = make([]DeviceData, len(devices))
	data.Devices = devices
	tflog.Debug(ctx, "DevicesDataSource Read: get devices", map[string]interface{}{"# Device": len(devices), "devices": data})
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}
