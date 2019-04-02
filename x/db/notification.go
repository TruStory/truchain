package db

import (
	"time"
)

// NotificationType represents a type of notification defiend by the system.
type NotificationType int

// Types of notifications.
const (
	NotificationStoryAction NotificationType = iota
	NotificationStoryUpdate
)

// NotificationEvent represents a notification sent to an user.
type NotificationEvent struct {
	ID               int64            `json:"id"`
	Address          string           `json:"address"`
	TwitterProfileID int64            `json:"profile_id"`
	TwitterProfile   *TwitterProfile  `json:"profile"`
	Message          string           `json:"message"`
	Timestamp        time.Time        `json:"timestamp"`
	SenderProfileID  int64            `json:"sender_profile_id" `
	SenderProfile    *TwitterProfile  `json:"sender_profile"  pg:"fk:twitter_profile_id"`
	Type             NotificationType `json:"type"`
	Read             bool             `json:"read"`
}

// NotificationEventsByAddress retrieves all notifications sent to an user.
// TODO: add pagination
func (c *Client) NotificationEventsByAddress(addr string) ([]NotificationEvent, error) {
	evts := make([]NotificationEvent, 0)

	err := c.Model(&evts).
		Column("notification_event.*", "TwitterProfile", "SenderProfile").
		Where("notification_event.address = ?", addr).Select()
	if err != nil {
		return nil, err
	}
	return evts, nil
}
