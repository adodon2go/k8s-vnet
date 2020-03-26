package main

import (
	"context"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/ligato/vpp-agent/api/configurator"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/networkservice"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/danielvladco/k8s-vnet/internal/cnf"
	"github.com/danielvladco/k8s-vnet/internal/ipam"
	"github.com/danielvladco/k8s-vnet/internal/vl3"
	"github.com/danielvladco/k8s-vnet/internal/vppagent"
	"github.com/danielvladco/k8s-vnet/pkg/nseconfig"
)

var (
	configPath        = flag.String("file", "/etc/vl3-nse/config.yaml", " full path to the configuration file")
	nsRegAddr         = envOrDefault("NSREGISTRY_ADDR", "nsmgr.nsm-system:5000")
	vppAgentEndpoint  = envOrDefault("VPP_AGENT_ENDPOINT", "localhost:9113")
	nseUniqueOctetStr = envOrDefault("NSE_IPAM_UNIQUE_OCTET", "")
	workspace         = envOrDefault(common.WorkspaceEnv, "/")
	nsRemoteIpList    = strings.Split(envOrDefault("NSM_REMOTE_NS_IP_LIST", ""), ",")
	logger            = logrus.New()
)

func init() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
}

func main() {
	flag.Parse()

	// basic configuration initialization
	file, err := os.Open(*configPath)
	fatalOnErr(err, "unable to locate config file")
	cnfConfig := &nseconfig.Config{}
	err = nseconfig.NewConfig(yaml.NewDecoder(file), cnfConfig)
	fatalOnErr(err, "config creation error")

	ctx := context.Background()

	err = tools.WaitForPortAvailable(ctx, "tcp", vppAgentEndpoint, 100*time.Millisecond)
	fatalOnErr(err, "error waiting for port")

	// grpc tracing and instrumenting
	tracer := opentracing.GlobalTracer()
	grpcOptions := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	}

	// vpp agent definition
	vppAgentGrpcConn, err := grpc.Dial(vppAgentEndpoint, grpcOptions...)
	fatalOnErr(err, "can't dial grpc server")
	defer logCloser(vppAgentGrpcConn.Close)
	var vppAgentConfiguratorClient = configurator.NewConfiguratorClient(vppAgentGrpcConn)
	vppagentSvc, err := vppagent.NewService(vppAgentConfiguratorClient, workspace)
	fatalOnErr(err, "vppagent initialization failed")

	// network service discovery client definition
	nsRegGrpcConn, err := tools.DialTCP(nsRegAddr, grpcOptions...)
	fatalOnErr(err, "unable to connect to discovery client")
	defer logCloser(nsRegGrpcConn.Close)
	var nsDiscoveryClient = registry.NewNetworkServiceDiscoveryClient(nsRegGrpcConn)

	nseUniqueOctet, err := strconv.Atoi(nseUniqueOctetStr)
	fatalOnErr(err, "NSE_IPAM_UNIQUE_OCTET env var is invalid")

	ipamGen := ipam.PrefixPoolFromPodIP(nseUniqueOctet)
	cleanup, err := cnf.InitAndStartNSEndpoints(cnfConfig.Endpoints,
		newMonitorEndpoint,
		newConnectionEndpoint,
		ipam.MakeNewIpamEndpoint(ipamGen),
		vl3.MakeNewVL3Endpoint(ipamGen, vppagentSvc, nsRemoteIpList, nsDiscoveryClient),
		ipam.NewRouteEndpoint,
		vppagent.MakeNewVPPEndpoint(vppagentSvc),
	)
	fatalOnErr(err, "endpoints init failed")
	defer logCloser(cleanup)

	logger.Infof("Starting endpoints")

	// Capture signals to cleanup before exiting
	<-tools.NewOSSignalChannel()
}
func newMonitorEndpoint(configuration *common.NSConfiguration, _ *string) (server networkservice.NetworkServiceServer, err error) {
	return endpoint.NewMonitorEndpoint(configuration), nil
}

func newConnectionEndpoint(configuration *common.NSConfiguration, _ *string) (server networkservice.NetworkServiceServer, err error) {
	return endpoint.NewConnectionEndpoint(configuration), nil
}

func fatalOnErr(err error, s string) {
	if err != nil {
		logger.Fatal(s, ": ", err.Error())
	}
}

func envOrDefault(envName, defaultValue string) string {
	val, ok := os.LookupEnv(envName)
	if ok {
		return val
	}

	return defaultValue
}

func logCloser(closer func() error) {
	err := closer()
	if err != nil {
		logger.Error("closed with errors: ", err)
	}
}
