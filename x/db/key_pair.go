package db

// KeyPair is the private key associated with an account
type KeyPair struct {
	ID               int64  `json:"id"`
	TwitterProfileID int64  `json:"twitter_profile_id"`
	PrivateKey       string `json:"private_key"`
	PublicKey        string `json:"public_key"`
}
