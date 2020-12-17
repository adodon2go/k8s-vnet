package vppagent

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ligato/vpp-agent/api/models/vpp"
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connection"
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/networkservice"
	"cisco-app-networking.github.io/networkservicemesh/sdk/common"
	"cisco-app-networking.github.io/networkservicemesh/sdk/endpoint"
	"github.com/sirupsen/logrus"

	"github.com/adodon2go/k8s-vnet/internal/cnf"
)

// UniversalCNFEndpoint is a Universal CNF Endpoint composite implementation
type vppEndpoint struct {
	backend     Service
	serviceName string
	ifname      string
	dpConfig    *vpp.ConfigData
}

// Request implements the request handler
func (uce *vppEndpoint) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	conn := request.GetConnection()

	if uce.dpConfig == nil {
		uce.dpConfig = &vpp.ConfigData{}
	}

	err := uce.backend.ProcessEndpointDP(ctx, &ProcessDataPlaneReq{
		Vppconfig:   uce.dpConfig,
		ServiceName: uce.serviceName,
		Ifname:      uce.ifname,
		Connection:  conn,
	})
	if err != nil {
		logrus.Errorf("Error processing dpconfig: %+v", uce.dpConfig)
		return nil, err
	}

	if endpoint.Next(ctx) != nil {
		return endpoint.Next(ctx).Request(ctx, request)
	}

	return request.GetConnection(), nil
}

// Close implements the close handler
func (uce *vppEndpoint) Close(ctx context.Context, connection *connection.Connection) (*empty.Empty, error) {
	logrus.Infof("VPP CNF DeleteConnection: %v", connection)

	if endpoint.Next(ctx) != nil {
		return endpoint.Next(ctx).Close(ctx, connection)
	}

	return &empty.Empty{}, nil
}

// Name returns the composite name
func (uce *vppEndpoint) Name() string {
	return "VPP endpoint"
}

func MakeNewVPPEndpoint(service Service) cnf.CompositeEndpointFactory {
	return func(cfg *common.NSConfiguration, _ *string) (server networkservice.NetworkServiceServer, err error) {
		return &vppEndpoint{
			serviceName: cfg.EndpointNetworkService,
			ifname:      cfg.NscInterfaceName,
			backend:     service,
		}, nil
	}
}
