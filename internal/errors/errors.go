package errors

import (
	"errors"
	"fmt"
)

func New(location, msg string) error {
	return fmt.Errorf("[%s] %s", location, msg)
}

func Newf(location, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("[%s] %s", location, msg)
}

func Wrap(err error, location, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("[%s] %s: %w", location, msg, err)
}

func Wrapf(err error, location, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("[%s] %s: %w", location, fmt.Sprintf(format, args...), err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
