package labels

import (
	"fmt"
	"strconv"

	"github.com/offen/docker-volume-backup/internal/errwrap"
	"github.com/robfig/cron/v3"
)

// Config holds the subset of configuration derived from volume labels.
type Config struct {
	BackupCronExpression string
	BackupArchive        string
	BackupRetentionDays  int32
}

// parseBasicLabels converts a set of volume labels into a Config.
// Supported keys are:
//   - schedule: cron expression controlling when to run a backup
//   - target:   backup archive path
//   - rotation: number of days to keep backups
func parseBasicLabels(labels map[string]string) (Config, error) {
	trimmed := map[string]string{}
	for key, value := range labels {
		if strings.HasPrefix(key, Prefix) {
			trimmed[strings.TrimPrefix(key, Prefix)] = value
			continue
		}
		trimmed[key] = value
	}

	var c Config

	if v, ok := trimmed["schedule"]; ok {
		if _, err := cron.ParseStandard(v); err != nil {
			return c, errwrap.Wrap(err, fmt.Sprintf("invalid schedule value %s", v))
		}
		c.BackupCronExpression = v
	}

	if v, ok := trimmed["target"]; ok {
		c.BackupArchive = v
	}

	if v, ok := trimmed["rotation"]; ok {
		days, err := strconv.Atoi(v)
		if err != nil {
			return c, errwrap.Wrap(err, fmt.Sprintf("invalid rotation value %s", v))
		}
		c.BackupRetentionDays = int32(days)
	}

	return c, nil
}
