# Slashing module specification

## Summary

The slashing module applies punishments to users who act badly on TruStory. After a user has been slashed a certain number of times, they are *jailed*.

## State

### Main Types

```go
type Slash struct {
    ID              uint64
    ArgumentID      uint64
    Creator         sdk.AccAddress
    CreatedTime     time.Time
}

// Params can be changed by governance vote
type Param struct {
    MinSlashCount        int
    SlashMagnitude       sdk.Dec            // 3x
    SlashMinStake        types.EarnedCoin   // 50 earned trustake
    SlashAdmins          []sdk.AccAddress   // list of admin addresses who can slash
    JailTime             time.Duration      // 7 days
}
```

## State Transitions

### Messages

`SlashArgumentMsg` increments the `SlashCount` of an `Argument`. If `SlashCount` exceeds the `MaxSlashCount` param, implement punishments. Only a user with an earned trustake of greater than `SlashMinStake` can slash *or* `SlashAdmins`. In the future, this value will be based on the total earned trustake in the community and user reputation.

Fail validation if the `SlashCount` already exceeds `MinSlashCount`, preventing further slashing on the argument.

If `SlashCount` is equal to `MinSlashCount`, then remove the amount of this stake from the total backing or challenge stake count on the claim.

Furthermore, the same user cannot slash the same argument more than once.

Punishment
* Slash total interest of each staker
* Slash 3x the total stake amount of each staker

A staker is someone who has backed, challenged, or upvoted (agreed).

Curator reward
* Each user who marked "Unhelpful" will get a reward of 25% of the staking pool, distributed evenly

When a user is punished, their stake should be removed from the `ActiveStakes` queue since it should no longer expire. Also, their `SlashCount` should be incremented. If it exceeds the value defined in the `AppAccount` params, mark the user as "jailed" and add set the `JailEndTime` on the user.

```go
type SlashArgumentMsg struct {
    StakeID     int64
    Type        SlashType
    Creator     sdk.AccAddress
}
```

The `SlashType` enum currently only has the `Unhelpful` state, but can be expanded with more slash types in the future.

```go
// slash type enum
type SlashType int
const (
    Unhelpful SlashType = iota  // 0
)

```
