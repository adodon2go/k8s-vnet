package ipam

import (
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/networkservice"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"

	"github.com/danielvladco/k8s-vnet/internal/cnf"
)

func MakeNewIpamEndpoint(ipamCidrGen PrefixPoolGenerator) cnf.CompositeEndpointFactory {
	return func(nsConfig *common.NSConfiguration, _ *string) (networkservice.NetworkServiceServer, error) {
		ipamCidr, err := ipamCidrGen(nsConfig)
		if err != nil {
			return nil, err
		}

		return endpoint.NewIpamEndpoint(&common.NSConfiguration{
			IPAddress: ipamCidr,
		}), nil
	}
}
