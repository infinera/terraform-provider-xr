package provider

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DeviceIdsDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceIdsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewDeviceIdsDataSource() datasource.DataSource {
	return &DeviceIdsDataSource{}
}

// coffeesDataSource is the data source implementation.
type DeviceIdsDataSource struct {
	client *xrcm_pf.Client
}

type DeviceID struct {
	N  types.String `tfsdk:"n"`
	Id types.String `tfsdk:"id"`
}

type DeviceIdsData struct {
	DeviceIds []DeviceID `tfsdk:"deviceids"`
}

// Metadata returns the data source type name.
func (d *DeviceIdsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_ids"
}

func (d *DeviceIdsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Device Ids.",
		Attributes: map[string]schema.Attribute{
			"deviceids": schema.ListNestedAttribute{
				Description: "List of device ids.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Device Id",
							Computed:    true,
						},
						"n": schema.StringAttribute{
							Description: "Device name",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DeviceIdsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *DeviceIdsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	deviceIdsData := DeviceIdsData{}
	diags := req.Config.Get(ctx, &deviceIdsData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "DeviceIdsDataSource: get Check, request", map[string]interface{}{"IdsData": deviceIdsData})

	path, _ := os.Getwd()

	tflog.Debug(ctx, "DeviceIdsDataSource: ", map[string]interface{}{"path": path + "/terraform.tfstate"})

	// Open our jsonFile
	jsonFile, err := os.Open(path + "/terraform.tfstate")
	// if we os.Open returns an error then handle it
	if err != nil {
		resp.Diagnostics.AddError(
			"DeviceIdsDataSource Failed!!", " Can not open TF state file: "+path+"/terraform.tfstate",
		)
		return
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	//tflog.Debug(ctx, "DeviceIdsDataSource: ", "state", result)

	if result == nil {
		tflog.Debug(ctx, "DeviceIdsDataSource: No devices ids in tf state file")
		resp.State.Set(ctx, &deviceIdsData)
		return
	}

	var deviceIds []DeviceID

	for _, res := range result["resources"].([]interface{}) {
		var resource = res.(map[string]interface{})
		if resource["type"].(string) != "xrcm_detaildevices" {
			continue
		}
		tflog.Debug(ctx, "DeviceIdsDataSource: res", map[string]interface{}{"res": res})
		for _, inst := range resource["instances"].([]interface{}) {
			var instance = inst.(map[string]interface{})
			tflog.Debug(ctx, "DeviceIdsDataSource: instance", map[string]interface{}{"instance": instance})
			for _, dev := range instance["attributes"].(map[string]interface{})["devices"].([]interface{}) {
				var device = dev.(map[string]interface{})
				tflog.Debug(ctx, "DeviceIdsDataSource: device", map[string]interface{}{"device": device})
				deviceId := DeviceID{}
				deviceId.N = types.StringValue(device["name"].(string))
				deviceId.Id = types.StringValue(device["deviceid"].(string))
				deviceIds = append(deviceIds, deviceId)
			}
		}
	}

	deviceIdsData.DeviceIds = make([]DeviceID, len(deviceIds))
	deviceIdsData.DeviceIds = deviceIds
	tflog.Debug(ctx, "DeviceIdsDataSource: get devices' ids", map[string]interface{}{"device IDs": deviceIds})

	diags = resp.State.Set(ctx, &deviceIdsData)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, "DeviceIdsDataSource Read: get devices", map[string]interface{}{"# Device IDs": len(deviceIds)})

}
