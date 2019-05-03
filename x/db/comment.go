package db

import (
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
