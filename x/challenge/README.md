# TruChain: Challenge Module

![](dep.png)

## Keeper

### Dependencies
* story keeper

### Stores
* "challenges"
    *  keyspace
        * `"challenges:id:5"` -> `Challenge`
            * main type storage
        * `"challenges:id:5:userAddr:0xdeadbeef"` -> `sdk.Coin`
            * list of challengers for each challenge
        * `"0x00"` -> `uint64`
            * unexpired challenge queue length
        * `"0x0100000000000000000001"` -> `int64`
            * unexpired challenge queue storage

## Notes

