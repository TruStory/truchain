# Authentication module specification

The authentication module registers and authenticates a TruStory user. It wraps the Cosmos auth module and provides TruStory-related functionality.

## State

`EarnedCoin` is a representation of a coin associated with a community, or "earned trustake".

```go
type EarnedCoin struct {
    sdk.Coin

    CommunityID int64
}

type EarnedCoins []EarnedCoin
```

`AppAccount` is the main account for a TruStory user. It contains the user's address, coins, public key, as well as the slash count and jail status.

```go
type AppAccount struct {
    auth.BaseAccount

    EarnedStake     type.EarnedCoins
    SlashCount      int
    IsJailed        bool
    JailEndTime     time.Time
    CreatedTime     time.Time
}

// Params can be changed by governance vote
type Param struct {
    MaxSlashCount   int
}
```

For reference, `auth.BaseAccount` is:
```go
type BaseAccount struct {
	Address       sdk.AccAddress
	Coins         sdk.Coins
	PubKey        crypto.PubKey
	AccountNumber uint64       
	Sequence      uint64       
}
```

The `Coins` stored in the base account represent spendable coins, not earned.

## State Transitions
### Messages

`RegisterKeyMsg` registers creates a new `AppAccount` for the user.

```go
type RegisterKeyMsg struct {
    Address    sdk.AccAddress
    PubKey     crypto.PubKey
    PubKeyAlgo string 
    Coins      sdk.Coins
}
```
