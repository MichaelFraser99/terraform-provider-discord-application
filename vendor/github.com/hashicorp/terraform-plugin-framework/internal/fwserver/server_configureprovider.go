package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// ConfigureProvider implements the framework server ConfigureProvider RPC.
func (s *Server) ConfigureProvider(ctx context.Context, req *provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	logging.FrameworkDebug(ctx, "Calling provider defined Provider Configure")

	if req != nil {
		s.Provider.Configure(ctx, *req, resp)
	} else {
		s.Provider.Configure(ctx, provider.ConfigureRequest{}, resp)
	}

	logging.FrameworkDebug(ctx, "Called provider defined Provider Configure")

	s.DataSourceConfigureData = resp.DataSourceData
	s.ResourceConfigureData = resp.ResourceData
}
