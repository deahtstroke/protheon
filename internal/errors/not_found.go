package protheonErrors

import "fmt"

type NotFoundError struct {
	Identifier string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("No config found with Id or alias '%s'", e.Identifier)
}
