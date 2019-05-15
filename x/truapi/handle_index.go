package truapi

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"os"
	"regexp"
	"strconv"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stripmd "github.com/writeas/go-strip-markdown"
)

const (
	defaultDescription = "TruStory is a social network to debate claims with skin in the game"
	defaultImage       = "https://s3-us-west-1.amazonaws.com/trustory/assets/Image+from+iOS.jpg"
)

var (
	storyRegex    = regexp.MustCompile("/story/([0-9]+)$")
	argumentRegex = regexp.MustCompile("/story/([0-9]+)/argument/([0-9]+)$")
	commentRegex  = regexp.MustCompile("/story/([0-9]+)/argument/([0-9]+)/comment/([0-9]+)$")
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

	// /story/xxx
	matches := storyRegex.FindStringSubmatch(route)
	if len(matches) == 2 {
		// replace placeholder with story details, where story id is in matches[1]
		storyID, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			// if error, return the default tags
			return compile(index, makeDefaultMetaTags(ta, route))
		}

		metaTags, err := makeStoryMetaTags(ta, route, storyID)
		if err != nil {
			return compile(index, makeDefaultMetaTags(ta, route))
		}
		return compile(index, *metaTags)
	}

	// /story/xxx/argument/xxx
	matches = argumentRegex.FindStringSubmatch(route)
	if len(matches) == 3 {
		// replace placeholder with story details, where story id is in matches[1]
		storyID, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			// if error, return the default tags
			return compile(index, makeDefaultMetaTags(ta, route))
		}
		argumentID, err := strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			// if error, return the default tags
			return compile(index, makeDefaultMetaTags(ta, route))
		}

		metaTags, err := makeArgumentMetaTags(ta, route, storyID, argumentID)
		if err != nil {
			return compile(index, makeDefaultMetaTags(ta, route))
		}
		return compile(index, *metaTags)
	}

	// /story/xxx/argument/xxx/comment/xxx
	matches = commentRegex.FindStringSubmatch(route)
	if len(matches) == 4 {
		// replace placeholder with story details, where story id is in matches[1]
		storyID, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			// if error, return the default tags
			return compile(index, makeDefaultMetaTags(ta, route))
		}
		argumentID, err := strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			// if error, return the default tags
			return compile(index, makeDefaultMetaTags(ta, route))
		}
		commentID, err := strconv.ParseInt(matches[3], 10, 64)
		if err != nil {
			// if error, return the default tags
			return compile(index, makeDefaultMetaTags(ta, route))
		}

		metaTags, err := makeCommentMetaTags(ta, route, storyID, argumentID, commentID)
		if err != nil {
			return compile(index, makeDefaultMetaTags(ta, route))
		}
		return compile(index, *metaTags)
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
		Title:       os.Getenv("APP_NAME"),
		Description: defaultDescription,
		Image:       defaultImage,
		URL:         os.Getenv("APP_URL") + route,
	}
}

// meta tags for a story
func makeStoryMetaTags(ta *TruAPI, route string, storyID int64) (*Tags, error) {
	ctx := context.Background()

	storyObj := ta.storyResolver(ctx, story.QueryStoryByIDParams{ID: storyID})
	backings := ta.backingsResolver(ctx, app.QueryByIDParams{ID: storyObj.ID})
	challenges := ta.challengesResolver(ctx, app.QueryByIDParams{ID: storyObj.ID})
	backingTotalAmount := ta.backingPoolResolver(ctx, storyObj)
	challengeTotalAmount := ta.challengePoolResolver(ctx, storyObj)

	storyState := "Active"
	if storyObj.Status == story.Expired {
		storyState = "Completed"
	}
	totalParticipants := len(backings) + len(challenges)
	totalParticipantsPlural := "s"
	if totalParticipants == 1 {
		totalParticipantsPlural = ""
	}
	totalStake := backingTotalAmount.Plus(challengeTotalAmount).Amount.Div(sdk.NewInt(app.Shanev))

	return &Tags{
		Title:       html.EscapeString(storyObj.Body),
		Description: fmt.Sprintf("%s: %d participant%s and %s TruStake", storyState, totalParticipants, totalParticipantsPlural, totalStake),
		Image:       defaultImage,
		URL:         os.Getenv("APP_URL") + route,
	}, nil
}

func makeArgumentMetaTags(ta *TruAPI, route string, storyID int64, argumentID int64) (*Tags, error) {
	ctx := context.Background()
	storyObj := ta.storyResolver(ctx, story.QueryStoryByIDParams{ID: storyID})
	categoryObj := ta.categoryResolver(ctx, category.QueryCategoryByIDParams{ID: storyObj.CategoryID})
	argumentObj := ta.argumentResolver(ctx, app.QueryArgumentByID{ID: argumentID})
	creatorObj, err := ta.DBClient.TwitterProfileByAddress(argumentObj.Creator.String())
	if err != nil {
		// if error, return default
		return nil, err
	}
	return &Tags{
		Title:       fmt.Sprintf("%s made an argument in %s", creatorObj.FullName, categoryObj.Title),
		Description: html.EscapeString(stripmd.Strip(argumentObj.Body)),
		Image:       defaultImage,
		URL:         os.Getenv("APP_URL") + route,
	}, nil
}

func makeCommentMetaTags(ta *TruAPI, route string, storyID int64, argumentID int64, commentID int64) (*Tags, error) {
	ctx := context.Background()
	storyObj := ta.storyResolver(ctx, story.QueryStoryByIDParams{ID: storyID})
	categoryObj := ta.categoryResolver(ctx, category.QueryCategoryByIDParams{ID: storyObj.CategoryID})
	argumentObj := ta.argumentResolver(ctx, app.QueryArgumentByID{ID: argumentID})
	comments := ta.commentsResolver(ctx, argumentObj)
	commentObj := db.Comment{}
	for _, comment := range comments {
		if comment.ID == commentID {
			commentObj = comment
		}
	}
	creatorObj, err := ta.DBClient.TwitterProfileByAddress(argumentObj.Creator.String())
	if err != nil {
		// if error, return default
		return nil, err
	}
	return &Tags{
		Title:       fmt.Sprintf("%s posted a comment in %s", creatorObj.FullName, categoryObj.Title),
		Description: html.EscapeString(stripmd.Strip(commentObj.Body)),
		Image:       defaultImage,
		URL:         os.Getenv("APP_URL") + route,
	}, nil
}
