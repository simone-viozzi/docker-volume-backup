package labels

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseAdvancedLabels(t *testing.T) {
	tests := []struct {
		name        string
		labels      map[string]string
		start       Config
		expected    Config
		expectError bool
	}{
		{
			name: "complete",
			labels: map[string]string{
				Prefix + "gpg-passphrase":      "secret",
				Prefix + "gpg-public-key-ring": "/keys/pubring.gpg",
				Prefix + "age-passphrase":      "agepass",
				Prefix + "age-public-keys":     "key1, key2",
				Prefix + "stop-during-backup":  "db",
				Prefix + "notification-urls":   "https://example.com,a://b",
				Prefix + "notification-level":  "INFO",
				Prefix + "smtp-port":           "2525",
				Prefix + "email-recipient":     "a@example.com",
				Prefix + "email-sender":        "b@example.com",
				Prefix + "smtp-host":           "mail.example.com",
				Prefix + "smtp-username":       "smtpuser",
				Prefix + "smtp-password":       "smtppass",
			},
			expected: Config{
				GpgPassphrase:               "secret",
				GpgPublicKeyRing:            "/keys/pubring.gpg",
				AgePassphrase:               "agepass",
				AgePublicKeys:               []string{"key1", "key2"},
				BackupStopDuringBackupLabel: "db",
				NotificationURLs:            []string{"https://example.com", "a://b"},
				NotificationLevel:           "info",
				EmailSMTPPort:               2525,
				EmailNotificationRecipient:  "a@example.com",
				EmailNotificationSender:     "b@example.com",
				EmailSMTPHost:               "mail.example.com",
				EmailSMTPUsername:           "smtpuser",
				EmailSMTPPassword:           "smtppass",
			},
		},
		{
			name: "invalid url",
			labels: map[string]string{
				Prefix + "notification-urls": "ht tp://wrong",
			},
			expectError: true,
		},
		{
			name: "invalid level",
			labels: map[string]string{
				"notification-level": "warn",
			},
			expectError: true,
		},
		{
			name: "invalid smtp port",
			labels: map[string]string{
				"smtp-port": "x",
			},
			expectError: true,
		},
		{
			name: "without prefix",
			labels: map[string]string{
				"gpg-passphrase":     "abc",
				"notification-level": "ERROR",
				"smtp-port":          "80",
			},
			expected: Config{
				GpgPassphrase:     "abc",
				NotificationLevel: "error",
				EmailSMTPPort:     80,
			},
		},
		{
			name: "empty lists",
			labels: map[string]string{
				Prefix + "age-public-keys":   "",
				Prefix + "notification-urls": "",
			},
			expected: Config{},
		},
		{
			name: "whitespace lists",
			labels: map[string]string{
				Prefix + "age-public-keys":   " , ",
				Prefix + "notification-urls": " ,",
			},
			expected: Config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := test.start
			err := ParseAdvancedLabels(test.labels, &cfg)
			if (err != nil) != test.expectError {
				t.Fatalf("unexpected error state %v", err)
			}
			if !test.expectError && !reflect.DeepEqual(cfg, test.expected) {
				t.Errorf("expected %#v, got %#v", test.expected, cfg)
			}
		})
	}
}

func Example_parseAdvancedLabels() {
	cfg := Config{}
	_ = ParseAdvancedLabels(map[string]string{
		Prefix + "age-public-keys":   "key1,key2",
		Prefix + "notification-urls": "https://example.com",
	}, &cfg)
	fmt.Println(cfg.AgePublicKeys)
	fmt.Println(cfg.NotificationURLs)
	// Output:
	// [key1 key2]
	// [https://example.com]
}
