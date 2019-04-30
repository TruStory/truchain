package db

import (
	"context"
	"time"

	"github.com/go-pg/pg/orm"
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
	UpsertDeviceToken(token *DeviceToken) error
	RemoveDeviceToken(address, token, platform string) error
	UpsertFlaggedStory(flaggedStory *FlaggedStory) error
	MarkAllNotificationEventsAsReadByAddress(addr string) error
	AddComment(comment *Comment) error
}

// Queries read from the database
type Queries interface {
	GenericQueries
	TwitterProfileByID(id int64) (TwitterProfile, error)
	TwitterProfileByAddress(addr string) (*TwitterProfile, error)
	TwitterProfileByUsername(username string) (*TwitterProfile, error)
	UsernamesByPrefix(prefix string) ([]string, error)
	KeyPairByTwitterProfileID(id int64) (KeyPair, error)
	DeviceTokensByAddress(addr string) ([]DeviceToken, error)
	NotificationEventsByAddress(addr string) ([]NotificationEvent, error)
	UnreadNotificationEventsCountByAddress(addr string) (*NotificationsCountResponse, error)
	FlaggedStoriesByStoryID(storyID int64) ([]FlaggedStory, error)
	CommentsByArgumentID(argumentID int64) ([]Comment, error)
}

// Timestamps carries the default timestamp fields for any derived model
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// BeforeInsert is the hook that fills in the created_at and updated_at fields
func (m *Timestamps) BeforeInsert(ctx context.Context, db orm.DB) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate is the hook that updates the updated_at field
func (m *Timestamps) BeforeUpdate(ctx context.Context, db orm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
