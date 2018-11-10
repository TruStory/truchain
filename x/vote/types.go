package vote

// NOTE: the main type for vote lives in the top-level `types` package.
// This is because it is embedded in multiple other types like
// `Backing` and `Challenge` in other modules. So it should exist
// at a higher level than this module.

// Params holds default parameters for a vote
type Params struct {
	MinCommentLength int // min number of chars for argument
	MaxCommentLength int // max number of chars for argument
	MinEvidenceCount int // min number of evidence URLs
	MaxEvidenceCount int // max number of evidence URLs
}

// DefaultParams creates a new Params type with defaults
func DefaultParams() Params {
	return Params{
		MinCommentLength: 10,
		MaxCommentLength: 340,
		MinEvidenceCount: 0,
		MaxEvidenceCount: 10,
	}
}
