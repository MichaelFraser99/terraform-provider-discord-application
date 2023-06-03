package provider

import (
	"context"
	sdkConfig "github.com/MichaelFraser99/discord-application-sdk/discord/config"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"os"

	"github.com/MichaelFraser99/discord-application-sdk/discord/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &DiscordProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DiscordProvider{
			version: version,
		}
	}
}

// DiscordProvider is the provider implementation.
type DiscordProvider struct {
	version string
}

type DiscordProviderModel struct {
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *DiscordProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "discord-application"
}

// Schema defines the provider-level schema for configuration data.
func (p *DiscordProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
	resp.Schema = schema.Schema{
		Description: "Interact with Discord's Application API",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "Application token to communicate with the Discord Application API. Can be a bot token or an OAuth Bearer token",
				Required:    true,
			},
		},
	}
}

// Configure prepares a Discord API client for data sources and resources.
func (p *DiscordProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config DiscordProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Discord Application API token",
			"Invalid or missing token",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	token := os.Getenv("DISCORD_APPLICATION_TOKEN")
	baseUrl := os.Getenv("DISCORD_APPLICATION_BASE_URL")

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Discord API Token",
			"The provider cannot authenticate without a token provided",
		)
	}

	if baseUrl == "" {
		baseUrl = "https://discord.com/api"
	}

	if resp.Diagnostics.HasError() {
		return
	}

	discordApplicationSdkConfig := sdkConfig.NewConfig(
		sdkConfig.TOKEN_TYPE_BOT,
		token,
		baseUrl,
		http.DefaultClient,
	)
	sdkClient := client.NewClient(discordApplicationSdkConfig)

	resp.DataSourceData = sdkClient
	resp.ResourceData = sdkClient
}

// DataSources defines the data sources implemented in the provider.
func (p *DiscordProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *DiscordProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCommandResource,
	}
}
