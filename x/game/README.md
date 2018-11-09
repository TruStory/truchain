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

This module doesn't have a codec and doesn't handle any messages. It is only used internally by TruChain to manage the validation game.
