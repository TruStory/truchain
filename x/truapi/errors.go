package truapi

import "errors"

// Errors for truapi module.
var (
	ErrFlaggedStoryEnvVar        = errors.New("Flagged Story Limit environment variable not set")
	ErrFlaggedStoryEnvVarParsing = errors.New("Error parsing flagged story environment variable")
)
