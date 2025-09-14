// Copyright 2021-2022 - offen.software <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
)

func main() {
	foreground := flag.Bool("foreground", false, "run the tool in the foreground")
	profile := flag.String("profile", "", "collect runtime metrics and log them periodically on the given cron expression")
	configStyle := flag.String("config-style", "envfile", "load configuration from envfile (default) or labels")
	flag.Parse()

	c := newCommand()
	if *foreground {
		opts := foregroundOpts{
			profileCronExpression: *profile,
		}
		c.must(c.runInForeground(opts, *configStyle))
	} else {
		c.must(c.runAsCommand(*configStyle))
	}
}
