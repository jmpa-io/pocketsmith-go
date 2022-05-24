package pocketsmith

import (
	"fmt"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		// should never reach here.
		return fmt.Sprintf("missing implementation: %d", int(l))
	}
}

// StringToLogLevel converts a string to a log level for this package.
// Returns an error if a log level is not found.
func StringToLogLevel(s string) (LogLevel, error) {
	switch s {
	case LogLevelDebug.String():
		return LogLevelDebug, nil
	case LogLevelInfo.String():
		return LogLevelInfo, nil
	case LogLevelWarn.String():
		return LogLevelWarn, nil
	case LogLevelError.String():
		return LogLevelError, nil
	}
	return LogLevelDebug, ErrFailedFindingLogLevel{s}
}
