module github.com/danielvladco/k8s-vnet

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/ligato/vpp-agent v2.3.0+incompatible
	github.com/networkservicemesh/networkservicemesh/controlplane/api v0.2.0
	github.com/networkservicemesh/networkservicemesh/pkg v0.2.0
	github.com/networkservicemesh/networkservicemesh/sdk v0.2.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	google.golang.org/grpc v1.23.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	gotest.tools v2.2.0+incompatible
)

replace (
	github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
	github.com/networkservicemesh/networkservicemesh => github.com/networkservicemesh/networkservicemesh v0.0.0-20191110020756-e6a23cc4eff7
	github.com/networkservicemesh/networkservicemesh/controlplane/api => github.com/networkservicemesh/networkservicemesh/controlplane/api v0.0.0-20191110020756-e6a23cc4eff7
	github.com/networkservicemesh/networkservicemesh/pkg => github.com/networkservicemesh/networkservicemesh/pkg v0.0.0-20191110020756-e6a23cc4eff7
	github.com/networkservicemesh/networkservicemesh/sdk => github.com/tiswanso/networkservicemesh/sdk v0.0.0-20200314223159-61ce0a4a8ec9
	github.com/networkservicemesh/networkservicemesh/utils => github.com/networkservicemesh/networkservicemesh/utils v0.0.0-20191110020756-e6a23cc4eff7
)
