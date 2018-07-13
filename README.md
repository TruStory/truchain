# TruStory Cosmos app

### State

```go
type Story struct {
    Body            string          // body of story
    Creator         sdk.Address     // address of creator
    SubmitBlock     int64           // block at which voting begins
    State           string          // "approved", "rejected", etc.
    YesVotes        int64           // total yes votes
    NoVotes         int64           // total no votes
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
