package db

import (
	"fmt"

	"github.com/go-pg/pg"
)

// TwitterProfile is the Twitter profile associated with an account
type TwitterProfile struct {
	Timestamps
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

// UsernamesByPrefix returns the first five usernames for the provided prefix string
func (c *Client) UsernamesByPrefix(prefix string) (usernames []string, err error) {
	var twitterProfiles []TwitterProfile
	sqlFragment := fmt.Sprintf("username LIKE '%s", prefix)
	err = c.Model(&twitterProfiles).Where(sqlFragment + "%'").Limit(5).Select()
	if err == pg.ErrNoRows {
		return usernames, nil
	}
	if err != nil {
		return usernames, err
	}
	for _, twitterProfile := range twitterProfiles {
		usernames = append(usernames, twitterProfile.Username)
	}

	return usernames, nil
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

// TwitterProfileByUsername implements `Datastore`
// Finds a Twitter profile by the given username
func (c *Client) TwitterProfileByUsername(username string) (*TwitterProfile, error) {
	twitterProfile := new(TwitterProfile)
	err := c.Model(twitterProfile).Where("username = ?", username).First()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return twitterProfile, err
	}

	return twitterProfile, nil
}

// UpsertTwitterProfile implements `Datastore`.
// Updates an existing Twitter profile or creates a new one.
func (c *Client) UpsertTwitterProfile(profile *TwitterProfile) error {
	_, err := c.Model(profile).
		OnConflict("(id) DO UPDATE").
		Set("address = EXCLUDED.address, username = EXCLUDED.username, full_name = EXCLUDED.full_name, avatar_uri = EXCLUDED.avatar_uri, email = EXCLUDED.email").
		Insert()

	return err
}
