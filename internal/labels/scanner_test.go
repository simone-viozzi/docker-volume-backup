package labels

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/volume"
	"github.com/offen/docker-volume-backup/internal/errwrap"
)

func TestScanVolumeLabels(t *testing.T) {
	tests := []struct {
		name        string
		client      *mockVolumeClient
		expected    map[string]map[string]string
		expectError bool
	}{
		{
			name: "success",
			client: &mockVolumeClient{
				result: volume.ListResponse{Volumes: []*volume.Volume{
					{Name: "foo", Labels: map[string]string{Prefix + "schedule": "daily", Prefix + "path": "/data", "ignored": "x"}},
					{Name: "bar", Labels: map[string]string{"other": "y", Prefix + "path": "/var"}},
				}},
			},
			expected: map[string]map[string]string{
				"foo": {"schedule": "daily", "path": "/data"},
				"bar": {"path": "/var"},
			},
		},
		{
			name: "none",
			client: &mockVolumeClient{
				result: volume.ListResponse{Volumes: []*volume.Volume{{Name: "foo"}}},
			},
			expected: map[string]map[string]string{},
		},
		{
			name:        "error",
			client:      &mockVolumeClient{err: errExample},
			expected:    nil,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := scanVolumeLabels(test.client)
			if (err != nil) != test.expectError {
				t.Errorf("unexpected error value %v", err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

var errExample = errwrap.Wrap(nil, "boom")
