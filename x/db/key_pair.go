package db

import "github.com/go-pg/pg"

// KeyPair is the private key associated with an account
type KeyPair struct {
	ID               int64  `json:"id"`
	TwitterProfileID int64  `json:"twitter_profile_id"`
	PrivateKey       string `json:"private_key"`
	PublicKey        string `json:"public_key"`
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
