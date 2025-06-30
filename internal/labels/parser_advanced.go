package labels

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/offen/docker-volume-backup/internal/errwrap"
)

// parseAdvancedLabels updates the given Config with additional values parsed
// from labels. Supported keys are:
//   - gpg-passphrase
//   - gpg-public-key-ring
//   - age-passphrase
//   - age-public-keys
//   - stop-during-backup
//   - notification-urls
//   - notification-level
//
// Unknown keys are ignored.
func parseAdvancedLabels(labels map[string]string, c *Config) error {
	trimmed := map[string]string{}
	for key, value := range labels {
		if strings.HasPrefix(key, Prefix) {
			trimmed[strings.TrimPrefix(key, Prefix)] = value
			continue
		}
		trimmed[key] = value
	}

	if v, ok := trimmed["gpg-passphrase"]; ok {
		c.GpgPassphrase = v
	}

	if v, ok := trimmed["gpg-public-key-ring"]; ok {
		c.GpgPublicKeyRing = v
	}

	if v, ok := trimmed["age-passphrase"]; ok {
		c.AgePassphrase = v
	}

	if v, ok := trimmed["age-public-keys"]; ok {
		if v != "" {
			keys := strings.Split(v, ",")
			c.AgePublicKeys = c.AgePublicKeys[:0]
			for _, k := range keys {
				k = strings.TrimSpace(k)
				if k != "" {
					c.AgePublicKeys = append(c.AgePublicKeys, k)
				}
			}
		} else {
			c.AgePublicKeys = nil
		}
	}

	if v, ok := trimmed["stop-during-backup"]; ok {
		c.BackupStopDuringBackupLabel = v
	}

	if v, ok := trimmed["notification-urls"]; ok {
		if v == "" {
			c.NotificationURLs = nil
		} else {
			urls := strings.Split(v, ",")
			c.NotificationURLs = c.NotificationURLs[:0]
			for _, u := range urls {
				u = strings.TrimSpace(u)
				if u == "" {
					continue
				}
				if _, err := url.ParseRequestURI(u); err != nil {
					return errwrap.Wrap(err, fmt.Sprintf("invalid notification url %s", u))
				}
				c.NotificationURLs = append(c.NotificationURLs, u)
			}
		}
	}

	if v, ok := trimmed["notification-level"]; ok {
		lower := strings.ToLower(v)
		switch lower {
		case "", "info", "error":
			c.NotificationLevel = lower
		default:
			return errwrap.Wrap(nil, fmt.Sprintf("invalid notification level %s", v))
		}
	}

	if v, ok := trimmed["smtp-port"]; ok {
		port, err := strconv.Atoi(v)
		if err != nil {
			return errwrap.Wrap(err, fmt.Sprintf("invalid smtp-port value %s", v))
		}
		c.EmailSMTPPort = port
	}

	if v, ok := trimmed["email-recipient"]; ok {
		c.EmailNotificationRecipient = v
	}

	if v, ok := trimmed["email-sender"]; ok {
		c.EmailNotificationSender = v
	}

	if v, ok := trimmed["smtp-host"]; ok {
		c.EmailSMTPHost = v
	}

	if v, ok := trimmed["smtp-username"]; ok {
		c.EmailSMTPUsername = v
	}

	if v, ok := trimmed["smtp-password"]; ok {
		c.EmailSMTPPassword = v
	}

	return nil
}
