package db

import (
	"encoding/hex"
	"fmt"
)

// NotFoundError is returned whenever a model with a specific ID should be found
// in the database but it is not.
type NotFoundError struct {
	ID []byte
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("could not find model with the given ID: %s", hex.EncodeToString(e.ID))
}

// AlreadyExistsError is returned whenever a model with a specific ID should not
// already exists in the database but it does.
type AlreadyExistsError struct {
	ID []byte
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("model already exists with the given ID: %s", hex.EncodeToString(e.ID))
}
