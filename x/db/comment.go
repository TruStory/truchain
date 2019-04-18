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
	ArgumentID int64     `json:"argument_id"`
	Body       string    `json:"body"`
	Creator    string    `json:"creator"`
	CreatedAt  time.Time `json:"created_at"`
}

// CommentsByArgumentID finds comments by argument id
func (c *Client) CommentsByArgumentID(argumentID int64) ([]Comment, error) {
	comments := make([]Comment, 0)
	err := c.Model(&comments).Where("argument_id = ?", argumentID).Select()
	if err != nil {
		return nil, err
	}
	transformedComments := make([]Comment, 0)
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

// replace @cosmosaddr with profile link [@username](https://app.trustory.io/profile/cosmosaddr)
func (c *Client) replaceAddressesWithProfileURLs(body string) (string, error) {
	var chainConfig ChainConfig
	err := envconfig.Process("chain", &chainConfig)
	profileURLPrefix := path.Join(chainConfig.Host, "profile")
	profileURLsByAddress, err := c.mapAddressesToProfileURLs(body, profileURLPrefix)
	if err != nil {
		return "", err
	}
	for address, profileURL := range profileURLsByAddress {
		body = strings.ReplaceAll(body, fmt.Sprintf("@%s", address), profileURL)
	}

	return body, nil
}

func (c *Client) mapAddressesToProfileURLs(body string, profileURLPrefix string) (map[string]string, error) {
	profileURLsByAddress := map[string]string{}
	addresses := parseMentions(body)
	for _, address := range addresses {
		twitterProfile, err := c.TwitterProfileByAddress(address)
		if err != nil {
			return profileURLsByAddress, err
		}
		profileURL := path.Join(profileURLPrefix, twitterProfile.Address)
		markdownProfileURL := fmt.Sprintf("[@%s](%s)", twitterProfile.Username, profileURL)
		profileURLsByAddress[address] = markdownProfileURL
	}

	return profileURLsByAddress, nil
}

// extract @mentions from text and return as slice
func parseMentions(body string) []string {
	return mention.GetTagsAsUniqueStrings('@', body)
}

// AddComment adds a new comment to the comments table
func (c *Client) AddComment(comment *Comment) error {
	transformedBody, err := c.replaceUsernamesWithAddress(comment.Body)
	if err != nil {
		return err
	}
	comment.Body = transformedBody
	err = c.Add(comment)
	if err != nil {
		return err
	}

	return nil
}

// replace @usernames with @cosmosaddr
func (c *Client) replaceUsernamesWithAddress(body string) (string, error) {
	addressByUsername := map[string]string{}
	usernames := parseMentions(body)
	for _, username := range usernames {
		twitterProfile, err := c.TwitterProfileByUsername(username)
		if err != nil {
			return body, err
		}
		addressByUsername[username] = twitterProfile.Address
	}
	for username, address := range addressByUsername {
		body = strings.ReplaceAll(body, username, address)
	}

	return body, nil
}
