package ipam

import (
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/networkservice"
	"cisco-app-networking.github.io/networkservicemesh/sdk/common"
	"cisco-app-networking.github.io/networkservicemesh/sdk/endpoint"

	"github.com/adodon2go/k8s-vnet/internal/cnf"
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
