module github.com/danielvladco/k8s-vnet

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/ligato/cn-infra v2.2.0+incompatible // indirect
	github.com/ligato/vpp-agent v2.3.0+incompatible
	github.com/networkservicemesh/networkservicemesh/controlplane v0.0.0-20200519133935-dd205487e66d // indirect
	github.com/networkservicemesh/networkservicemesh/controlplane/api v0.3.0
	github.com/networkservicemesh/networkservicemesh/pkg v0.3.0
	github.com/networkservicemesh/networkservicemesh/sdk v0.3.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	gotest.tools v2.2.0+incompatible
)

replace (
	github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
	github.com/networkservicemesh/networkservicemesh => github.com/tiswanso/networkservicemesh v0.0.0-20200515015809-416c5355f322
	github.com/networkservicemesh/networkservicemesh/controlplane/api => github.com/tiswanso/networkservicemesh/controlplane/api v0.0.0-20200515015809-416c5355f322
	github.com/networkservicemesh/networkservicemesh/forwarder/api => github.com/tiswanso/networkservicemesh/forwarder/api v0.0.0-20200515015809-416c5355f322
	github.com/networkservicemesh/networkservicemesh/pkg => github.com/tiswanso/networkservicemesh/pkg v0.0.0-20200515015809-416c5355f322
	github.com/networkservicemesh/networkservicemesh/sdk => github.com/tiswanso/networkservicemesh/sdk v0.0.0-20200515015809-416c5355f322
	github.com/networkservicemesh/networkservicemesh/utils => github.com/tiswanso/networkservicemesh/utils v0.0.0-20200515015809-416c5355f322
)
