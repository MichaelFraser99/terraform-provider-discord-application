package applicationCommands

import (
	"context"
	"github.com/MichaelFraser99/discord-application-sdk/discord/model"
	"net/http"
)

type ApplicationCommandsAPI interface {
	GetCommands(ctx context.Context, applicationID string) (output *[]model.ApplicationCommand, resp *http.Response, err error)
	
	GetCommand(ctx context.Context, applicationID, commandID string) (output *model.ApplicationCommand, resp *http.Response, err error)

	CreateCommand(ctx context.Context, applicationID string, request *model.CreateApplicationCommand) (output *model.ApplicationCommand, resp *http.Response, err error)

	PatchCommand(ctx context.Context, applicationID, commandID string, request *model.PatchApplicationCommand) (output *model.ApplicationCommand, resp *http.Response, err error)

	DeleteCommand(ctx context.Context, applicationID, commandID string) (resp *http.Response, err error)
}
