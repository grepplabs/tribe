package errors

import "github.com/pkg/errors"

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func WithStack(err error) error {
	if _, ok := err.(StackTracer); ok {
		return err
	}

	return errors.WithStack(err)
}
