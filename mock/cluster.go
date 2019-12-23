package mock

import (
	"context"
)

// LKECluster is a mock implementation of linodego's Kubernetes LKECluster type
type LKECluster struct {
	ID      int `json:"id"`
	Created string
	Updated string
	Label   string
	Region  string
	Tags    []string
	Status  string
	Version string
}

// LKEClusterPool is a mock implementation of linodego's Kubernetes LKEClusterPool type
type LKEClusterPool struct {
	ID      int                    `json:"id"`
	Label   string                 `json:"label"`
	Count   int                    `json:"count"`
	Type    string                 `json:"string"`
	Linodes []LKEClusterPoolLinode `json:"linodes"`
}

// LKEClusterPoolLinode is a mock implementation of linodego's Kubernetes LKEClusterPoolLinode type
type LKEClusterPoolLinode struct {
	ID     *int
	Status LKELinodeStatus
}

// NewLKEClusterPoolLinode returns a new LKEClusterPoolLinode used because the ID is a *int
func NewLKEClusterPoolLinode(ID int) LKEClusterPoolLinode {
	return LKEClusterPoolLinode{
		ID:     &ID,
		Status: "active",
	}
}

// LKELinodeStatus is a mock implementation of linodego's LKELinodeStatus enum
type LKELinodeStatus string

// LKEClusterPoolStatus constants reflect the current status of an LKEClusterPool
const (
	LKELinodeReady    LKELinodeStatus = "ready"
	LKELinodeNotReady LKELinodeStatus = "not_ready"
)

// ListLKEClusters is a mock implementation of linodego's Kubernetes ListLKEClusters function
func (c *Client) ListLKEClusters(ctx context.Context, opts interface{}) ([]LKECluster, error) {
	return []LKECluster{
		LKECluster{
			ID:    660,
			Label: "linode-exporter",
		},
		LKECluster{
			ID:    661,
			Label: "bigmachine",
		},
	}, nil
}

// ListLKEClusterPools is a mock implementation of linodego's Kubernetes ListLKEClusterPools function
func (c *Client) ListLKEClusterPools(ctx context.Context, ID int, opts interface{}) ([]LKEClusterPool, error) {
	return []LKEClusterPool{
		LKEClusterPool{
			ID:    880,
			Count: 3,
			Type:  "g6-standard-1",
			Linodes: []LKEClusterPoolLinode{
				NewLKEClusterPoolLinode(1),
				NewLKEClusterPoolLinode(2),
				NewLKEClusterPoolLinode(3),
			},
		},
	}, nil
}
