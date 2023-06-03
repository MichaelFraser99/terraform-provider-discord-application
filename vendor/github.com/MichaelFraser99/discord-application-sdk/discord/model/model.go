package model

type ApplicationCommand struct {
	ID                       string                      `json:"id"`
	Type                     int                         `json:"type,omitempty"`
	ApplicationID            string                      `json:"application_id,omitempty"`
	GuildID                  *string                     `json:"guild_id,omitempty"`
	Name                     string                      `json:"name"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	Description              string                      `json:"description"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	Options                  *[]ApplicationCommandOption `json:"options,omitempty"` //When defining multiple, ensure required values are listed before optional values.
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DmPermission             bool                        `json:"dm_permission"`
	DefaultPermission        *bool                       `json:"default_permission,omitempty"`
	Nsfw                     bool                        `json:"nsfw"`
	Version                  string                      `json:"version"`
}

type ApplicationCommandOption struct {
	Type                     int                               `json:"type"`
	Name                     string                            `json:"name"`
	NameLocalizations        *map[string]string                `json:"name_localizations,omitempty"`
	Description              string                            `json:"description"`
	DescriptionLocalizations *map[string]string                `json:"description_localizations,omitempty"`
	Required                 *bool                             `json:"required,omitempty"`
	Choices                  *[]ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options                  *[]ApplicationCommandOption       `json:"options,omitempty"` //When defining multiple, ensure required values are listed before optional values.
	ChannelTypes             *[]int                            `json:"channel_types,omitempty"`
	MinValue                 *int                              `json:"min_value,omitempty"`
	MaxValue                 *int                              `json:"max_value,omitempty"`
	MinLength                *int                              `json:"min_length,omitempty"`
	MaxLength                *int                              `json:"max_length,omitempty"`
	AutoComplete             *bool                             `json:"autocomplete,omitempty"` //Must be false if choices are defined
}

type ApplicationCommandOptionChoice struct {
	Name              string             `json:"name"`
	NameLocalizations *map[string]string `json:"name_localizations,omitempty"`
	Value             interface{}        `json:"value"` //Can be string, integer, or double. Note: if string, max length is 100
}

type CreateApplicationCommand struct {
	Name                     string                      `json:"name"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	Description              string                      `json:"description"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	Options                  *[]ApplicationCommandOption `json:"options,omitempty"` //When defining multiple, ensure required values are listed before optional values.
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DmPermission             bool                        `json:"dm_permission"`
	DefaultPermission        *bool                       `json:"default_permission,omitempty"`
	Type                     *int                        `json:"type,omitempty"` //defaults to 1
	Nsfw                     bool                        `json:"nsfw"`
}

type PatchApplicationCommand struct {
	Name                     *string                     `json:"name,omitempty"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	Options                  *[]ApplicationCommandOption `json:"options,omitempty"` //When defining multiple, ensure required values are listed before optional values.
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DmPermission             *bool                       `json:"dm_permission,omitempty"`
	DefaultPermission        *bool                       `json:"default_permission,omitempty"`
	Nsfw                     *bool                       `json:"nsfw,omitempty"`
}
