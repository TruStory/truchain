# Slashing module specification

## Summary

The slashing module applies punishments to users who act badly on TruStory. After a user has been slashed a certain number of times, they are *jailed*. A user *unjails* themselves by creating a valid claim.

## State

### Main Types

```go
type Slash struct {
    ID          int64
    StakeID     int64
    Creator     sdk.AccAddress
}


// Params can be changed by governance vote
type Param struct {
    MaxStakeSlashCount   int
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

`SlashUserMsg` increments the `SlashCount` of an `Argument`. If `SlashCount` exceeds the `MaxStakeSlashCount` param, implement punishments. Only a user with an earned trustake of greater than 100 can slash. In the future, this value will be based on the total earned trustake in the community and user reputation.

Punishment
* Slash total interest of each staker
* Slash 2x the total stake amount of each staker

Curator reward
* Each user who marked "Unhelpful" will get a reward of 25% of the staking pool, distributed evenly

When a user is punished, their stake should be removed from the `ActiveStakes` queue since it should no longer expire. Also, their `SlashCount` should be incremented. If it exceeds the value defined in the `AppAccount` params, mark the user as "jailed". The user has to create a claim to "unjail" themselves.

```go
type SlashUserMsg struct {
    StakeID     int64
    Type        int
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