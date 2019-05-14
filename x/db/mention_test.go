package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMentionParse(t *testing.T) {
	text := "@d-truth, @shanev you hear me?"
	mentions := parseMentions(text)
	assert.Equal(t, 2, len(mentions))
}

func TestMentionWithNewLine(t *testing.T) {
	testComment := "@d-truth\n@shanev you hear me?"
	mentions := parseMentions(testComment)
	assert.Equal(t, 2, len(mentions))
}

func TestMentionWithPunctuations(t *testing.T) {
	testComment := "@user-a. (@shanev you hear me @userb?) @user-c, @userd! and @userX:"
	mentions := parseMentions(testComment)
	assert.Equal(t, 6, len(mentions))
}
