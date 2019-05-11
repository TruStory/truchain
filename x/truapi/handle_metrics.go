package truapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/truapi/render"
	trubank "github.com/TruStory/truchain/x/trubank"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Metrics represents metrics for the platform.
type Metrics struct {
	Users map[string]*UserMetrics `json:"users"`
}

// AccumulatedUserCred tracks accumulated cred by day.
type AccumulatedUserCred map[string]sdk.Coin

// GetUserMetrics gets user metrics or initializes one if not in the map.
func (m *Metrics) GetUserMetrics(address string) *UserMetrics {
	userMetrics, ok := m.Users[address]
	if !ok {
		userMetrics = &UserMetrics{
			InterestEarned:     sdk.NewCoin(app.StakeDenom, sdk.NewInt(0)),
			StakeLost:          sdk.NewCoin(app.StakeDenom, sdk.NewInt(0)),
			StakeEarned:        sdk.NewCoin(app.StakeDenom, sdk.NewInt(0)),
			TotalAmountAtStake: sdk.NewCoin(app.StakeDenom, sdk.NewInt(0)),
			TotalAmountStaked:  sdk.NewCoin(app.StakeDenom, sdk.NewInt(0)),
			CredEarned:         make(map[string]AccumulatedUserCred),
		}
	}
	m.setUserMetrics(address, userMetrics)
	return userMetrics
}

func (m *Metrics) setUserMetrics(address string, userMetrics *UserMetrics) {
	m.Users[address] = userMetrics
}

// UserMetrics a summary of different metrics per user
type UserMetrics struct {
	TotalClaims            int64 `json:"total_claims"`
	TotalArguments         int64 `json:"total_arguments"`
	TotalGivenEndorsements int64 `json:"total_given_endorsments"`
	// Tracked by day
	CredEarned         map[string]AccumulatedUserCred
	InterestEarned     sdk.Coin `json:"intereset_earned"`
	StakeLost          sdk.Coin `json:"stake_lost"`
	StakeEarned        sdk.Coin `json:"stake_earned"`
	TotalAmountAtStake sdk.Coin `json:"total_amount_at_stake"`
	TotalAmountStaked  sdk.Coin `json:"total_amount_staked"`
}

func (um *UserMetrics) increaseArgumentsCount() {
	um.TotalArguments = um.TotalArguments + 1
}
func (um *UserMetrics) increaseClaimsCount() {
	um.TotalClaims = um.TotalClaims + 1
}

func (um *UserMetrics) increaseGivenEndorsmentsCount() {
	um.TotalGivenEndorsements = um.TotalGivenEndorsements + 1
}

func (um *UserMetrics) addAmoutAtStake(amount sdk.Coin) {
	um.TotalAmountAtStake = um.TotalAmountAtStake.Plus(amount)
}

func (um *UserMetrics) addStakedAmount(amount sdk.Coin) {
	um.TotalAmountStaked = um.TotalAmountStaked.Plus(amount)
}

func (um *UserMetrics) addStakeLost(amount sdk.Coin) {
	um.StakeLost = um.StakeLost.Plus(amount)
}

func (um *UserMetrics) addInterestEarned(amount sdk.Coin) {
	um.InterestEarned = um.InterestEarned.Plus(amount)
}

func (um *UserMetrics) addStakeEarned(amount sdk.Coin) {
	um.StakeEarned = um.StakeEarned.Plus(amount)
}

func (um *UserMetrics) trackCreadEarned(tx trubank.Transaction) {
	date := tx.Timestamp.CreatedTime.Format("2006-01-02")
	userCred, ok := um.CredEarned[date]
	if !ok {
		userCred = AccumulatedUserCred(make(map[string]sdk.Coin))
	}
	um.CredEarned[date] = userCred
	categoryTotal, ok := userCred[tx.Amount.Denom]

	if !ok {
		categoryTotal = sdk.NewCoin(tx.Amount.Denom, sdk.NewInt(0))
	}

	userCred[tx.Amount.Denom] = categoryTotal.Plus(tx.Amount)
}

// HandleMetrics dumps metrics per user basis.
func (ta *TruAPI) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	stories := make([]story.Story, 0)

	res := ta.RunQuery("stories/all", struct{}{})
	err := json.Unmarshal(res.Value, &stories)
	if err != nil {
		render.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	m := &Metrics{
		Users: make(map[string]*UserMetrics),
	}
	mappedStories := make(map[int64]int)
	mapUserStakeByStoryIDKey := func(user string, storyID int64) string {
		return fmt.Sprintf("%s:%d", user, storyID)
	}
	mapUserStakeByStoryID := make(map[string]sdk.Coin)
	for idx, s := range stories {
		mappedStories[s.ID] = idx
		backingAmount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(0))
		challengeAmount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(0))
		m.GetUserMetrics(s.Creator.String()).increaseClaimsCount()

		// get backings and challenges
		backings := ta.backingsResolver(r.Context(), app.QueryByIDParams{ID: s.ID})
		for _, b := range backings {
			backingAmount = backingAmount.Plus(b.Amount())
			creator := b.Creator().String()
			mapUserStakeByStoryID[mapUserStakeByStoryIDKey(creator, b.StoryID())] = b.Amount()
			if b.StoryID() == 235 {
				fmt.Println("backed", mapUserStakeByStoryIDKey(creator, b.StoryID()), b.Amount())
			}
			backerMetrics := m.GetUserMetrics(creator)
			backerMetrics.addStakedAmount(b.Amount())
			if s.Status == story.Pending {
				backerMetrics.addAmoutAtStake(b.Amount())
			}
			argument := ta.argumentResolver(r.Context(), app.QueryArgumentByID{ID: b.ArgumentID, Raw: true})

			if argument.ID == 0 {
				continue
			}

			if argument.Creator.String() == creator {
				backerMetrics.increaseArgumentsCount()
			}

			if argument.Creator.String() != creator {
				backerMetrics.increaseGivenEndorsmentsCount()
			}

		}

		challenges := ta.challengesResolver(r.Context(), app.QueryByIDParams{ID: s.ID})
		for _, c := range challenges {
			challengeAmount = challengeAmount.Plus(c.Amount())
			creator := c.Creator().String()
			if c.StoryID() == 235 {
				fmt.Println("challenge", mapUserStakeByStoryIDKey(creator, c.StoryID()), c.Amount())
			}
			mapUserStakeByStoryID[mapUserStakeByStoryIDKey(creator, c.StoryID())] = c.Amount()
			challengerMetrics := m.GetUserMetrics(creator)
			challengerMetrics.addStakedAmount(c.Amount())
			if s.Status == story.Pending {
				challengerMetrics.addAmoutAtStake(c.Amount())
			}

			argument := ta.argumentResolver(r.Context(), app.QueryArgumentByID{ID: c.ArgumentID, Raw: true})

			if argument.ID == 0 {
				continue
			}

			if argument.Creator.String() == creator {
				challengerMetrics.increaseArgumentsCount()
			}

			if argument.Creator.String() != creator {
				challengerMetrics.increaseGivenEndorsmentsCount()
			}
		}
		// only check expired
		if s.Status == story.Pending {
			continue
		}
		// Check if backings lost
		if backingAmount.IsLT(challengeAmount) {
			for _, b := range backings {
				m.GetUserMetrics(b.Creator().String()).addStakeLost(b.Amount())
			}
		}

		// Check if challenges lost
		if challengeAmount.IsLT(backingAmount) {
			for _, c := range challenges {
				m.GetUserMetrics(c.Creator().String()).addStakeLost(c.Amount())
			}
		}
	}

	type storyRewardResult struct {
		Reward        *sdk.Coin
		StakeReturned *sdk.Coin
	}
	for userAddress, userMetrics := range m.Users {
		txs := ta.transactionsResolver(r.Context(), app.QueryByCreatorParams{Creator: userAddress})
		userStoryResults := make(map[int64]*storyRewardResult)
		for _, tx := range txs {
			switch tx.TransactionType {
			case trubank.Interest:
				userMetrics.addInterestEarned(tx.Amount)
			case trubank.BackingLike:
				userMetrics.trackCreadEarned(tx)
			case trubank.ChallengeLike:
				userMetrics.trackCreadEarned(tx)
			// this three transactions are related to finished expired stories.
			case trubank.RewardPool:
				fallthrough
			case trubank.BackingReturned:
				fallthrough
			case trubank.ChallengeReturned:
				i, ok := mappedStories[tx.GroupID]
				if !ok {
					continue
				}
				s := stories[i]
				if s.Status != story.Expired {
					continue
				}
				storyReward, ok := userStoryResults[tx.GroupID]
				if !ok {
					storyReward = &storyRewardResult{}
					userStoryResults[tx.GroupID] = storyReward
				}
				if tx.TransactionType == trubank.RewardPool {

					reward := sdk.NewCoin(tx.Amount.Denom, sdk.NewInt(tx.Amount.Amount.Int64()))
					storyReward.Reward = &reward
				}

				if tx.TransactionType == trubank.BackingReturned || tx.TransactionType == trubank.ChallengeReturned {
					returned := sdk.NewCoin(tx.Amount.Denom, sdk.NewInt(tx.Amount.Amount.Int64()))
					storyReward.StakeReturned = &returned
				}
			}
		}

		for storyID, storyResult := range userStoryResults {
			// majority was not reached and we performed a refund
			if storyResult.Reward == nil {
				continue
			}
			// this is the case after we introduced two transactions to reward an user
			if storyResult.StakeReturned != nil {
				userMetrics.addStakeEarned(*storyResult.Reward)
				continue
			}
			// this will be the case where we will need to deduct staked amount from reward to get net value
			stakedAmount, ok := mapUserStakeByStoryID[mapUserStakeByStoryIDKey(userAddress, storyID)]
			if !ok {
				continue
			}
			reward := storyResult.Reward.Minus(stakedAmount)
			// stake was returned
			if reward.IsZero() {
				continue
			}

			// this should not happend for any reason but just adding a safe check point.
			if reward.IsNegative() {
				continue
			}
			userMetrics.addStakeEarned(reward)

		}
	}
	render.JSON(w, r, m, http.StatusOK)

}
