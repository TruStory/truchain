package story

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// // GenesisState struct
// type GenesisState struct {
// 	Stories []GenesisStory `json:"stories`
// }

// // NewGenesisState
// func NewGenesisState(stories []GenesisStory) GenesisState {
// 	return GenesisState{
// 		Stories: stories,
// 	}
// }

// type GenesisStory struct {
// 	ID         int64          `json:"id"`
// 	Argument   string         `json:"arguments,omitempty"`
// 	Body       string         `json:"body"`
// 	CategoryID int64          `json:"category_id"`
// 	Creator    sdk.AccAddress `json:"creator"`
// 	Flagged    bool           `json:"flagged,omitempty"`
// 	GameID     int64          `json:"game_id,omitempty"`
// 	Source     url.URL        `json:"source,omitempty"`
// 	State      State          `json:"state"`
// 	Type       Type           `json:"type"`
// 	Timestamp  app.Timestamp  `json:"timestamp"`
// }

// // NewGenesisStoryI gets the sate addresses and cins
// func NewGenesisStoryI(st Story) GenesisStory {
// 	gst := GenesisStory{
// 		ID:         st.ID,
// 		Argument:   st.Argument,
// 		Body:       st.Body,
// 		CategoryID: st.CategoryID,
// 		Creator:    st.Creator,
// 		Flagged:    st.Flagged,
// 		GameID:     st.GameID,
// 		Source:     st.Source,
// 		State:      st.State,
// 		Type:       st.Type,
// 		Timestamp:  st.Timestamp,
// 	}

// 	return gst
// }

// ExportGenesis gets all the current stories and calls app.WriteJSONtoNodeHome() to write data to file.
// []GenesisStory
func ExportGenesis(ctx sdk.Context, sk WriteKeeper) []Story {

	// st := Story{}
	// fmt.Printf("%+v\n", st)
	// stories := []GenesisStory{}
	// appendStoriesfn := func(st Story) bool {
	// 	story := NewGenesisStoryI(st)
	// 	stories = append(stories, story)
	// 	return false
	// }
	// fmt.Printf("%+v\n", stories)
	s := sk.StoriesNoSort(ctx)
	// fmt.Printf("%+v\n", s)

	// sk.IterateStories(ctx, appendStoriesfn)
	// s := sk.StoriesNoSort(ctx)
	// stories = append(stories, s)
	// app.WriteJSONtoNodeHome(stories, dnh, bh, fmt.Sprintf("%s.json", k.GetStoreKey().Name()))
	return s
}
