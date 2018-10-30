# TruChain: Story Module

![](dep.png)

## Keeper

### Dependencies
* category keeper

### Stores
* "stories"
    * keys
        * `"stories:id:[StoryID]"` -> `Story`

* "storiesByCategory"
    * keys
        * `"categories:id:[CategoryID]:stories:id:[StoryID]"` -> `[StoryID]`

* "challengedStoriesByCategory"
    * keys
        * `"categories:id:[CategoryID]:stories:id:[StoryID]"` -> `[StoryID]`

## Notes


