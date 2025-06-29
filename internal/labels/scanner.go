package labels

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/volume"
	"github.com/offen/docker-volume-backup/internal/errwrap"
)

func scanVolumeLabels(cli interface {
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
}) (map[string]map[string]string, error) {
	vols, err := listVolumes(cli)
	if err != nil {
		return nil, errwrap.Wrap(err, "error scanning volume labels")
	}
	result := map[string]map[string]string{}
	for _, v := range vols {
		labels := map[string]string{}
		for key, value := range v.Labels {
			if strings.HasPrefix(key, Prefix) {
				labels[strings.TrimPrefix(key, Prefix)] = value
			}
		}
		if len(labels) > 0 {
			result[v.Name] = labels
		}
	}
	return result, nil
}

// ScanVolumeLabels returns a mapping of volume names to their relevant labels.
// The provided client is used to look up the list of volumes.
func ScanVolumeLabels(cli interface {
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
}) (map[string]map[string]string, error) {
	return scanVolumeLabels(cli)
}
