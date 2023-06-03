package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/MichaelFraser99/discord-application-sdk/discord/config"
	"net/http"
)

type HTTPMethods string

const (
	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_PUT    = "PUT"
	METHOD_PATCH  = "PATCH"
	METHOD_DELETE = "DELETE"
)

type HTTP struct {
	Config   config.Config
	request  *http.Request
	response *http.Response
}

func NewHTTP(config config.Config) *HTTP {
	return &HTTP{
		Config: config,
	}
}

func (h *HTTP) WithRequest(request *http.Request) *HTTP {
	h.request = request
	return h
}

func (h *HTTP) GetResponseAndClear() http.Response {
	response := *h.response
	h.response = nil
	return response
}

func (h *HTTP) Do(ctx context.Context) error {
	if h.request == nil {
		return errors.New("cannot perform request without a request object")
	}
	if h.Config.HTTPClient == nil {
		return errors.New("cannot perform request without a http client")
	}

	h.request = h.request.WithContext(ctx)
	h.request.Header.Set("Content-Type", "application/json")
	h.request.Header.Set("Authorization", fmt.Sprintf("%s %s", h.Config.TokenType, h.Config.Token))

	do, err := h.Config.HTTPClient.Do(h.request)
	h.response = do
	if err != nil {
		return err
	}
	return nil
}
