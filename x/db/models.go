package db

import (
	"fmt"

	"github.com/go-pg/pg"
)

// Datastore defines all operations on the DB
// This interface can be mocked out for tests, etc.
type Datastore interface {
	Mutations
	Queries
}

// Mutations write to the database
type Mutations interface {
	GenericMutations
	UpsertTwitterProfile(profile *TwitterProfile) error
}

// Queries read from the database
type Queries interface {
	GenericQueries
	TwitterProfileByAddress(addr string) (TwitterProfile, error)
	KeyPairByTwitterProfileID(id int64) (KeyPair, error)
}

// TwitterProfile is the Twitter profile associated with an account
type TwitterProfile struct {
	ID        int64  `json:"id"`
	Address   string `json:"address"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	AvatarURI string `json:"avatar_uri"`
}

// KeyPair is the private key associated with an account
type KeyPair struct {
	ID               int64  `json:"id"`
	TwitterProfileID int64  `json:"twitter_profile_id"`
	PrivateKey       string `json:"private_key"`
	PublicKey        string `json:"public_key"`
}

func (t TwitterProfile) String() string {
	return fmt.Sprintf(
		"Twitter Profile<%d %s %s %s %s>",
		t.ID, t.Address, t.Username, t.FullName, t.AvatarURI)
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

// KeyPairByTwitterProfileID returns the key-pair for the account
func (c *Client) KeyPairByTwitterProfileID(id int64) (KeyPair, error) {
	keyPair := new(KeyPair)
	err := c.Model(keyPair).Where("twitter_profile_id = ?", id).First()

	if err == pg.ErrNoRows {
		return *keyPair, nil
	}

	if err != nil {
		return *keyPair, err
	}

	return *keyPair, nil
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
