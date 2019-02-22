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

//func NotifyStoryCreated..


//func NotifyStoryBacked..

func NotifyStoryChallenged(storyID int64, argument string) (err sdk.Error) {
  //lookup Story for creators address
  //create Draft Notif
  draft := DraftNotif{ "recipient", "payloaddd", "TAgtaag" }
  //pass to QueueNotif
  QueueNotif(draft)
}


func queueNotif(n *DraftNotif) {
  // query db for deviceToken using DraftNotif.recipientAddress
  // Insert PushNotif into push_notifs  db table queue
}
