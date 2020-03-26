package ipam

import (
	"fmt"
	"net"

	"github.com/networkservicemesh/networkservicemesh/sdk/common"
)

type PrefixPoolGenerator func(cfg *common.NSConfiguration) (string, error)

func PrefixPoolFromPodIP(ipamUniqueOctet int) PrefixPoolGenerator {
	return func(nsConfig *common.NSConfiguration) (s string, err error) {
		if nsConfig.IPAddress == "" {
			return "", fmt.Errorf("NSConfiguration.IPAddress is empty")
		}

		prefixPoolIP, _, err := net.ParseCIDR(nsConfig.IPAddress)
		if err != nil {
			return "", fmt.Errorf("failed to parse configured prefix pool ip: %w", err)
		}

		return fmt.Sprintf("%d.%d.%d.%d/24", prefixPoolIP.To4()[0], prefixPoolIP.To4()[1], ipamUniqueOctet, 0), nil
	}
}
