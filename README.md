# TruStory Cosmos app

### Types

All state is stored in these types.

```go
type Bond struct {
    ID              int64               // id of the bond
    Amount          sdk.Coin            // amount of the bond
    CreatedBlock    int64               // block at which bond was created
    Creator         sdk.AccAddress      // address of the creator
    Period          time.Duration       // time period of the bond (days)    
    StoryID         int64               // id of the associated story
}

type Comment struct {
    ID              int64               // id of comment
    Body            string              // body of comment
    Creator         sdk.AccAddress      // address of creator
    StoryID         int64               // associated story
}

type Evidence struct {
    ID              int64               // id of evidence
    Creator         sdk.Address         // address of creator
    StoryID         int64               // id of associated story
    URI             string              // uri of evidence
}

type Story struct {
    ID              int64               // id of story
    Body            string              // body of story
    BondIDs         []int64             // bonds ids for story
    Category        string              // category slug, "btc", "eth", "cp" (consensus protocols)
    CommentIDs      []int64             // ids of comments
    CreatedBlock    int64               // block at which story was created
    Creator         sdk.AccAddress      // address of creator
    EvidenceIDs     []int64             // ids of evidence
    Expiration      time.Time           // expiration time of story
    Rank            float64             // story rank
    State           string              // "created", "validated", "rejected", "unverifiable", "challenged", "revoked"
    SubmitBlock     int64               // block at which voting begins
    Thread          []int64             // associated story ids
    Type            string              // "identity", "recovery", "default"
    UpdatedBlock    int64               // block at which story was updated
    Users           []sdk.AccAddress    // users mention in story body
    VoteIDs         []int64             // vote ids for story
}

type Vote struct {
    ID              int64               // id of vote
    CreatedBlock    int64               // block at which vote was cast
    Creator         sdk.AccAddress      // address of creator
    StoryID         int64               // id of associated story
    Vote            bool                // yes or no
}
```

### Messages

These are the messages needed to modify the above state.

- `PlaceBondMsg`: to place a bond on a story
- `AddCommentMsg`: to add a comment to a story
- `SubmitEvidenceMsg`: to submit evidence for a story
- `SubmitStoryMsg`: to submit a story
- `VoteMsg`: to vote on a story

```go
type PlaceBondMsg struct {
    StoryID         int64           // id of the story
    Stake           sdk.Coin        // amount of bond
    Creator         sdk.AccAddress  // person who is placing the bond
    Period          time.Duration   // time period of bond (days)
}

type AddCommentMsg struct {
    StoryID         int64           // id of the story
    Body            string          // body of comment
    Creator         sdk.AccAddress  // creator of comment
}

type SubmitEvidenceMsg struct {
    StoryID         int64           // id of the story
    Creator         sdk.AccAddress  // creator of evidence submission
    URI             string          // uri of evidence
}

type SubmitStoryMsg struct {
    Body            string              // body of story
    Category        string              // category of story
    Creator         sdk.AccAddress      // creator of story
    StoryType       string              // type of story
    Users           []sdk.AccAddress    // addresses of mentioned users
}

type VoteMsg struct {
    StoryID         int64               // if of the story
    Creator         sdk.AccAddress      // creator of vote
    Vote            bool                // value of vote
}
```
