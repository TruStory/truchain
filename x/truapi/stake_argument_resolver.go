package truapi

import (
	"context"
	"encoding/json"
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
)

func (ta *TruAPI) backingArgumentResolver(ctx context.Context, q app.QueryByIDParams) argument.Argument {
	res := ta.RunQuery("backings/id", app.QueryByIDParams{ID: q.ID})
	if res.Code != 0 {
		fmt.Println("error getting backing", res)
		return argument.Argument{}
	}
	backing := backing.Backing{}
	err := json.Unmarshal(res.Value, &backing)
	if err != nil {
		panic(err)
	}
	return ta.argumentResolver(ctx, app.QueryArgumentByID{ID: backing.ArgumentID, Raw: true})
}

func (ta *TruAPI) challengeArgumentResolver(ctx context.Context, q app.QueryByIDParams) argument.Argument {
	res := ta.RunQuery("challenges/id", app.QueryByIDParams{ID: q.ID})
	argument := argument.Argument{}
	if res.Code != 0 {
		fmt.Println("error getting challenge", res)
		return argument
	}
	challenge := challenge.Challenge{}
	err := json.Unmarshal(res.Value, &challenge)
	if err != nil {
		panic(err)
	}
	return ta.argumentResolver(ctx, app.QueryArgumentByID{ID: challenge.ArgumentID, Raw: true})
}

func (ta *TruAPI) stakeArgumentResolver(
	ctx context.Context, q app.QueryStakeArgumentByIDAndType) StakeArgument {
	if q.Backing {
		return StakeArgument{ta.backingArgumentResolver(ctx, app.QueryByIDParams{ID: q.StakeID})}
	}
	return StakeArgument{ta.challengeArgumentResolver(ctx, app.QueryByIDParams{ID: q.StakeID})}
}
