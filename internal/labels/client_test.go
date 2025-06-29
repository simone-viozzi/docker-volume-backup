// Copyright 2024 - offen.software <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package labels

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/volume"
)

type mockVolumeClient struct {
	result volume.ListResponse
	err    error
}

func (m *mockVolumeClient) VolumeList(context.Context, volume.ListOptions) (volume.ListResponse, error) {
	return m.result, m.err
}

func TestListVolumesInternal(t *testing.T) {
	tests := []struct {
		name        string
		client      *mockVolumeClient
		expected    []*volume.Volume
		expectError bool
	}{
		{
			name: "success",
			client: &mockVolumeClient{
				result: volume.ListResponse{Volumes: []*volume.Volume{{Name: "foo"}, {Name: "bar"}}},
			},
			expected: []*volume.Volume{{Name: "foo"}, {Name: "bar"}},
		},
		{
			name:        "error",
			client:      &mockVolumeClient{err: errors.New("explode")},
			expected:    nil,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := listVolumes(test.client)
			if (err != nil) != test.expectError {
				t.Errorf("unexpected error value %v", err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
