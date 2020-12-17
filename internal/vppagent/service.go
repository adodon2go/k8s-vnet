package vppagent

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"

	"github.com/ligato/vpp-agent/api/configurator"
	"github.com/ligato/vpp-agent/api/models/vpp"
	interfaces "github.com/ligato/vpp-agent/api/models/vpp/interfaces"
	vpp_l3 "github.com/ligato/vpp-agent/api/models/vpp/l3"
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connection"
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connection/mechanisms/memif"
	"github.com/sirupsen/logrus"
)

type Service interface {
	ProcessClientDP(ctx context.Context, req *ProcessDataPlaneReq) error
	ProcessEndpointDP(ctx context.Context, req *ProcessDataPlaneReq) error
}

func NewService(client configurator.ConfiguratorClient, workspace string) (Service, error) {
	if _, err := client.Update(context.Background(), &configurator.UpdateRequest{
		Update:     &configurator.Config{},
		FullResync: true,
	}); err != nil {
		return nil, fmt.Errorf("failed to reset vppagent: %w", err)
	}

	return service{client: client, workspace: workspace, endpointIfID: map[string]int{}}, nil
}

type service struct {
	client       configurator.ConfiguratorClient
	endpointIfID map[string]int
	workspace    string
}

func (s service) ProcessClientDP(ctx context.Context, req *ProcessDataPlaneReq) error {
	srcIP := req.Connection.GetContext().GetIpContext().GetSrcIpAddr()
	dstIP, _, _ := net.ParseCIDR(req.Connection.GetContext().GetIpContext().GetDstIpAddr())
	socketFilename := path.Join(s.workspace, memif.ToMechanism(req.Connection.GetMechanism()).GetSocketFilename())

	var ipAddresses []string
	if len(srcIP) > 4 {
		ipAddresses = append(ipAddresses, srcIP)
	}

	req.Vppconfig.Interfaces = append(req.Vppconfig.Interfaces,
		&interfaces.Interface{
			Name:        req.Ifname,
			Type:        interfaces.Interface_MEMIF,
			Enabled:     true,
			IpAddresses: ipAddresses,
			Link: &interfaces.Interface_Memif{
				Memif: &interfaces.MemifLink{
					Master:         false, // The client is not the master in MEMIF
					SocketFilename: socketFilename,
				},
			},
		})

	// Process static routes
	for _, route := range req.Connection.GetContext().GetIpContext().GetDstRoutes() {
		route := &vpp.Route{
			Type:        vpp_l3.Route_INTER_VRF,
			DstNetwork:  route.Prefix,
			NextHopAddr: dstIP.String(),
		}
		req.Vppconfig.Routes = append(req.Vppconfig.Routes, route)
	}

	return s.send(ctx, req.Vppconfig)
}

type ProcessDataPlaneReq struct {
	Vppconfig   *vpp.ConfigData
	ServiceName string
	Ifname      string
	Connection  *connection.Connection
}

func (s service) ProcessEndpointDP(ctx context.Context, req *ProcessDataPlaneReq) error {
	srcIP, _, _ := net.ParseCIDR(req.Connection.GetContext().GetIpContext().GetSrcIpAddr())
	dstIP := req.Connection.GetContext().GetIpContext().GetDstIpAddr()
	socketFilename := path.Join(s.workspace, memif.ToMechanism(req.Connection.GetMechanism()).GetSocketFilename())

	var ipAddresses []string
	if len(dstIP) > 4 {
		ipAddresses = append(ipAddresses, dstIP)
	}

	req.Vppconfig.Interfaces = append(req.Vppconfig.Interfaces,
		&interfaces.Interface{
			Name:        req.Ifname + s.getEndpointIfID(req.ServiceName),
			Type:        interfaces.Interface_MEMIF,
			Enabled:     true,
			IpAddresses: ipAddresses,
			Link: &interfaces.Interface_Memif{
				Memif: &interfaces.MemifLink{
					Master:         true, // The endpoint is always the master in MEMIF
					SocketFilename: socketFilename,
				},
			},
		})

	if err := os.MkdirAll(path.Dir(socketFilename), os.ModePerm); err != nil {
		return err
	}

	// Process static routes
	for _, route := range req.Connection.GetContext().GetIpContext().GetSrcRoutes() {
		route := &vpp.Route{
			Type:        vpp_l3.Route_INTER_VRF,
			DstNetwork:  route.Prefix,
			NextHopAddr: srcIP.String(),
		}
		req.Vppconfig.Routes = append(req.Vppconfig.Routes, route)
	}

	return s.send(ctx, req.Vppconfig)
}

func (s service) send(ctx context.Context, vppconfig *vpp.ConfigData) error {
	_, err := s.client.Update(ctx, &configurator.UpdateRequest{
		Update: &configurator.Config{
			VppConfig: vppconfig,
		},
	})
	if err != nil {
		logrus.Error(err)
		_, err = s.client.Delete(ctx, &configurator.DeleteRequest{
			Delete: &configurator.Config{
				VppConfig: vppconfig,
			},
		})
	}
	return err
}

// GetEndpointIfID generates a new interface ID from the service name
func (s *service) getEndpointIfID(serviceName string) string {
	if _, ok := s.endpointIfID[serviceName]; !ok {
		s.endpointIfID[serviceName] = 0
	} else {
		s.endpointIfID[serviceName]++
	}

	return "/" + strconv.Itoa(s.endpointIfID[serviceName])
}
