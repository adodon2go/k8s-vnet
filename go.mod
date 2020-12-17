module github.com/adodon2go/k8s-vnet

go 1.13

require (
	cisco-app-networking.github.io/networkservicemesh/controlplane/api v1.0.10
	cisco-app-networking.github.io/networkservicemesh/pkg v1.0.10
	cisco-app-networking.github.io/networkservicemesh/sdk v1.0.10
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/ligato/cn-infra v2.2.0+incompatible // indirect
	github.com/ligato/vpp-agent v2.3.0+incompatible
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	gotest.tools v2.2.0+incompatible
)

replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
