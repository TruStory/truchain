package db

import (
	"fmt"
	"time"
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
	InsertDeviceToken(token *DeviceToken) error
}

// Queries read from the database
type Queries interface {
	GenericQueries
	TwitterProfileByAddress(addr string) (TwitterProfile, error)
  DeviceTokenByAddress(addr string) (DeviceToken, error)
}


/////////////////// TWITTER PROFILE ///////////////////

// TwitterProfile is the Twitter profile associated with an account
type TwitterProfile struct {
	ID        int64  `json:"id"`
	Address   string `json:"address"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	AvatarURI string `json:"avatar_uri"`
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

// UpsertTwitterProfile implements `Datastore`.
// Updates an existing Twitter profile or creates a new one.
func (c *Client) UpsertTwitterProfile(profile *TwitterProfile) error {
	_, err := c.Model(profile).
		OnConflict("(id) DO UPDATE").
		Set("address = EXCLUDED.address, username = EXCLUDED.username, full_name = EXCLUDED.full_name, avatar_uri = EXCLUDED.avatar_uri").
		Insert()

	return err
}

/////////////////// DEVICE TOKEN  ///////////////////

// DeviceToken is q device token associated with an account
type DeviceToken struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
	Token   string `json:"token"`
}

func (d DeviceToken) String() string {
	return fmt.Sprintf("Device Token<%d %s %s>", d.ID, d.Address, d.Token)
}

// Multiple tokens per address allow users to use multiple devices
func (c *Client) InsertDeviceToken(deviceToken *DeviceToken) error {
	_, err := c.Model(deviceToken).Insert()
	return err
}

// Finds a Device Token by the given address
func (c *Client) DeviceTokenByAddress(addr string) (DeviceToken, error) {
	deviceToken := new(DeviceToken)
	err := c.Model(deviceToken).Where("address = ?", addr).Select()
	if err != nil {
		return *deviceToken, err
	}

	return *deviceToken, nil
}



/////////////////// PUSH NOTIF ///////////////////

// stores a push notif queued for delivery
type PushNotif struct {
	ID        int64
	Token     string
	Payload   string
  Tag       string
	Scheduled time.Time
	Sent      time.Time
}

// PushNotif implements `String`
func (p PushNotif) String() string {
	return fmt.Sprintf("Push Notif<%d %s %s %s %s %s>", p.ID, p.Token, p.Payload, p.Tag, p.Scheduled, p.Sent)
}

// InsertPushNotif implements `Datastore`.
// Inserts a new PushNotif
func (c *Client) InsertPushNotif(notif *PushNotif) error {
	_, err := c.Model(notif).Insert()
	return err
}

// UpdatePushNotif implements `Datastore`.
// Updates an existing notif to mark as sent.
func (c *Client) UpdatePushNotif(notif *PushNotif) error {
	_, err := c.Model(notif).Update()
	return err
}
