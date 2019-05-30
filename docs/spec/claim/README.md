# Claim module specification

## Summary

The claim module's responsibility is to create and store claims. Claims are the top-most level content in TruStory. Claims are short pieces of text that provide content for the start of a debate.

## State

### Main Types

`Claim` objects are used to track the life-cycle of a claim. They contain the content of a claim, and other fields which track the mutable state of debating a claim.

`TotalBacked` and `TotalChallenged` are modified by the staking and slashing modules when an argument is created or curated.

```go
type Claim struct {
    ID                  int64
    CommunityID         int64
    Body                string
    Creator             sdk.AccAddress
    Source              url.URL
    TotalParticipants   int64
    TotalBackingStake   sdk.Coin
    TotalChallengeStake sdk.Coin
    CreatedTime         time.Time
}

// Params can be voted on by governance
type Params struct {
    MinBodyLength int
    MaxBodyLength int
}
```

### Associations

`ClaimArguments` maintains a list of all arguments on a claim.

```go
// "claim:id:XX:argument:id:XX:creator:XX"
type ClaimArguments types.UserList
```

## State Transitions
### Messages

`CreateClaimMsg` creates a claim in the module's key-value store. 

When creating a claim, check if the `Creator` has been jailed. Jailed users cannot create claims.

```go
type CreateClaimMsg struct {
    CommunityID     int64
    Body            string
    Creator         sdk.AccAddress
}
```

A claim can be deleted with `DeleteClaimMsg` as long as `TotalBacked` and `TotalChallenged` are zero. The `Creator` of this message must be the same creator of the claim to delete.

```go
type DeleteClaimMsg struct {
    ClaimID     int64
    Creator     sdk.AccAddress
}
```
