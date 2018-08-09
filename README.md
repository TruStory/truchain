# TruStory Cosmos app

### Types

All state is stored in these types.

```go
type Bond struct {
    ID              int64               // id of the bond
    Amount          float64             // amount of the bond
    CreatedBlock    int64               // block at which bond was created
    Creator         sdk.AccAddress      // address of the creator
    Period          int64               // period of the bond    
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

- `SubmitStoryMsg`: to submit stories
- `VoteMsg`: to vote on stories

```go
type SubmitStoryMsg struct {
    Body            string          // body of story
    Creator         sdk.Address     // address of creator
}
```

```go
type VoteMsg struct {
    StoryID         int64           // id of the story
    Option          string          // "yes" or "no"
    Stake           sdk.Coins       // stake for vote
    Voter           sdk.Address     // address of voter
}
```
