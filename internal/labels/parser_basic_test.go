package labels

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseBasicLabels(t *testing.T) {
	tests := []struct {
		name        string
		labels      map[string]string
		expected    Config
		expectError bool
	}{
		{
			name: "complete",
			labels: map[string]string{
				Prefix + "schedule": "@daily",
				Prefix + "target":   "/archive",
				Prefix + "rotation": "7",
			},
			expected: Config{
				BackupCronExpression: "@daily",
				BackupArchive:        "/archive",
				BackupRetentionDays:  7,
			},
		},
		{
			name: "partial",
			labels: map[string]string{
				Prefix + "schedule": "0 0 * * *",
			},
			expected: Config{
				BackupCronExpression: "0 0 * * *",
			},
		},
		{
			name: "invalid schedule",
			labels: map[string]string{
				Prefix + "schedule": "never",
			},
			expectError: true,
		},
		{
			name: "invalid rotation",
			labels: map[string]string{
				Prefix + "rotation": "abc",
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseBasicLabels(test.labels)
			if (err != nil) != test.expectError {
				t.Fatalf("unexpected error state %v", err)
			}
			if !test.expectError && !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func Example_parseBasicLabels() {
	cfg, _ := ParseBasicLabels(map[string]string{
		Prefix + "schedule": "@hourly",
		Prefix + "target":   "/archive",
		Prefix + "rotation": "3",
	})
	fmt.Println(cfg.BackupCronExpression)
	fmt.Println(cfg.BackupArchive)
	fmt.Println(cfg.BackupRetentionDays)
	// Output:
	// @hourly
	// /archive
	// 3
}
