package collector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/linode/linodego"
)

func GetLinodeRegions(client linodego.Client, ctx context.Context) ([]string, error) {
	regions, err := client.ListRegions(ctx, nil)
	if err != nil {
		return nil, err
	}

	regionIDs := make([]string, len(regions))
	for i, region := range regions {
		regionIDs[i] = region.ID
	}
	return regionIDs, nil
}
func GetLinodeObjectStorageEndpoints(client linodego.Client, ctx context.Context) ([]linodego.ObjectStorageEndpoint, error) {

	objectStorageRegions, err := client.ListObjectStorageEndpoints(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Filter out endpoints where S3Endpoint is nil (Not possible with Linode API)
	filtered := make([]linodego.ObjectStorageEndpoint, 0, len(objectStorageRegions))
	for _, ep := range objectStorageRegions {
		if ep.S3Endpoint != nil {
			filtered = append(filtered, ep)
		}
	}

	// Print the object storage regions as json
	jsonData, err := json.Marshal(filtered)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jsonData))
	return filtered, nil
}
