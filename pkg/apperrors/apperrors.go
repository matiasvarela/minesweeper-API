package apperrors

import (
	"github.com/matiasvarela/errors"
)

var (
	Internal     = errors.Define("internal")
	NotFound     = errors.Define("not_found")
	InvalidInput = errors.Define("invalid_input")
)
