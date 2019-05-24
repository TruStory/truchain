# Staking module specification

## Summary

The staking module contains all data types and state transitions needed to stake on Arguments. Arguments are how claims are supported or refuted. Arguments contain a text body that serves to *Back* or *Challenge* a specific claim with a specific amount of trustake.

Furthermore, Arguments can be staked with an *Agree* designation. It serves as kind of an upvote, to further enhance the standing of a claim, making the case for it stronger.

## State

### Main Types

The `Stake` type stores all data necessarily for a staking action. It is designed to be embedded by every type that requires staking functionality. For example, it is embedded in the `Argument` type.

```go
type Stake struct {
    ID           int64
    Type         int
    Amount       sdk.Coin
    Creator      sdk.AccAddress
    CreatedTime  time.Time
    EndTime      time.Time
}

// stake type enum
type StakeType int
const (
    Back StakeType = iota  // 0
    Challenge              // 1
    Agree                  // 2
)

// Params can be voted on by governance
type Params struct {
    Period                  time.Time   // default = 3 days
    ArgumentCreationStake   sdk.Coin    // default = 50 trustake
    AgreeStake              sdk.Coin    // default = 10 trustake
}
```

An `Argument` contains all data for an argument that either supports (back) or refutes (challenge) a claim. It also embeds a `Stake` type for staking data storage.

```go
type Argument struct {
    *Stake

    ClaimID              int64
    TLDR                 string
    Body                 string
    TotalAgreed          sdk.Coin
    SlashCount           int
}
```

### Associations

`ArgumentStakes` maintains an easily accessible list of all stakes for each argument.

```go
// "argument:id:XX:staketype:XX:creator:XX" -> nil
type ArgumentStakes app.UserList
```

### Queues

`ActiveStakes` maintains a queue of all currently active stakes, sorted by `EndTime`.

```go
type ActiveStakes Queue
```

## State Transitions
### Messages

`CreateArgumentMsg` creates an `Argument` in the module's key-value store. The only allowed values for `StakeType` are 0 (back), and 1 (challenge). 

`EndTime` is currently fixed at 3 days. 

The stake `Amount` is currently fixed at 50 trustake. In the future, this will be a value algorithmically determined based on various factors such as the current amount staked on the claim, total supply of trustake, and the health of the community associated with the claim.

```go
type CreateArgumentMsg struct {
    ClaimID       int64
    TLDR          string
    Body          string
    StakeType     int             // back  or challenge
    Creator       sdk.AccAddress
}
```

An argument currently cannot be edited or deleted.

An argument's standing can be enhanced with an `AgreeArgumentMsg` with some stake. 

`EndTime` is currently fixed at 3 days. 

The stake `Amount` is currently fixed at 10 trustake.

```go
type AgreeArgumentMsg struct {
    ArgumentID    int64
    Creator       sdk.AccAddress
}
```

Staking via `CreateArgumentMsg` and `AgreeAgrumentMsg` should fail validation if the creator has already staked over 66% of their total trustake within a 7-day rolling period. 

## Block Triggers

### End Block

After each block is processed, check the `ActiveStakes` queue for expiring stakes. After a stake has ended, distribute rewards.

Rewards:
* argument creators get 70% interest reward from each staker
* stakers keep 30% interest

This incentive structure heavily rewards argument creation. Agreeing is a lightweight way to earn 30% interest. But to earn full interest and rewards, content creators are encouraged to write arguments.
