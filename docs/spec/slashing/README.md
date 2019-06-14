# Slashing module specification

## Summary

The slashing module applies punishments to users who act badly on TruStory. After a user has been slashed a certain number of times, they are _jailed_. A user _unjails_ themselves by creating a valid claim.

## State

### Main Types

```go
type Slash struct {
    ID              int64
    StakeID         int64
    Creator         sdk.AccAddress
    CreatedTime     time.Time
}


// Params can be changed by governance vote
type Param struct {
    MaxStakeSlashCount   int
    SlashMagnitude       sdk.Dec            // 3x
    SlashMinStake        types.EarnedCoin   // 50 earned trustake
    SlashAdmins          []sdk.AccAddress   // list of admin addresses who can slash
    JailTime             time.Duration      // 7 days
}
```

### Associations

`SlashedStakes` maintains an easily accessible list of all slashes for each stake.

```go
// "stake:id:XX:slash:id:XX:creator:XX" -> nil
type SlashedStakes app.UserList
```

## State Transitions

### Messages

`SlashArgumentMsg` increments the `SlashCount` of an `Argument`. If `SlashCount` exceeds the `MaxStakeSlashCount` param, implement punishments. Only a user with an earned trustake of greater than `SlashMinStake` can slash _or_ `SlashAdmins`. In the future, this value will be based on the total earned trustake in the community and user reputation.

Fail validation if the `SlashCount` already exceeds `MaxStakeSlashCount`, preventing further slashing on the argument.

If `SlashCount` is equal to `MaxStakeSlashCount`, then remove the amount of this stake from the total backing or challenge stake count on the claim.

Futhermore, the same user cannot slash the same argument more than once. Verify with the `SlashedStakes` list first.

Punishment

- Slash total interest of each staker
- Slash 3x the total stake amount of each staker

Curator reward

- Each user who marked "Unhelpful" will get a reward of 25% of the staking pool, distributed evenly

When a user is punished, their stake should be removed from the `ActiveStakes` queue since it should no longer expire. Also, their `SlashCount` should be incremented. If it exceeds the value defined in the `AppAccount` params, mark the user as "jailed" and add set the `JailEndTime` on the user. The user has to create an argument with at least X upvotes to unjail themselves. The parameter for X is defined in the staking module param store.

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
