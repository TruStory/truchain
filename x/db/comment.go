package db

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/gernest/mention"
	"github.com/kelseyhightower/envconfig"
)

// ChainConfig represents chain env vars
type ChainConfig struct {
	Host               string `required:"true"`
	LetsEncryptEnabled bool   `split_words:"true"`
}

// Comment represents a comment in the DB
type Comment struct {
	Timestamps
	ID         int64  `json:"id"`
	ParentID   int64  `json:"parent_id"`
	ArgumentID int64  `json:"argument_id"`
	Body       string `json:"body"`
	Creator    string `json:"creator"`
}

// CommentsByArgumentID finds comments by argument id
func (c *Client) CommentsByArgumentID(argumentID int64) ([]Comment, error) {
	var chainConfig ChainConfig
	err := envconfig.Process("chain", &chainConfig)
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, 0)
	err = c.Model(&comments).Where("argument_id = ?", argumentID).Select()
	if err != nil {
		return nil, err
	}
	transformedComments := make([]Comment, 0)
	for _, comment := range comments {
		transformedComment := comment
		transformedBody, err := c.replaceAddressesWithProfileURLs(chainConfig, comment.Body)
		if err != nil {
			return transformedComments, err
		}
		transformedComment.Body = transformedBody
		transformedComments = append(transformedComments, transformedComment)
	}

	return transformedComments, nil
}

// replace @cosmosaddr with profile link [@username](https://app.trustory.io/profile/cosmosaddr)
func (c *Client) replaceAddressesWithProfileURLs(config ChainConfig, body string) (string, error) {
	profileURLPrefix := path.Join(config.Host, "profile")
	profileURLsByAddress, err := c.mapAddressesToProfileURLs(config, body, profileURLPrefix)
	if err != nil {
		return "", err
	}
	for address, profileURL := range profileURLsByAddress {
		body = strings.ReplaceAll(body, fmt.Sprintf("@%s", address), profileURL)
	}

	return body, nil
}

func (c *Client) mapAddressesToProfileURLs(config ChainConfig, body string, profileURLPrefix string) (map[string]string, error) {
	profileURLsByAddress := map[string]string{}
	addresses := parseMentions(body)
	for _, address := range addresses {
		twitterProfile, err := c.TwitterProfileByAddress(address)
		if err != nil {
			return profileURLsByAddress, err
		}
		if twitterProfile == nil {
			profileURLsByAddress[address] = address
			continue
		}
		profileURLString := path.Join(profileURLPrefix, twitterProfile.Address)
		profileURL, err := url.Parse(profileURLString)
		if err != nil {
			return profileURLsByAddress, err
		}

		httpPrefix := "http://"
		if config.LetsEncryptEnabled == true {
			httpPrefix = "https://"
		}
		markdownProfileURL := fmt.Sprintf("[@%s](%s%s)", twitterProfile.Username, httpPrefix, profileURL)
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
		if twitterProfile == nil {
			addressByUsername[username] = username
		} else {
			addressByUsername[username] = twitterProfile.Address
		}
	}
	for username, address := range addressByUsername {
		body = strings.ReplaceAll(body, username, address)
	}

	return body, nil
}

// CommentsParticipantsByArgumentID gets the list of users participating on a argument thread.
func (c *Client) CommentsParticipantsByArgumentID(argumentID int64) ([]string, error) {
	comments := make([]Comment, 0)
	addresses := make([]string, 0)
	err := c.Model(&comments).ColumnExpr("DISTINCT creator").Where("argument_id = ?", argumentID).Select()
	if err != nil {
		return nil, err
	}
	for _, c := range comments {
		addresses = append(addresses, c.Creator)
	}
	return addresses, nil
}

// CommentByID returns the comment for specific pk.
func (c *Client) CommentByID(id int64) (*Comment, error) {
	comment := new(Comment)
	err := c.Model(comment).Where("id = ?", id).Select()
	if err != nil {
		return comment, err
	}
	return comment, nil
}
