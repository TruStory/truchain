package db

import (
	"time"
)

// NotificationType represents a type of notification defiend by the system.
type NotificationType int

// NotificationsCountResponse is the interface to respond graphQL query
type NotificationsCountResponse struct {
	Count int64 `json:"count"`
}

// Types of notifications.
const (
	NotificationStoryAction NotificationType = iota
	NotificationCommentAction
)

// NotificationMeta  contains extra payload information.
type NotificationMeta struct {
	ArgumentID *int64 `json:"argumentId,omitempty" graphql:"argumentId"`
	StoryID    *int64 `json:"storyId,omitempty" graphql:"storyId"`
	CommentID  *int64 `json:"commentId,omitempty" graphql:"commentId"`
}

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
	Meta             NotificationMeta `json:"meta"`
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

// UnreadNotificationEventsCountByAddress retrieves the number of unread notifications sent to an user.
func (c *Client) UnreadNotificationEventsCountByAddress(addr string) (*NotificationsCountResponse, error) {
	notificationEvent := new(NotificationEvent)

	count, err := c.Model(notificationEvent).
		Where("notification_event.address = ?", addr).
		Where("read = ?", false).Count()
	if err != nil {
		return &NotificationsCountResponse{
			Count: 0,
		}, err
	}

	return &NotificationsCountResponse{
		Count: int64(count),
	}, nil
}

// MarkAllNotificationEventsAsReadByAddress retrieves the number of unread notifications sent to an user.
func (c *Client) MarkAllNotificationEventsAsReadByAddress(addr string) error {
	notificationEvent := new(NotificationEvent)

	_, err := c.Model(notificationEvent).
		Where("notification_event.address = ?", addr).
		Where("read = ?", false).
		Set("read = ?", true).
		Update()
	if err != nil {
		return err
	}

	return nil
}
