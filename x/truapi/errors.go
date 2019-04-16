package truapi

import "errors"

// Errors for truapi module.
var (
	ErrFlaggedStoryEnvVarParsing = errors.New("Error parsing flagged story environment variable")
	Err404                       = errors.New("Resource not found")
	Err401NotAuthenticated       = errors.New("User not authenticated")
)
