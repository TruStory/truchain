package truapi

import (
  "context"
  "encoding/json"
  "fmt"

  "github.com/TruStory/truchain/x/story"
)

func (ta *TruApi) storyResolver(_ context.Context, q story.QueryStoriesByIdParams) []story.Story {
  res := ta.RunQuery("story", q)

  if res.Code != 0 {
    fmt.Println("Resolver err: ", res)
    return []story.Story{}
  }

  s := new([]story.Story)
  _ = json.Unmarshal(res.Value, s)

  return *s
}
