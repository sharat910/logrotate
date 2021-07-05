package logrotator

import (
	"errors"
	"time"
)

// OptFunc sets configuration options for LogRotator
type OptFunc func(lr *LogRotator) error

// WithImmediateFlush enables flushing to file on every Write (useful for logs)
func WithImmediateFlush(lr *LogRotator) error {
	lr.immediateFlush = true
	return nil
}

// WithCompression enables compression on closed log files
func WithCompression(lr *LogRotator) error {
	lr.compress = true
	return nil
}

// PrependTimeFormat sets the time fmt string and delim which is added as prefix to filename
func PrependTimeFormat(tf string, delim string) OptFunc {
	return func(lr *LogRotator) error {
		lr.prependTimeFormat = tf + delim
		return nil
	}
}

// StartHour sets the first hour to start rotating log files from
func StartHour(h int) OptFunc {
	return func(lr *LogRotator) error {
		if h <= 0 || h >= 23 {
			return errors.New("invalid startHour: valid range: [0,23]")
		}
		lr.startHour = h
		return nil
	}
}

// RotateInterval sets the max interval after which log files will be rotated
func RotateInterval(i time.Duration) OptFunc {
	return func(lr *LogRotator) error {
		if i < 0 {
			return errors.New("invalid rotateInterval: negative duration")
		}
		lr.rotateInterval = i
		return nil
	}
}

// Header sets a byte slice that will be written on every new log file creation
// Useful while rotating csv files.
func Header(h []byte) OptFunc {
	return func(lr *LogRotator) error {
		lr.header = h
		return nil
	}
}

func RotateCallback(cb func(t time.Time)) OptFunc {
	return func(lr *LogRotator) error {
		lr.rotateCallback = cb
		return nil
	}
}
