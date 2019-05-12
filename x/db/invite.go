package db

// Invite represents an invite from a friend in the DB
type Invite struct {
	ID                    int64  `json:"id"`
	Creator               string `json:"creator"`
	FriendTwitterUsername string `json:"friend_twitter_username"`
	FriendEmail           string `json:"friend_email"`
	Paid                  bool   `json:"paid"`
	Timestamps
}

// Invites returns all invites in theDB
func (c *Client) Invites() ([]Invite, error) {
	invites := make([]Invite, 0)
	err := c.Model(&invites).Select()
	if err != nil {
		return nil, err
	}

	return invites, nil
}

// AddInvite inserts an invitation
func (c *Client) AddInvite(invite *Invite) error {
	_, err := c.Model(invite).
		OnConflict("DO NOTHING").
		Insert()

	return err
}
