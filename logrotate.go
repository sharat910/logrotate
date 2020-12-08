package logrotator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	filePath          string
	prependTimeFormat string
	header            []byte
	startHour         int
	rotateInterval    time.Duration
	compress          bool
	nowFunc           func() time.Time
}

// LogRotator rotates log files every (configured) time interval
// It implements io.WriteCloser and checks for rotation on every Write call
type LogRotator struct {
	config
	w        *bufio.Writer
	f        *os.File
	logStart time.Time
}

func defaultLogRotator(filePath string) *LogRotator {
	return &LogRotator{
		config: config{
			filePath:          filePath,
			prependTimeFormat: "2006-01-02",
			header:            nil,
			rotateInterval:    24 * time.Hour,
			compress:          false,
			nowFunc:           time.Now,
		},
	}
}

// New returns a LogRotator writing to files of fmt /your/path/<time-fmt>-yourfilename.log
// @Input filePath: General filePath for log files. E.g. /your/path/yourfilename.log
// @Input options: Provide configuration options while instantiating a LogRotator
func New(filePath string, options ...OptFunc) (*LogRotator, error) {
	lr := defaultLogRotator(filePath)
	for _, opt := range options {
		err := opt(lr)
		if err != nil {
			return nil, err
		}
	}
	err := lr.rotate()
	if err != nil {
		return nil, err
	}
	return lr, nil
}

// Write checks if file should be rotated and writes to a bufio.Writer
func (lr *LogRotator) Write(p []byte) (n int, err error) {
	if lr.shouldRotate() {
		err := lr.rotate()
		if err != nil {
			return 0, err
		}
	}
	return lr.w.Write(p)
}

// Flush flushes the underlying bufio.Writer
func (lr *LogRotator) Flush() error {
	return lr.w.Flush()
}

// Close flushes bufio.Writer and closes the log file
func (lr *LogRotator) Close() error {
	// Close existing writer
	if lr.w != nil {
		err := lr.w.Flush()
		if err != nil {
			return err
		}
	}

	// Close existing file
	if lr.f != nil {
		err := lr.f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (lr *LogRotator) initLogStart() {
	now := lr.nowFunc()
	lr.logStart = time.Date(now.Year(), now.Month(), now.Day(), lr.startHour, 0, 0, 0, now.Location())
	for lr.logStart.After(now) {
		lr.logStart = lr.logStart.Add(-lr.rotateInterval)
	}
	for lr.logStart.Add(lr.rotateInterval).Before(now) {
		lr.logStart = lr.logStart.Add(lr.rotateInterval)
	}
}

func (lr *LogRotator) shouldRotate() bool {
	return lr.nowFunc().Sub(lr.logStart) > lr.rotateInterval
}

func (lr *LogRotator) rotate() error {

	if lr.logStart.IsZero() {
		lr.initLogStart()
	} else {
		lr.logStart = lr.logStart.Add(lr.rotateInterval)
	}

	if err := lr.Close(); err != nil {
		return err
	}

	fp, err := lr.getFormattedFilepath(lr.logStart)
	if err != nil {
		return err
	}

	appendMode := fileExists(fp)
	lr.f, err = createFile(fp)
	if err != nil {
		return err
	}
	lr.w = bufio.NewWriter(lr.f)

	if lr.header != nil && !appendMode {
		_, err := lr.w.Write(lr.header)
		if err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) getFormattedFilepath(t time.Time) (string, error) {
	last := filepath.Base(lr.filePath)
	if last == "." {
		return "", errors.New("unable to get filename from filepath")
	}
	prependTimeString := t.Format(lr.prependTimeFormat)
	fp := filepath.Join(filepath.Dir(lr.filePath), fmt.Sprintf("%s-%s", prependTimeString, last))
	return fp, nil
}
