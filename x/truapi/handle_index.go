package truapi

import (
	"bytes"
	"regexp"
)

// CompileIndexFile replaces the placeholders for the social sharing
func CompileIndexFile(index []byte, route string) string {

	// /story/detail/xxx
	r, err := regexp.Compile("/story/detail/([0-9]+)")
	if err != nil {
		panic(err)
	}
	matches := r.FindStringSubmatch(route)
	if len(matches) == 2 {
		// replace placeholder with story details, where story id is in matches[1]
		compiled := bytes.Replace(index, []byte("$PLACEHOLDER__TITLE"), []byte("TruStory - Story page"), -1)
		return string(compiled)
	}

	// default case
	compiled := bytes.Replace(index, []byte("$PLACEHOLDER__TITLE"), []byte("TruStory"), -1)
	return string(compiled)
}
