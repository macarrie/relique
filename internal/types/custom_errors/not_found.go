package custom_errors

import "fmt"

type DBNotFoundError struct {
	ID       int64
	ItemType string
}

func (e *DBNotFoundError) Error() string {
	return fmt.Sprintf("no item '%s' with ID '%d' found in DB", e.ItemType, e.ID)
}

func IsDBNotFoundError(e error) bool {
	_, ok := e.(*DBNotFoundError)
	return ok
}
