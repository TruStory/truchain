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
    Challenge
    Upvote
    BackingReturned
    ChallengeReturned
    UpvoteReturned
    RewardPool
    Interest
)
```

## State Transitions
### Messages

None.

Currently the bank module is only used internally and doesn't allow transfer in or out of TruStory.

### Functions
* AddCoin
* SubtractCoin
* MintAndAddCoin

