package provider

import (
	"context"
	"os"

	"terraform-provider-xrcm/internal/xrcm_pf"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &XRProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &XRProvider{}
}

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type XRProvider struct {
	version string
}

// providerData can be used to store data from the Terraform configuration.
type XRProviderModel struct {
	Username types.String `tfsdk:"username"`
	Host     types.String `tfsdk:"host"`
	Password types.String `tfsdk:"password"`
}

func (p *XRProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "xrcm"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *XRProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with XR",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for XR API. May also be provided via XR_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for XR API. May also be provided via XR_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for XR API. May also be provided via XR_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *XRProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config XRProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	tflog.Debug(ctx, "XRProvider: Configure", map[string]interface{}{"config": config})

	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a user to the provider
	if config.Host.IsUnknown() {
		// Cannot connect to client with an unknown Host value
		resp.Diagnostics.AddAttributeError(
			path.Root("Host"),
			"Unknown client Host",
			"The provider cannot create the XR API client as there is an unknown configuration value for the API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the XR_HOST environment variable.",
		)
	}

	// User must provide a user to the provider
	if config.Username.IsUnknown() {
		// Cannot connect to client with an unknown  Username value
		resp.Diagnostics.AddAttributeError(
			path.Root("Username"),
			"Unknown client Username",
			"The provider cannot create the XR API client as there is an unknown configuration value for the API Username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the XR_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		// Cannot connect to client with an unknown Password value
		resp.Diagnostics.AddAttributeError(
			path.Root("Password"),
			"Unknown client Password",
			"The provider cannot create the XR API client as there is an unknown configuration value for the API Password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the XR_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("XR_HOST")
	username := os.Getenv("XR_USERNAME")
	password := os.Getenv("XR_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing XR API Host",
			"The provider cannot create the XR API client as there is a missing or empty value for the XR API host. "+
				"Set the host value in the configuration or use the XR_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing XR API Username",
			"The provider cannot create the XR API client as there is a missing or empty value for the XR API username. "+
				"Set the username value in the configuration or use the XR_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing XR API Password",
			"The provider cannot create the XR API client as there is a missing or empty value for the XR API password. "+
				"Set the password value in the configuration or use the XR_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "xr_host", host)
	ctx = tflog.SetField(ctx, "xr_username", username)
	ctx = tflog.SetField(ctx, "xr_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "xr_password")

	tflog.Debug(ctx, "Creating XR client")

	// Create a new XRCM client and set it to the provider client
	client, err := xrcm_pf.NewClient(&host, &username, &password)

	if err != nil {
		resp.Diagnostics.AddError(
			"provider: Unable to create client",
			"Unable to create XRCM client:\n\n"+err.Error(),
		)
		return
	}

	// Make the XR client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Debug(ctx, "provider: XRCM - successful connection request")
	client.Devicemap = make(map[string]string)
}

// DataSources defines the data sources implemented in the provider.
func (p *XRProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewACsDataSource,
		NewCarriersDataSource,
		NewCheckResourcesDataSource,
		NewCheckDataSource,
		NewChecksDataSource,
		NewDetailDevicesDataSource,
		NewDevicesDataSource,
		NewDSCGsDataSource,
		NewDSCsDataSource,
		NewHostNeighborsDataSource,
		NewEthernetsDataSource,
		NewLCsDataSource,
		NewLineNeighborDataSource,
		NewDeviceIdsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *XRProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCarrierResource,
		NewCarrierDiagResource,
		NewCfgResource,
		NewACResource,
		NewDSCDiagResource,
		NewDSCResource,
		NewDSCGResource,
		NewEthernetDiagResource,
		NewEthernetLLDPResource,
		NewEthernetResource,
		NewLCResource,
		NewODUResource,
		NewOTUResource,
		NewLinePTPResource,
	}
}
