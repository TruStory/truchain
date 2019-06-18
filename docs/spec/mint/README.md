# Mint module specification

The mint module is responsible for calculating the inflation per block.

Note: Much of this module is based on Cosmos SDK's [mint module](https://github.com/cosmos/cosmos-sdk/tree/master/x/mint).

## State

The minter holds current inflation information.

```go
type Minter struct {
   Inflation            sdk.Dec  // annual inflation rate
   AnnualProvisions     sdk.Dec  // current annual expected provisions
}
```

Parameters for inflation calculations.

```go
type Params struct {
    MintDenom               string
    InflationRateChange     sdk.Dec  // 15%
    InflationMax            sdk.Dec  // 25%
    InflationMin            sdk.Dec  // 10%
    GoalStaked              sdk.Dec  // 67%
    BlocksPerYear           uint64   // 5,256,000 for 6 second blocks
}
```

BlocksPerYear = 60 min / 6 sec block time = 10 blocks per min * 60 min = 600 per hour * 24 hours = 14,400 per day * 365 days = 5,256,000 

## Block Triggers

### Begin-Block

At the beginning of each block, inflation is re-calcuated.

#### Inflation Rate

The target annual inflation rate is re-calculated each block. It changes based on the distance from a target desired ratio of 67%. If the staked ratio is below 67%, the inflation rate approaches a cap of 25% to encourage more staking. If the ratio is above 67%, the rate approaches a cap of 10%. The maximum possible annual rate change is 15% per year.

```go
NextInflationRate(params Params, stakedRatio sdk.Dec) (inflation sdk.Dec) {
    inflationRateChangePerYear = (1 - stakedRatio / Params.GoalStaked) * Params.InflationRateChange
    inflationRateChange = inflationRateChangePerYear / Params.BlocksPerYear

    inflation += inflationRateChange
    if inflation > Params.InflationMax {
        inflation = Params.InflationMax
    }
    if inflation < Params.InflationMin {
        inflation  = Params.InflationMin
    }

    return inflation
}
```

#### Annual Provisions

Calculate the annual provisions based on the current total supply and inflation rate.

```go
NextAnnualProvisions(totalSupply sdk.Dec) (provisions sdk.Dec) {
    return inflation + totalSupply
}
```

This value is saved, and later will be handled by the distribution module end block for actual distribution to the new user pool, staked users, and staked validators.

#### Block Provisions

Calculate the total provisions at the block level, based on the annual provisions.

```go
BlockProvisions(params Params) sdk.Coin {
    provisionAmount := AnnualProvisions / params.BocksPerYear
    return sdk.NewCoin(params.MintDenom, provisionAmount.Truncate())
}
```
