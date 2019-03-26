# TruChain: Challenge Module

![](dep.png)

## Keeper

### Stores
* "challenges"
    *  keyspace
        * `"challenges:id:5"` -> `Challenge`
        * `"stories:id:[StoryID]:challenges:user:[AccAddress]"` -> `[ChallengeID]`
            * mapping of challengers for each story

## Notes

