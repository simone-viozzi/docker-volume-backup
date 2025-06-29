// Copyright 2024 - offen.software <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package labels

import (
	"context"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/offen/docker-volume-backup/internal/errwrap"
)

func newDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errwrap.Wrap(err, "error creating docker client")
	}
	return cli, nil
}

func listVolumes(cli interface {
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
}) ([]*volume.Volume, error) {
	resp, err := cli.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return nil, errwrap.Wrap(err, "error listing volumes")
	}
	return resp.Volumes, nil
}

// ListVolumes returns all Docker volumes using a new Docker client derived from
// the current environment.
func ListVolumes() ([]*volume.Volume, error) {
	cli, err := newDockerClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	return listVolumes(cli)
}
