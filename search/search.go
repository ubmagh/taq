package search

import (
	"strings"

	"github.com/ubmagh/taq/types"
)

func SearchHostsByLabel(hosts []types.Host, searchToken string) ([]types.Host, error) {
	var matchedHosts []types.Host
	for _, host := range hosts {
		if strings.Contains(host.Name, searchToken) || strings.Contains(host.Address, searchToken) {
			matchedHosts = append(matchedHosts, host)
			continue
		}

		if val, ok := host.Labels[searchToken]; ok && val == searchToken {
			matchedHosts = append(matchedHosts, host)
		}
	}
	return matchedHosts, nil
}
