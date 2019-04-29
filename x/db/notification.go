package db

import (
	"time"
)

// NotificationType represents a type of notification defiend by the system.
type NotificationType int

// Types of notifications.
const (
	NotificationStoryAction NotificationType = iota
	NotificationCommentAction
)

// NotificationEvent represents a notification sent to an user.
type NotificationEvent struct {
	Timestamps
	ID               int64            `json:"id"`
	TypeID           int64            `json:"type_id"`
	Address          string           `json:"address"`
	TwitterProfileID int64            `json:"profile_id"`
	TwitterProfile   *TwitterProfile  `json:"profile"`
	Message          string           `json:"message"`
	Timestamp        time.Time        `json:"timestamp"`
	SenderProfileID  int64            `json:"sender_profile_id" `
	SenderProfile    *TwitterProfile  `json:"sender_profile"`
	Type             NotificationType `json:"type" sql:",notnull"`
	Read             bool             `json:"read"`
}

// NotificationEventsByAddress retrieves all notifications sent to an user.
// TODO (issue #435): add pagination
func (c *Client) NotificationEventsByAddress(addr string) ([]NotificationEvent, error) {
	evts := make([]NotificationEvent, 0)

	err := c.Model(&evts).
		Column("notification_event.*", "TwitterProfile", "SenderProfile").
		Where("notification_event.address = ?", addr).Order("timestamp DESC").Select()
	if err != nil {
		return nil, err
	}
	return evts, nil
}
