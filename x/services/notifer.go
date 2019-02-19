package services

import ( 
  "errors"
  "fmt"
  "github.com/TruStory/truchain/x/db"
)

// used to lookup the recipients device token
type DraftNotif struct {
  RecipientAddress  string
  Payload           string
  Tag               string
}


func QueueNotif(n *DraftNotif) {

  // query db for deviceToken
  // save PushNotif
}
