# Bank module specification

## Summary

All transactions in TruStory will go through the bank module. It is responsible for keeping a ledger of all transactions, type of transactions, and preventing earned trustake from being traded.

## State

```go
type Transaction struct {
    ID              int64
    TransactionType TransactionType
    GroupID         int64
    ReferenceID     int64
    Amount          sdk.Coin
    Creator         sdk.AccAddress
    CreatedTime     time.Time
}

type TransactionType int8

const (
    Backing TransactionType = iota
    BackingReturned
    Challenge
    ChallengeReturned
    Upvote
    UpvoteReturned
    Interest
    InviteAFriend
    RewardPayout
)
```

## State Transitions
### Messages

`PayRewardMsg` is used as a reward for inviting friends to TruStory.

```go
type PayRewardMsg struct {
    Creator   sdk.AccAddress
    Recipient sdk.AccAddress
    Reward    sdk.Coin
    InviteID  int64
}
```
Currently the bank module doesn't allow transfer out of TruStory.
