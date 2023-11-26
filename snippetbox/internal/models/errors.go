package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrValidation = errors.New("Validation error")
