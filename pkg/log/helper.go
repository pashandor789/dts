package log

import (
	"time"
)

func IfError(logger Logger, err error, msg string, args ...interface{}) error {
	if err != nil {
		logger.WithError(err).Errorf(msg, args...)
	}
	return err
}

func DurationToMs(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / float64(time.Millisecond)
}

func Elapsed(logger Logger, startedAt time.Time, msg string, args ...interface{}) {
	logger.WithPrefix(LibPrefix).WithField("elapsed_ms", elapsedMs(startedAt)).Infof(msg, args...)
}

func elapsedMs(since time.Time) float64 {
	return DurationToMs(time.Since(since))
}
