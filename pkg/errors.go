package pkg

import "fmt"

type ErrIllegalArgument struct {
	Reason string
}

func (e ErrIllegalArgument) Error() string {
	return fmt.Sprintf("Illegal argument: %q", e.Reason)
}
