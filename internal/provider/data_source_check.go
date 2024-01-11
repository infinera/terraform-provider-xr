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
	_ datasource.DataSource              = &CheckDataSource{}
	_ datasource.DataSourceWithConfigure = &CheckDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewCheckDataSource() datasource.DataSource {
	return &CheckDataSource{}
}

// coffeesDataSource is the data source implementation.
type CheckDataSource struct {
	client *xrcm_pf.Client
}

type CheckDataSourceData struct {
	Condition   types.Bool   `tfsdk:"condition"`
	Description types.String `tfsdk:"description"`
	Throw       types.String `tfsdk:"throw"`
}

// Metadata returns the data source type name.
func (d *CheckDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check"
}

func (d CheckDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of neighbors.",
		Attributes: map[string]schema.Attribute{
			"condition": schema.BoolAttribute{
				Description: "condition",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "description",
				Optional:    true,
			},
			"throw": schema.StringAttribute{
				Description: "throw",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *CheckDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*xrcm_pf.Client)
}

func (d *CheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	checkData := CheckDataSourceData{}
	diags := req.Config.Get(ctx, &checkData)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "CheckDataSource: get Check, request", map[string]interface{}{"checkData": checkData})

	if checkData.Condition.ValueBool() {
		resp.Diagnostics.AddError(
			"Check Condition Failed!! << "+checkData.Description.ValueString()+" >>",
			"Error: "+checkData.Throw.ValueString(),
		)
	}
}
