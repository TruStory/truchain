package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// –––––––––––– Flags ––––––––––––––––

// nolint
const (
	FlagVerifiedStory = "verified"
)

// GetCmdQueryStories gets the command to get all stories
func GetCmdQueryStories(storeName string, cdc *codec.Codec) *cobra.Command {
	// cmd := &cobra.Command{
	// 	Use:   "stories",
	// 	Short: "Query all stories",
	// 	Args:  cobra.ExactArgs(0),
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		ctx := context.NewCoreContextFromViper()

	// 		resKVs, err := ctx.QuerySubspace(cdc, []byte("stories"), storeName)
	// 		if err != nil {
	// 			// return err
	// 		}

	// 		// if viper.IsSet(FlagVerifiedStory) {
	// 		// 	isVerified := viper.GetBool(FlagVerifiedStory)
	// 		// } else {
	// 		// 	isVerified := false
	// 		// }

	// 		// parse out the stories
	// 		var stories []ts.Story
	// 		for _, KV := range resKVs {
	// 			var story ts.Story
	// 			cdc.MustUnmarshalBinary(KV.Value, &story)
	// 			stories = append(stories, story)
	// 		}

	// 		output, err := codec.MarshalJSONIndent(cdc, stories)
	// 		if err != nil {
	// 			// return err
	// 		}
	// 		fmt.Println(string(output))
	// 		// return nil
	// 	},
	// }
	// cmd.Flags().Bool(FlagVerifiedStory, false, "Query only verified stories (default: true)")
	// return cmd
	return nil
}
