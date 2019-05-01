package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComment(t *testing.T) {
	testComment := "@d-truth, @shanev you hear me?"
	mentions := parseMentions(testComment)
	assert.Equal(t, 2, len(mentions))
}

func TestCommentWithNewline(t *testing.T) {
	testComment := "@d-truth\n@shanev you hear me?"
	mentions := parseMentions(testComment)
	assert.Equal(t, 2, len(mentions))
}
