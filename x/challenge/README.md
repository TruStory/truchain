# TruChain: Challenge Module

![](dep.png)

## Keeper

### Dependencies
* story keeper

### Stores
* "challenges"
    *  keyspace
        * `"challenges:id:5"` -> `Challenge`
        * `"games:id:[GameID]:challenges:user:[AccAddress]"` -> `[ChallengeID]`
            * mapping of challenges for each game

## Notes

