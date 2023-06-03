package config

import (
	"fmt"
	"github.com/MichaelFraser99/discord-application-sdk/discord/utils"
	"net/http"
)

type TokenType string

const (
	TOKEN_TYPE_BOT    TokenType = "Bot"
	TOKEN_TYPE_BEARER           = "Bearer"
)

type Config struct {
	TokenType  TokenType
	Token      string
	BaseUrl    string  `default:"https://discord.com/api"`
	apiVersion *string // omitting will default to the current default version detailed within the discord API documentation
	HTTPClient *http.Client
}

func NewConfig(tokenType TokenType, token string, baseUrl string, httpClient *http.Client) *Config {
	return &Config{
		TokenType:  tokenType,
		Token:      token,
		BaseUrl:    baseUrl,
		apiVersion: utils.String("10"),
		HTTPClient: httpClient,
	}
}

func (c *Config) GetVersionedUrl() string {
	return fmt.Sprintf("%s/v%s", c.BaseUrl, *c.apiVersion)
}
