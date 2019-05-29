# Community module specification

The community module specifies data for a TruStory community.

## State

```go
type Community struct {
    ID                  int64         
    Name                string        
    Slug                string        
    Description         string        
    TotalEarnedStake    sdk.Coin      
}
```

## State Transitions

Currently, communities are created at genesis. The only way to add more is via import/export. In the future they may be created via governance vote or by reputable users.
