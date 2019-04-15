package db

import (
	"time"
)

// FlaggedStory represents a flagged story in the DB
type FlaggedStory struct {
	Timestamps
	ID        int64     `json:"id"`
	StoryID   int64     `json:"story_id"`
	Creator   string    `json:"creator"`
	CreatedOn time.Time `json:"created_on"`
}

// FlaggedStoriesByStoryID finds flagged stories by story id
func (c *Client) FlaggedStoriesByStoryID(storyID int64) ([]FlaggedStory, error) {
	flaggedStories := make([]FlaggedStory, 0)
	err := c.Model(&flaggedStories).Where("story_id = ?", storyID).Select()
	if err != nil {
		return nil, err
	}

	return flaggedStories, nil
}

// UpsertFlaggedStory implements `Datastore`.
// Updates an existing `FlaggedStory` or creates a new one.
func (c *Client) UpsertFlaggedStory(flaggedStory *FlaggedStory) error {
	_, err := c.Model(flaggedStory).
		Where("story_id = ?", flaggedStory.StoryID).
		Where("creator = ?", flaggedStory.Creator).
		OnConflict("DO NOTHING").
		SelectOrInsert()

	return err
}
