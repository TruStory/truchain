package db

import (
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
}

// Queries read from the database
type Queries interface {
	GenericQueries
	TwitterProfileByID(id int64) (TwitterProfile, error)
	TwitterProfileByAddress(addr string) (TwitterProfile, error)
	KeyPairByTwitterProfileID(id int64) (KeyPair, error)
	DeviceTokensByAddress(addr string) ([]DeviceToken, error)
	NotificationEventsByAddress(addr string) ([]NotificationEvent, error)
}

// TimestampedModel carries the default timestamp fields for any derived model
type TimestampedModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// BeforeInsert is the hook that fills in the created_at and updated_at fields
func (m *TimestampedModel) BeforeInsert(db orm.DB) error {
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
func (m *TimestampedModel) BeforeUpdate(db orm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
