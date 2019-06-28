package community

import (
	"fmt"
	"time"
)

// Defines module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName
)

// Community represents the state of a community on TruStory
type Community struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
}

// Communities is a slice of communites
type Communities []Community

// NewCommunity creates a new Community
func NewCommunity(id, name, description string, createdTime time.Time) Community {
	return Community{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedTime: createdTime,
	}
}

func (c Community) String() string {
	return fmt.Sprintf(`Community:
   ID: 			    %s
   Name: 			%s
   Description:  	%s
   CreatedTime: 	%s`,
		c.ID, c.Name, c.Description, c.CreatedTime.String())
}
