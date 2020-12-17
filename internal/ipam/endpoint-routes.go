package ipam

import (
	"context"

	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connection"
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connectioncontext"
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/networkservice"
	"cisco-app-networking.github.io/networkservicemesh/sdk/common"
	"cisco-app-networking.github.io/networkservicemesh/sdk/endpoint"
)

func NewRouteEndpoint(cfg *common.NSConfiguration, _ *string) (server networkservice.NetworkServiceServer, err error) {
	return endpoint.NewCustomFuncEndpoint("route", func(ctx context.Context, c *connection.Connection) error {
		for _, r := range cfg.Routes {
			c.GetContext().GetIpContext().DstRoutes = append(c.GetContext().GetIpContext().DstRoutes, &connectioncontext.Route{
				Prefix: r,
			})
		}
		return nil
	}), nil
}
