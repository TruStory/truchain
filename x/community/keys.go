package community

// Keys for community store
// Items are stored with the following key: values
//
// - 0x00<communityID_Bytes>: Community{} bytes
var (
	CommunityKeyPrefix = []byte{0x00}
)

// key for getting a specific community from the store
func key(id string) []byte {
	return append(CommunityKeyPrefix, []byte(id)...)
}
