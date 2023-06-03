package provider

import (
	"context"
	"fmt"
	"github.com/MichaelFraser99/discord-application-sdk/discord/client"
	"github.com/MichaelFraser99/discord-application-sdk/discord/model"
	"github.com/MichaelFraser99/discord-application-sdk/discord/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource = &commandResource{}
)

func NewCommandResource() resource.Resource {
	return &commandResource{}
}

type commandResourceModel struct {
	ApplicationID types.String `tfsdk:"application_id"`
	CommandID     types.String `tfsdk:"command_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Type          types.Int64  `tfsdk:"type"`
	LastUpdated   types.String `tfsdk:"last_updated"`
}

func (c *commandResourceModel) fromCommand(command *model.ApplicationCommand) {
	c.ApplicationID = types.StringValue(command.ApplicationID)
	c.CommandID = types.StringValue(command.ID)
	c.Name = types.StringValue(command.Name)
	c.Description = types.StringValue(command.Description)
	c.Type = types.Int64Value(int64(command.Type))
	c.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
}

type commandResource struct {
	client *client.Client
}

//todo: implement rest of commands api
//todo: optimise this code a bit - lots of repeated code

func (c *commandResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "discord-application_command"
}

func (c *commandResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Discord application command",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				Description: "The application ID that the command belongs to",
				Required:    true,
			},
			"command_id": schema.StringAttribute{
				Description: "The ID of the command",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the command - matches the command a user would type in discord",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the command - displayed in discord",
				Required:    true,
			},
			"type": schema.Int64Attribute{
				Description: "The type of command - see discord application API documentation for more info",
				Required:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "The last time the command was updated",
				Computed:    true,
			},
		},
	}
}

func (c *commandResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	c.client = client
}

func (c *commandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, "-")
	resource.ImportStatePassthroughID(ctx, path.Root("application_id"), resource.ImportStateRequest{ID: ids[0]}, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("command_id"), resource.ImportStateRequest{ID: ids[1]}, resp)
}

func (c *commandResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// Retrieve values from plan
	var plan commandResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createApplicationCommand := &model.CreateApplicationCommand{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        utils.Int(int(plan.Type.ValueInt64())),
	}

	// Create new command
	command, sdkResponse, err := c.client.ApplicationCommand.CreateCommand(ctx, plan.ApplicationID.ValueString(), createApplicationCommand)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating command",
			"Could not create command, unexpected error: "+err.Error(),
		)
		return
	}

	if sdkResponse.StatusCode != 201 && sdkResponse.StatusCode != 200 {
		response.Diagnostics.AddError(
			"Error creating command",
			fmt.Sprintf("Could not create command, unexpected status code: %d", sdkResponse.StatusCode),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.fromCommand(command)

	// Set state to fully populated data
	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *commandResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state commandResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Get refreshed command value from discord
	command, sdkResponse, err := c.client.ApplicationCommand.GetCommand(ctx, state.ApplicationID.ValueString(), state.CommandID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Discord Application Command",
			"Could not read Discord Application Command | ID: "+state.CommandID.ValueString()+" | Application ID: "+state.ApplicationID.ValueString()+" | Error: "+err.Error(),
		)
		return
	}

	if sdkResponse.StatusCode != 200 {
		response.Diagnostics.AddError(
			"Error Reading Discord Application Command",
			"Could not read Discord Application Command | ID: "+state.CommandID.ValueString()+" | Application ID: "+state.ApplicationID.ValueString()+": "+sdkResponse.Status,
		)
		return
	}

	// Overwrite items with refreshed state
	state.fromCommand(command)

	// Set refreshed state
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *commandResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state commandResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan commandResourceModel
	diags = request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	command := model.PatchApplicationCommand{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
	}

	// Update existing command
	updatedCommand, sdkResponse, err := c.client.ApplicationCommand.PatchCommand(ctx, plan.ApplicationID.ValueString(), state.CommandID.ValueString(), &command)
	if err != nil {
		response.Diagnostics.AddError(
			"Error Updating Discord Application Command",
			"Could not update Discord Application Command ID: "+state.CommandID.ValueString()+": "+err.Error(),
		)
		return
	}

	if sdkResponse.StatusCode != 200 {
		response.Diagnostics.AddError(
			"Error Updating Discord Application Command",
			"Could not update Discord Application Command ID: "+plan.CommandID.ValueString()+": "+sdkResponse.Status,
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.fromCommand(updatedCommand)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *commandResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state commandResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	sdkResponse, err := c.client.ApplicationCommand.DeleteCommand(ctx, state.ApplicationID.ValueString(), state.CommandID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error Deleting Discord Application Command",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}

	if sdkResponse.StatusCode != 204 {
		response.Diagnostics.AddError(
			"Error Deleting Discord Application Command",
			"Could not delete Discord Application Command ID "+state.CommandID.ValueString()+": "+sdkResponse.Status,
		)
		return
	}
}
