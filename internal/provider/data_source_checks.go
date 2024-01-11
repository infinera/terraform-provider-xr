package provider

import (
	"context"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ChecksDataSource{}
	_ datasource.DataSourceWithConfigure = &ChecksDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewChecksDataSource() datasource.DataSource {
	return &ChecksDataSource{}
}

// coffeesDataSource is the data source implementation.
type ChecksDataSource struct {
	client *xrcm_pf.Client
}

type CheckData struct {
	Condition   types.Bool   `tfsdk:"condition"`
	Description types.String `tfsdk:"description"`
	Throw       types.String `tfsdk:"throw"`
}

type ChecksDataSourceData struct {
	Checks []CheckData `tfsdk:"checks"`
}

// Metadata returns the data source type name.
func (d *ChecksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_checks"
}

func (d *ChecksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: " Multiple Checks",
		Attributes: map[string]schema.Attribute{
			"checks": schema.ListNestedAttribute{
				Description: "List of device ids.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"condition": schema.BoolAttribute{
							Description: "condition",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "description",
							Computed:    true,
						},
						"throw": schema.StringAttribute{
							Description: "throw",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ChecksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}
func (d *ChecksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	ChecksData := ChecksDataSourceData{}
	diags := req.Config.Get(ctx, &ChecksData)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "ChecksDataSource: get Checks, request", map[string]interface{}{"ChecksData": ChecksData})

	for _, check := range ChecksData.Checks {
		if check.Condition.ValueBool() {
			resp.Diagnostics.AddError(
				"Checks Condition Failed!!! << "+check.Description.ValueString()+" >>",
				"Error: "+check.Throw.ValueString(),
			)
			return
		}
	}
}
