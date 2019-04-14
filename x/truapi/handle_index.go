package truapi

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
)

const (
	defaultTitle       = "TruStory"
	defaultDescription = "Trustory is a social network for experts to identify what is true and what isn't."
	defaultImage       = "pathtoimageurl"
)

// Tags defines the struct containing all the request Meta Tags for a page
type Tags struct {
	Title       string
	Description string
	Image       string
	URL         string
}

// CompileIndexFile replaces the placeholders for the social sharing
func CompileIndexFile(ta *TruAPI, index []byte, route string) string {

	// /story/detail/xxx
	r, err := regexp.Compile("/story/detail/([0-9]+)")
	if err != nil {
		panic(err)
	}
	matches := r.FindStringSubmatch(route)
	if len(matches) == 2 {
		// replace placeholder with story details, where story id is in matches[1]
		storyID, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			panic(err)
		}

		return compile(index, makeStoryMetaTags(ta, route, storyID))
	}

	return compile(index, makeDefaultMetaTags(ta, route))
}

// compiles the index file with the variables
func compile(index []byte, tags Tags) string {
	compiled := bytes.Replace(index, []byte("$PLACEHOLDER__TITLE"), []byte(tags.Title), -1)
	compiled = bytes.Replace(compiled, []byte("$PLACEHOLDER__DESCRIPTION"), []byte(tags.Description), -1)
	compiled = bytes.Replace(compiled, []byte("$PLACEHOLDER__IMAGE"), []byte(tags.Image), -1)
	compiled = bytes.Replace(compiled, []byte("$PLACEHOLDER__URL"), []byte(tags.URL), -1)

	return string(compiled)
}

// makes the default meta tags
func makeDefaultMetaTags(ta *TruAPI, route string) Tags {
	return Tags{
		Title:       defaultTitle,
		Description: defaultDescription,
		Image:       defaultImage,
		URL:         os.Getenv("APP_URL") + route,
	}
}

// meta tags for a story
func makeStoryMetaTags(ta *TruAPI, route string, storyID int64) Tags {
	ctx := context.Background()

	storyObj := ta.storyResolver(ctx, story.QueryStoryByIDParams{ID: storyID})
	categoryObj := ta.categoryResolver(ctx, category.QueryCategoryByIDParams{ID: storyObj.CategoryID})
	creatorObj, err := ta.DBClient.TwitterProfileByAddress(storyObj.Creator.String())
	if err != nil {
		panic(err)
	}
	return Tags{
		Title:       fmt.Sprintf("%s made a claim in %s on TruStory", creatorObj.FullName, categoryObj.Title),
		Description: storyObj.Body,
		Image:       defaultImage,
		URL:         os.Getenv("APP_URL") + route,
	}
}
