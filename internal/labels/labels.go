// Copyright 2024 - offen.software <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package labels

// Prefix is the namespace used when scanning Docker volumes for configuration labels.
// It can be overridden at runtime using the `--label-prefix` flag or the
// `LABEL_PREFIX` environment variable.
var Prefix = "dvbackup."
