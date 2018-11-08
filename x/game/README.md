# TruChain: Game Module

![](dep.png)

## Keeper

Stores all data pertaining to a validation game.

### Dependencies
* story keeper

### Stores
* "games"
    * keys
        * `"games:id:[GameID]"` -> `Game`
        * `"games:len"` -> `[int64]`
        * `"stories:id:[StoryID]:games:time:[Time]"` -> `[GameID]`

## Notes
