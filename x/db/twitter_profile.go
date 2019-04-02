package db

import (
	"fmt"

	"github.com/go-pg/pg"
)

// TwitterProfile is the Twitter profile associated with an account
type TwitterProfile struct {
	ID        int64  `json:"id"`
	Address   string `json:"address"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	AvatarURI string `json:"avatar_uri"`
}

func (t TwitterProfile) String() string {
	return fmt.Sprintf(
		"Twitter Profile<%d %s %s %s %s>",
		t.ID, t.Address, t.Username, t.FullName, t.AvatarURI)
}

// TwitterProfileByID implements `Datastore`
// Finds a Twitter profile by the given twitter profile id
func (c *Client) TwitterProfileByID(id int64) (TwitterProfile, error) {
	twitterProfile := new(TwitterProfile)
	err := c.Model(twitterProfile).Where("id = ?", id).Select()

	if err == pg.ErrNoRows {
		return *twitterProfile, nil
	}

	if err != nil {
		return *twitterProfile, err
	}

	return *twitterProfile, nil
}

// TwitterProfileByAddress implements `Datastore`
// Finds a Twitter profile by the given address
func (c *Client) TwitterProfileByAddress(addr string) (TwitterProfile, error) {
	twitterProfile := new(TwitterProfile)
	err := c.Model(twitterProfile).Where("address = ?", addr).Select()
	if err != nil {
		return *twitterProfile, err
	}

	return *twitterProfile, nil
}

// UpsertTwitterProfile implements `Datastore`.
// Updates an existing Twitter profile or creates a new one.
func (c *Client) UpsertTwitterProfile(profile *TwitterProfile) error {
	_, err := c.Model(profile).
		OnConflict("(id) DO UPDATE").
		Set("address = EXCLUDED.address, username = EXCLUDED.username, full_name = EXCLUDED.full_name, avatar_uri = EXCLUDED.avatar_uri").
		Insert()

	return err
}
