package db

import (
	"time"

	"github.com/go-pg/pg"
)

// ReactionType represents the type of reaction left on the argument
type ReactionType int64

// ReactionableType represents the type of model that is reacted on
type ReactionableType string

// Reactionable represents any model that can be reacted upon (polymorphism design pattern)
type Reactionable struct {
	Type ReactionableType
	ID   int64
}

const (
	// GotAnIdea represents a reaction where the reacting user has got a new idea based on the ReactionableType
	GotAnIdea ReactionType = iota + 1

	// ChangedMyMind represents a reaction where the reacting user has changed their minds based on the ReactionableType
	ChangedMyMind
)

const (
	// Argument represents a type of reactionable
	Argument ReactionableType = "arguments"
)

// Reaction represents a reaction left by a user
type Reaction struct {
	Timestamps

	ID               int64            `json:"id"`
	ReactionableType ReactionableType `json:"reactionable_type"`
	ReactionableID   int64            `json:"reactionable_id"`
	ReactionType     ReactionType     `json:"reaction_type"`
	Creator          string           `json:"creator"`
}

// ReactionsCount represents the structure to contain the count of each reaction left on a reactionable
type ReactionsCount struct {
	Type  ReactionType `json:"type"`
	Count int64        `json:"count"`
}

// ReactionsByReactionable returns all the reactions left by all the users on a particular reactionable
func (c *Client) ReactionsByReactionable(reactionable Reactionable) ([]Reaction, error) {
	reactions := make([]Reaction, 0)

	err := c.Model(&reactions).
		Where("reactionable_type = ?", reactionable.Type).
		Where("reactionable_id = ?", reactionable.ID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Select()

	if err != nil {
		return nil, err
	}
	return reactions, nil
}

// ReactionsCountByReactionable returns the cound of each reaction type for a particular reactionable
func (c *Client) ReactionsCountByReactionable(reactionable Reactionable) ([]ReactionsCount, error) {

	var result []ReactionsCount

	err := c.Model((*Reaction)(nil)).
		ColumnExpr("reaction_type as type").
		ColumnExpr("count(*) AS count").
		Group("reaction_type").
		Order("reaction_type").
		Select(&result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReactionsByAddress returns all the reactions left by a user on any reactionable
func (c *Client) ReactionsByAddress(addr string) ([]Reaction, error) {
	reactions := make([]Reaction, 0)

	err := c.Model(&reactions).
		Where("creator = ?", addr).
		Where("deleted_at IS NOT NULL").
		Order("created_at DESC").
		Select()

	if err != nil {
		return nil, err
	}
	return reactions, nil
}

// ReactionByAddressAndReactionable returns the reaction left by a particular user on a particular reactionable
func (c *Client) ReactionByAddressAndReactionable(addr string, reaction ReactionType, reactionable Reactionable) (*Reaction, error) {
	rxn := new(Reaction)

	err := c.Model(rxn).
		Where("creator = ?", addr).
		Where("reactionable_type = ?", reactionable.Type).
		Where("reactionable_id = ?", reactionable.ID).
		Where("reaction_type = ?", reaction).
		Where("deleted_at IS NULL").
		First()

	if err == pg.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return rxn, nil
}

// ReactOnReactionable leaves a reaction by a user on a reactionable
func (c *Client) ReactOnReactionable(addr string, reaction ReactionType, reactionable Reactionable) error {
	// checking if already exists
	rxn, err := c.ReactionByAddressAndReactionable(addr, reaction, reactionable)
	if err != nil {
		return err
	}

	// if found, then, do nothing
	if rxn != nil {
		return nil
	}

	rxn = &Reaction{
		ReactionableType: reactionable.Type,
		ReactionableID:   reactionable.ID,
		ReactionType:     reaction,
		Creator:          addr,
	}
	err = c.Insert(rxn)
	if err != nil {
		return err
	}

	return nil
}

// UnreactByAddressAndID removes a reaction by a user on a reactionable
// We are avoiding using just ID to protect our database from abuse.
// IDs are auto-incrementing numbers, thus, easier to guess and abuse.
func (c *Client) UnreactByAddressAndID(addr string, id int64) error {
	rxn := new(Reaction)

	_, err := c.Model(rxn).
		Where("id = ?", id).
		Where("creator = ?", addr).
		Where("deleted_at IS NULL").
		Set("deleted_at = ?", time.Now()).
		Update()

	if err != nil {
		return err
	}

	return nil
}
