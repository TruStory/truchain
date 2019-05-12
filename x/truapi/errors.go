package truapi

import "errors"

// Errors for truapi module.
var (
	ErrFlaggedStoryEnvVarParsing = errors.New("Error parsing flagged story environment variable")
	Err400MissingParameter       = errors.New("Missing parameter")
	Err401NotAuthenticated       = errors.New("User not authenticated")
	Err404ResourceNotFound       = errors.New("Resource not found")
	Err422UnprocessableEntity    = errors.New("Unprocessable entity")
	Err500InternalServerError    = errors.New("Something went wrong")
)
