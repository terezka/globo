//
// utils.go
// 2016 giulio <giulioungaretti@me.com>
//

package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//LogLevel sets the level of the logging with a flag
var LogLevel = flag.String("log_level", "warn", "pick log level (error|warn|debug)")

// ParseLogLevel parses and returns the desired log level
func ParseLogLevel(logLevel string) (log.Level, error) {
	switch logLevel {
	case "debug":
		return log.DebugLevel, nil
	case "warn":
		return log.WarnLevel, nil
	case "error":
		return log.ErrorLevel, nil
	}

	return log.ErrorLevel, fmt.Errorf("Incorrect log-level setting:  %v", logLevel)
}
