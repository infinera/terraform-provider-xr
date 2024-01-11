package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &CheckResourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &CheckResourcesDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewCheckResourcesDataSource() datasource.DataSource {
	return &CheckResourcesDataSource{}
}

// coffeesDataSource is the data source implementation.
type CheckResourcesDataSource struct {
	client *xrcm_pf.Client
}

type AttributeValuesData struct {
	Attribute              types.String `tfsdk:"attribute"`
	IntentValue            types.String `tfsdk:"intentvalue"`
	DeviceValue            types.String `tfsdk:"devicevalue"`
	ControlAttribute       types.String `tfsdk:"controlattribute"`
	IsValueMatch           types.Bool   `tfsdk:"isvaluematch"`
	AttributeControlByHost types.Bool   `tfsdk:"attributecontrolbyhost"`
}

type ResourceData struct {
	Id              types.String          `tfsdk:"id"`
	GrandparentId   types.String          `tfsdk:"grandparentid"`
	ParentId        types.String          `tfsdk:"parentid"`
	ResourceId      types.String          `tfsdk:"resourceid"`
	Aid             types.String          `tfsdk:"aid"`
	AttributeValues []AttributeValuesData `tfsdk:"attributevalues"`
}

type ResourceDataSourceData struct {
	N            types.String   `tfsdk:"n"`
	DeviceId     types.String   `tfsdk:"deviceid"`
	ResourceType types.String   `tfsdk:"resourcetype"`
	Resources    []ResourceData `tfsdk:"resources"`
}

type ResourcesDataSourceData struct {
	Queries []ResourceDataSourceData `tfsdk:"queries"`
}

// Metadata returns the data source type name.
func (d *CheckResourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_resources"
}

func (d CheckResourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Check Resources",
		Attributes: map[string]schema.Attribute{
			"queries": schema.ListNestedAttribute{
				Description: "List of resources to check",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"n": schema.StringAttribute{
							Description: "Device Name",
							Required:    true,
						},
						"deviceid": schema.StringAttribute{
							Description: "device id",
							Computed:    true,
						},
						"resourcetype": schema.StringAttribute{
							Description: "resource type",
							Required:    true,
						},
						"resources": schema.ListNestedAttribute{
							Description: "List of resources",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "resource identifier",
										Computed:    true,
									},
									"aid": schema.StringAttribute{
										Description: "resource aid",
										Computed:    true,
									},
									"parentid": schema.StringAttribute{
										Description: "parent resource identifier",
										Optional:    true,
										Required:    false,
									},
									"grandparentid": schema.StringAttribute{
										Description: "grand resource identifier",
										Optional:    true,
										Required:    false,
									},
									"resourceid": schema.StringAttribute{
										Description: "resource id",
										Optional:    true,
										Required:    false,
									},
									"attributevalues": schema.ListNestedAttribute{
										Description: "List of attribute values",
										Optional:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"attribute": schema.StringAttribute{
													Description: "attribute",
													Optional:    true,
													Required:    false,
												},
												"intentvalue": schema.StringAttribute{
													Description: "intent value",
													Optional:    true,
													Required:    false,
												},
												"devicevalue": schema.StringAttribute{
													Description: "Device value",
													Computed:    true,
												},
												"controlattribute": schema.StringAttribute{
													Description: "control attribute",
													Optional:    true,
													Required:    false,
												},
												"isvaluematch": schema.BoolAttribute{
													Description: "is value match",
													Computed:    true,
												},
												"attributecontrolbyhost": schema.BoolAttribute{
													Description: "attribute control by host",
													Computed:    true,
												},
											},
										},
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
func (d *CheckResourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *CheckResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	resourcesDataSourceData := ResourcesDataSourceData{}

	diags := req.Config.Get(ctx, &resourcesDataSourceData)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "CheckResources: get Resources'Queries", map[string]interface{}{"querysData": resourcesDataSourceData.Queries})

	for index1, deviceQuery := range resourcesDataSourceData.Queries {

		tflog.Debug(ctx, "CheckResources: get Resources for Device", map[string]interface{}{"Device": deviceQuery.N.ValueString(), "Intent Resources": deviceQuery.Resources})

		for index2, intentResource := range deviceQuery.Resources {
			if intentResource.AttributeValues == nil || len(intentResource.AttributeValues) == 0 {
				continue
			}
			var data2 = map[string]interface{}{}
			var deviceId string
			var err = errors.New("")
			switch deviceQuery.ResourceType.ValueString() {
			case "Carrier":
				if !intentResource.ParentId.IsNull() {
					data2, deviceId, err = GetResource(ctx, d.client, deviceQuery.N.ValueString(), "resources/lineptps/"+intentResource.ParentId.ValueString()+"/carriers/"+intentResource.ResourceId.ValueString())
				}
			case "DSC":
				if !intentResource.ParentId.IsNull() || !intentResource.GrandparentId.IsNull() {
					data2, deviceId, err = GetResource(ctx, d.client, deviceQuery.N.ValueString(), "resources/lineptps/"+intentResource.GrandparentId.ValueString()+"/carriers/"+intentResource.ParentId.ValueString()+"/dscs/"+intentResource.ResourceId.ValueString())
				}
			case "DSCG":
				if !intentResource.ParentId.IsNull() || !intentResource.GrandparentId.IsNull() {
					data2, deviceId, err = GetResource(ctx, d.client, deviceQuery.N.ValueString(), "resources/lineptps/"+intentResource.GrandparentId.ValueString()+"/carriers/"+intentResource.ParentId.ValueString()+"/dscgs/"+intentResource.ResourceId.ValueString())
				}
			case "Ethernet":
				data2, deviceId, err = GetResource(ctx, d.client, deviceQuery.N.ValueString(), "resources/ethernets/"+intentResource.ResourceId.ValueString())
			case "AC":
				if !intentResource.ParentId.IsNull() {
					data2, deviceId, err = GetResource(ctx, d.client, deviceQuery.N.ValueString(), "resources/ethernets/"+intentResource.ParentId.ValueString()+"/acs/"+intentResource.ResourceId.ValueString())
				}
			case "LC":
				data2, deviceId, err = GetResource(ctx, d.client, deviceQuery.N.ValueString(), "resources/lcs/"+intentResource.ResourceId.ValueString())
			case "Config":
				data2, deviceId, err = GetResource(ctx, d.client, deviceQuery.N.ValueString(), "resources/cfg")
			case "Device":
				deviceId, found := d.client.GetDeviceIdFromName(deviceQuery.N.ValueString())
				if !found {
					diags.AddError(
						"Error CheckResources",
						"CheckResources: Could not GET Device ID: "+deviceQuery.N.ValueString()+", error = "+err.Error(),
					)
					return
				}

				//data2, deviceId, err = GetResource(ctx, d.provider.client, deviceQuery.N.Value, "devices/"+deviceId)
				body, error := d.client.ExecuteHttpCommand("GET", "devices/"+deviceId, nil)
				if error != nil {
					diags.AddError(
						"Error CheckResources",
						"CheckResources: Could not GET Device: "+deviceQuery.N.ValueString()+", error = "+err.Error(),
					)
					return
				}
				data2 = make(map[string]interface{})
				err = json.Unmarshal(body, &data2)
			default:
				diags.AddError(
					"Error Read CheckResources",
					"CheckResources: Invalid Resource Type.",
				)
				return
			}

			if err != nil && !strings.Contains(err.Error(), "status: 404") {
				diags.AddError(
					"Error Read CheckResources",
					"CheckResources: Could not GET Resources, unexpected error: "+err.Error(),
				)
				return
			}

			//deviceQuery.DeviceId = types.StringValue(deviceId}
			resourcesDataSourceData.Queries[index1].DeviceId = types.StringValue(deviceId)

			tflog.Debug(ctx, "CheckResources: Get Device Resource", map[string]interface{}{"Device Resource": data2})

			data3 := data2["data"].(map[string]interface{})
			resource := data3["content"].(map[string]interface{})
			deviceQuery.Resources[index2].Id = types.StringValue(deviceQuery.N.ValueString() + data3["resourceId"].(map[string]interface{})["href"].(string))
			if resource["aid"] != nil {
				deviceQuery.Resources[index2].Aid = types.StringValue(resource["aid"].(string))
			}
			for index3, v := range intentResource.AttributeValues {
				tflog.Debug(ctx, "CheckResources: Get intent Resource IntentValue", map[string]interface{}{"IntentValue": v})
				var rawValue = resource[v.Attribute.ValueString()]
				var match bool = true
				var strValue string = ""
				if !v.Attribute.IsNull() && rawValue != nil {
					valueType := reflect.TypeOf(rawValue).String()
					//tflog.Debug(ctx, "CheckResources: Get intent Resource IntentValue 1111", "valueType", valueType)
					if valueType == "int" {
						value, _ := strconv.ParseInt(v.IntentValue.ValueString(), 10, 0)
						if rawValue != value {
							match = false
						}
						strValue = strconv.Itoa(rawValue.(int))
					} else if valueType == "float64" {
						value, _ := strconv.ParseFloat(v.IntentValue.ValueString(), 64)
						if rawValue != value {
							match = false
						}
						strValue = fmt.Sprintf("%.2f", rawValue.(float64))
					} else if valueType == "bool" {
						value, _ := strconv.ParseBool(v.IntentValue.ValueString())
						if rawValue != value {
							match = false
						}
						strValue = strconv.FormatBool(rawValue.(bool))
					} else {
						if v.IntentValue.ValueString() != rawValue.(string) {
							match = false
						}
						strValue = rawValue.(string)
					}
					//tflog.Debug(ctx, "CheckResources: Get intent Resource IntentValue 2222", "v.IntentValue.Value", v.IntentValue.Value, "match", match)
					deviceQuery.Resources[index2].AttributeValues[index3].IsValueMatch = types.BoolValue(match)
					deviceQuery.Resources[index2].AttributeValues[index3].DeviceValue = types.StringValue(strValue)
				}
				deviceQuery.Resources[index2].AttributeValues[index3].AttributeControlByHost = types.BoolValue(false)
				if !v.ControlAttribute.IsNull() { //&& resource[v.ControlAttribute.Value] != nil {
					attributeControl := os.Getenv("TF_" + v.ControlAttribute.ValueString())
					if attributeControl == "Host" {
						deviceQuery.Resources[index2].AttributeValues[index3].AttributeControlByHost = types.BoolValue(true)
					}
				}
			}
		}

		tflog.Debug(ctx, "CheckResources: Check Resources", map[string]interface{}{"resourcesDataSourceData": resourcesDataSourceData})
		diags = resp.State.Set(ctx, &resourcesDataSourceData)
		resp.Diagnostics.Append(diags...)
	}
}
