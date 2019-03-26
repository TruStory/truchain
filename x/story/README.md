# TruChain: Story Module

![](dep.png)

## Keeper

### Stores
* "stories"
    * keys
        * `"stories:id:[StoryID]"` -> `Story`
        * `"categories:id:[CategoryID]:stories:time:[time.Time]"` -> `[StoryID]`
        * `"stories:len"` -> `[int64]`

## Notes
