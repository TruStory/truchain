package argument

// NewHandler creates a function to handle argument messages
// func NewHandler(k Keeper) sdk.Handler {
// 	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
// 		switch msg := msg.(type) {
// 		case LikeArgumentMsg:
// 			return handleLikeArgumentMsg(ctx, k, msg)
// 		default:
// 			return app.ErrMsgHandler(msg)
// 		}
// 	}
// }

// func handleLikeArgumentMsg(ctx sdk.Context, k Keeper, msg LikeArgumentMsg) sdk.Result {
// 	if err := msg.ValidateBasic(); err != nil {
// 		return err.Result()
// 	}

// 	id, err := k.Like(ctx, msg.ArgumentID, msg.Creator, msg.VoteChoice, msg.StakeID)
// 	if err != nil {
// 		return err.Result()
// 	}

// 	return app.Result(id)
// }
