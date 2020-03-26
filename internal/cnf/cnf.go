package cnf

import (
	"bytes"
	"context"
	"fmt"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/networkservice"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"
	"github.com/sirupsen/logrus"

	"github.com/danielvladco/k8s-vnet/pkg/nseconfig"
)

type CompositeEndpointFactory func(cfg *common.NSConfiguration, serviceEndpointName *string) (networkservice.NetworkServiceServer, error)

// NewProcessEndpoints returns a new ProcessInitCommands struct
func InitAndStartNSEndpoints(endpoints []*nseconfig.Endpoint, endpointFactories ...CompositeEndpointFactory) (cleaner, error) {
	var cleaners []cleaner
	for _, e := range endpoints {
		configuration := e.ToNSConfiguration()
		// Build the list of composites
		var compositeEndpoints []networkservice.NetworkServiceServer
		// Invoke any additional composite endpoint constructors via the add-on interface
		var serviceEndpointName *string
		for _, addon := range endpointFactories {
			addCompositeEndpoints, err := addon(configuration, serviceEndpointName)
			if err != nil {
				return nil, err
			}
			if addCompositeEndpoints != nil {
				compositeEndpoints = append(compositeEndpoints, addCompositeEndpoints)
			}
		}

		clean, err := StartNSEndpoint(configuration, serviceEndpointName, compositeEndpoints...)
		if err != nil {
			return nil, err
		}

		cleaners = append(cleaners, clean)
	}

	return func() error {
		var errs errors
		for _, clean := range cleaners {
			if err := clean(); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return errs
		}
		return nil
	}, nil
}

// Process iterates over the init commands and applies them
func StartNSEndpoint(NSConfiguration *common.NSConfiguration, serviceEndpointName *string, compositeEndpoints ...networkservice.NetworkServiceServer) (cleaner, error) {
	nsEndpoint, err := endpoint.NewNSMEndpoint(context.TODO(), NSConfiguration, endpoint.NewCompositeEndpoint(compositeEndpoints...))
	if err != nil {
		return nil, err
	}
	if err = nsEndpoint.Start(); err != nil {
		return nil, err
	}
	v := nsEndpoint.GetName()
	serviceEndpointName = &v
	logrus.Infof("Started endpoint %s", nsEndpoint.GetName())
	return nsEndpoint.Delete, nil
}

type errors []error

func (v errors) Error() string {
	b := bytes.NewBufferString("errors: \n")
	for _, err := range v {
		fmt.Fprintf(b, "\t%s\n", err)
	}

	return b.String()
}

type cleaner func() error
