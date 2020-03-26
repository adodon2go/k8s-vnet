package ipam

import (
	"context"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connectioncontext"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/networkservice"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"
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
