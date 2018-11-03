# TruChain: Story Module

![](dep.png)

## Keeper

### Dependencies
* category keeper

### Stores
* "stories"
    * keys
        * `"stories:id:[StoryID]"` -> `Story`
        * `"categories:id:[CategoryID]:stories:time:[time.Time]"` -> `[StoryID]`
        * `"categories:id:[CategoryID]:stories:challenged:time:[time.Time]"` -> `[StoryID]`
        * `"stories:len"` -> `[int64]`

## Notes
