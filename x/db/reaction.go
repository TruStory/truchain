package db

import (
	"fmt"
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

// ReactionsByReactionable returns all the reactions left by all the users on a particular reactionable
func (c *Client) ReactionsByReactionable(reactionable Reactionable) ([]Reaction, error) {
	reactions := make([]Reaction, 0)

	err := c.Model(&reactions).
		Column("reactions.*").
		Where("reactions.reactionable_type = ?", reactionable.Type).
		Where("reactions.reactionable_id = ?", reactionable.ID).
		Where("reactions.deleted_at IS NOT NULL").
		Order("timestamp DESC").
		Select()

	if err != nil {
		return nil, err
	}
	return reactions, nil
}

// ReactionsByAddress returns all the reactions left by a user on any reactionable
func (c *Client) ReactionsByAddress(addr string) ([]Reaction, error) {
	reactions := make([]Reaction, 0)

	err := c.Model(&reactions).
		Column("reactions.*").
		Where("reactions.creator = ?", addr).
		Where("reactions.deleted_at IS NOT NULL").
		Order("timestamp DESC").
		Select()

	if err != nil {
		return nil, err
	}
	return reactions, nil
}

// ReactionByAddressAndReactionable returns the reaction left by a particular user on a particular reactionable
func (c *Client) ReactionByAddressAndReactionable(addr string, reactionable Reactionable) (*Reaction, error) {
	rxn := new(Reaction)

	err := c.Model(rxn).
		Where("creator = ?", addr).
		Where("reactionable_type = ?", reactionable.Type).
		Where("reactionable_id = ?", reactionable.ID).
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
	rxn, err := c.ReactionByAddressAndReactionable(addr, reactionable)
	fmt.Printf("\n\nRXN -- %v\n\n", rxn)
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
