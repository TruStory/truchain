package db

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/gernest/mention"
	"github.com/kelseyhightower/envconfig"
)

// ChainConfig represents chain env vars
type ChainConfig struct {
	Host string
}

// Comment represents a comment in the DB
type Comment struct {
	ID         int64     `json:"id"`
	ParentID   int64     `json:"parent_id"`
	ArgumentID int64     `json:"argument_id" sql:"notnull"`
	Body       string    `json:"body" sql:"notnull"`
	Creator    string    `json:"creator" sql:"notnull"`
	CreatedAt  time.Time `json:"created_at" sql:"notnull"`
}

// CommentsByArgumentID finds comments by argument id
func (c *Client) CommentsByArgumentID(argumentID int64) ([]Comment, error) {
	comments := make([]Comment, 0)
	err := c.Model(&comments).Where("argument_id = ?", argumentID).Select()
	if err != nil {
		return nil, err
	}

	// replace @cosmosaddr with profile link [@username](https://app.trustory.io/profile/username)
	transformedComments := make([]Comment, len(comments))
	for _, comment := range comments {
		transformedComment := comment
		transformedBody, err := c.replaceAddressesWithProfileURLs(comment.Body)
		if err != nil {
			return transformedComments, err
		}
		transformedComment.Body = transformedBody
		transformedComments = append(transformedComments, transformedComment)
	}

	return transformedComments, nil
}

// AddComment adds a new comment to the comments table
func (c *Client) AddComment(comment *Comment) error {
	// TODO: replace @mentions with @cosmosaddr
	err := c.Add(comment)
	if err != nil {
		return err
	}
	return nil
}

// replace @mentions with @cosmosaddr
func (c *Client) replaceUsernameWithAddress(body string) string {
	return ""
}

// replace @cosmosaddr with profile link [@username](https://app.trustory.io/profile/username)
func (c *Client) replaceAddressesWithProfileURLs(body string) (string, error) {
	var chainConfig ChainConfig
	err := envconfig.Process("chain", &chainConfig)
	profileURLPrefix := path.Join(chainConfig.Host, "profile")
	profileURLsByAddress, err := c.mapAddressesToProfileURLs(body, profileURLPrefix)
	if err != nil {
		return "", err
	}
	for address, profileURL := range profileURLsByAddress {
		markdownLink := fmt.Sprintf("[%s](%s)")
		body = strings.ReplaceAll(body, fmt.Sprintf("@%s", address), profileURL)
	}

	return body, nil
}

func (c *Client) mapAddressesToProfileURLs(body string, profileURLPrefix string) (map[string]string, error) {
	profileURLsByAddress := map[string]string{}
	usernames := parseMentions(body)
	for _, username := range usernames {
		twitterProfile, err := c.TwitterProfileByUsername(username)
		if err != nil {
			return profileURLsByAddress, err
		}
		profileURLsByAddress[twitterProfile.Address] = path.Join(profileURLPrefix, twitterProfile.Address)
	}

	return profileURLsByAddress, nil
}

func parseMentions(body string) []string {
	return mention.GetTagsAsUniqueStrings('@', body)
}
