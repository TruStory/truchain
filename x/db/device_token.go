package db

import "github.com/go-pg/pg"

// DeviceToken is the association between a cosmos address and a device token used for
// push notifications.
type DeviceToken struct {
	Timestamps
	ID int64 `json:"id"`

	// Address is the cosmos address
	Address string `json:"address"  sql:"unique:device_address_token,notnull"`
	// Token represents the DeviceToken (iOS), RegistrationId (android)
	Token string `json:"token"  sql:"unique:device_address_token,notnull"`
	// Platform indicates to which platform the token belongs to : android, ios
	Platform string `json:"platform"  sql:"unique:device_address_token,notnull"`
}

// UpsertDeviceToken implements `Datastore`.
// Updates an existing DeviceToken or creates a new one.
func (c *Client) UpsertDeviceToken(token *DeviceToken) error {
	_, err := c.TwitterProfileByAddress(token.Address)
	if err == pg.ErrNoRows {
		return ErrInvalidAddress
	}
	if err != nil {
		return err
	}
	_, err = c.Model(token).
		Where("address = ? ", token.Address).
		Where("token = ?", token.Token).
		Where("platform = ?", token.Platform).
		OnConflict("DO NOTHING").
		SelectOrInsert()
	return err
}

// RemoveDeviceToken performs a delete for a specific token.
func (c *Client) RemoveDeviceToken(address, token, platform string) error {
	t := &DeviceToken{}
	_, err := c.Model(t).Where("address = ? ", address).
		Where("token = ?", token).
		Where("platform = ?", platform).Delete()
	return err
}

// DeviceTokensByAddress implements `Datastore`
// Finds a Device Tokens by the given address
func (c *Client) DeviceTokensByAddress(addr string) ([]DeviceToken, error) {
	deviceTokens := make([]DeviceToken, 0)
	err := c.Model(&deviceTokens).Where("address = ?", addr).Select()
	if err != nil {
		return nil, err
	}
	return deviceTokens, nil
}
