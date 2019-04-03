package db

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
