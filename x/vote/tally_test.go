package vote

// func TestProcessGame(t *testing.T) {
// 	ctx, _, k := fakeValidationGame()

// 	gameID := int64(1)
// 	game, _ := k.gameKeeper.Game(ctx, gameID)

// 	err := processGame(ctx, k, game)
// 	assert.Nil(t, err)
// }

// func TestTally(t *testing.T) {
// 	ctx, _, k := fakeValidationGame()

// 	gameID := int64(1)
// 	game, _ := k.gameKeeper.Game(ctx, gameID)

// 	votes, _ := tally(ctx, k, game)

// 	assert.Equal(t, 6, len(votes.trueVotes))
// 	assert.Equal(t, 5, len(votes.falseVotes))
// }

// func TestRewardPool(t *testing.T) {
// 	ctx, votes, _ := fakeValidationGame()

// 	expectedPool := sdk.NewCoin(params.StakeDenom, sdk.NewInt(5500000000000))

// 	pool, _ := rewardPool(ctx, votes, true)
// 	assert.Equal(t, expectedPool.String(), pool.String())
// }

// func TestDistributeRewards(t *testing.T) {
// 	ctx, votes, k := fakeValidationGame()

// 	pool := sdk.NewCoin(params.StakeDenom, sdk.NewInt(35))

// 	err := distributeRewards(
// 		ctx, k.backingKeeper, k.bankKeeper, pool, votes, true, "trudex")
// 	assert.Nil(t, err)
// }

// func TestConfirmStory(t *testing.T) {
// 	ctx, votes, k := fakeValidationGame()

// 	confirmed, _ := confirmStory(ctx, k.accountKeeper, votes, "trudex")
// 	assert.True(t, confirmed)
// }

// func TestWeightedVote(t *testing.T) {
// 	ctx, votes, k := fakeValidationGame()

// 	trueWeights, _ := weightedVote(ctx, k.accountKeeper, votes.trueVotes, "trudex")
// 	falseWeights, _ := weightedVote(ctx, k.accountKeeper, votes.falseVotes, "trudex")

// 	// 5 true, 1 preethi added each due to cold-start
// 	assert.Equal(t, "6", trueWeights.String())

// 	// 4 false, 1 preethi added each due to cold-start
// 	assert.Equal(t, "5", falseWeights.String())
// }

// func TestConfirmedStoryRewardPool(t *testing.T) {
// 	ctx, votes, _ := fakeValidationGame()

// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())

// 	confirmedPool(ctx, votes.falseVotes, &pool)
// 	assert.Equal(t, "5500000000000trusteak", pool.String())
// }

// func TestConfirmedStoryRewardPool2(t *testing.T) {
// 	ctx, votes, _ := fakeValidationGame2()

// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())

// 	confirmedPool(ctx, votes.falseVotes, &pool)
// 	assert.Equal(t, "265000000000trusteak", pool.String())
// }

// func TestDistributeRewardsConfirmed(t *testing.T) {
// 	ctx, votes, k := fakeValidationGame()
// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())
// 	confirmedPool(ctx, votes.falseVotes, &pool)

// 	cred := "trudex"

// 	err := distributeRewardsConfirmed(
// 		ctx, k.backingKeeper, k.bankKeeper, votes, pool, cred)
// 	assert.Nil(t, err)

// 	coins := sdk.Coins{}

// 	winningBacker1 := votes.trueVotes[0].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "1416666666666", coins.AmountOf(cred).String())

// 	winningBacker2 := votes.trueVotes[1].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, winningBacker2.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "1416666666666", coins.AmountOf(cred).String())

// 	winningBacker3 := votes.trueVotes[2].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, winningBacker3.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "1416666666666", coins.AmountOf(cred).String())

// 	winningVoter1 := votes.trueVotes[3].(TokenVote)
// 	coins = k.bankKeeper.GetCoins(ctx, winningVoter1.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "916666666666", coins.AmountOf(cred).String())

// 	winningVoter2 := votes.trueVotes[4].(TokenVote)
// 	coins = k.bankKeeper.GetCoins(ctx, winningVoter2.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "916666666666", coins.AmountOf(cred).String())

// 	losingBacker1 := votes.falseVotes[0].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, losingBacker1.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingChallenger1 := votes.falseVotes[1].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, losingChallenger1.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingChallenger2 := votes.falseVotes[2].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, losingChallenger2.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingChallenger3 := votes.falseVotes[3].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, losingChallenger3.Creator())
// 	assert.Equal(t, "0", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingVoter1 := votes.falseVotes[4].(TokenVote)
// 	coins = k.bankKeeper.GetCoins(ctx, losingVoter1.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())
// }

// func TestDistributeRewardsConfirmed2(t *testing.T) {
// 	ctx, votes, k := fakeValidationGame2()
// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())
// 	confirmedPool(ctx, votes.falseVotes, &pool)

// 	cred := "trudex"

// 	err := distributeRewardsConfirmed(
// 		ctx, k.backingKeeper, k.bankKeeper, votes, pool, cred)
// 	assert.Nil(t, err)

// 	coins := sdk.Coins{}

// 	winningBacker1 := votes.trueVotes[0].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "95003666333", coins.AmountOf(cred).String())

// 	winningBacker2 := votes.trueVotes[1].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, winningBacker2.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "92001934065", coins.AmountOf(cred).String())

// 	winningBacker3 := votes.trueVotes[2].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, winningBacker3.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "95003666333", coins.AmountOf(cred).String())

// 	losingChallenger1 := votes.falseVotes[0].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, losingChallenger1.Creator())
// 	assert.Equal(t, "1828000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingChallenger2 := votes.falseVotes[1].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, losingChallenger2.Creator())
// 	assert.Equal(t, "1917000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingChallenger3 := votes.falseVotes[2].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, losingChallenger3.Creator())
// 	assert.Equal(t, "1990000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())
// }

// func TestRejectedStoryRewardPool(t *testing.T) {
// 	ctx, votes, _ := fakeValidationGame()

// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())

// 	rejectedPool(ctx, votes, &pool)
// 	assert.Equal(t, "8000000000000trusteak", pool.String())
// }

// func TestChallengerPool(t *testing.T) {
// 	ctx, votes, _ := fakeValidationGame()
// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())
// 	rejectedPool(ctx, votes, &pool)

// 	coin := calculateChallengerPool(pool, DefaultParams())
// 	assert.Equal(t, "6000000000000trusteak", coin.String())
// }

// func TestVoterPool(t *testing.T) {
// 	ctx, votes, _ := fakeValidationGame()
// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())
// 	rejectedPool(ctx, votes, &pool)

// 	coin := calculateVoterPool(pool, DefaultParams())
// 	assert.Equal(t, "2000000000000trusteak", coin.String())
// }

// func TestCount(t *testing.T) {
// 	ctx, votes, _ := fakeValidationGame()
// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())
// 	rejectedPool(ctx, votes, &pool)

// 	cAmount, cCount, vCount, _ := winnerInfo(votes.falseVotes)
// 	assert.Equal(t, int64(4000000000000), cAmount.Int64())
// 	assert.Equal(t, int64(3), cCount)
// 	assert.Equal(t, int64(1), vCount)
// }

// func TestChallengerRewardAmount(t *testing.T) {
// 	coin := challengerRewardAmount(
// 		sdk.NewCoin("trudex", sdk.NewInt(1000000000000)),
// 		sdk.NewInt(2000000000000),
// 		sdk.NewCoin("trudex", sdk.NewInt(5250000000000)))

// 	assert.Equal(t, "2625000000000", coin.String())
// }

// func TestChallengerRewardAmount2(t *testing.T) {
// 	coin := challengerRewardAmount(
// 		sdk.NewCoin("trudex", sdk.NewInt(1500000000000)),
// 		sdk.NewInt(2000000000000),
// 		sdk.NewCoin("trudex", sdk.NewInt(5250000000000)))

// 	assert.Equal(t, "3937500000000", coin.String())
// }

// func TestChallengerRewardAmount3(t *testing.T) {
// 	coin := challengerRewardAmount(
// 		sdk.NewCoin("trudex", sdk.NewInt(500000000000)),
// 		sdk.NewInt(2000000000000),
// 		sdk.NewCoin("trudex", sdk.NewInt(5250000000000)))

// 	assert.Equal(t, "1312500000000", coin.String())
// }

// func TestDistributeRewardsRejected(t *testing.T) {
// 	ctx, votes, k := fakeValidationGame()
// 	pool := sdk.NewCoin(params.StakeDenom, sdk.ZeroInt())
// 	rejectedPool(ctx, votes, &pool)

// 	cred := "trudex"

// 	err := distributeRewardsRejected(
// 		ctx, k.backingKeeper, k.bankKeeper, votes, pool, cred)
// 	assert.Nil(t, err)

// 	coins := sdk.Coins{}

// 	winningBacker1 := votes.falseVotes[0].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())

// 	winningChallenger1 := votes.falseVotes[1].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, winningChallenger1.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "1500000000000", coins.AmountOf(cred).String())

// 	winningChallenger2 := votes.falseVotes[2].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, winningChallenger2.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "1500000000000", coins.AmountOf(cred).String())

// 	winningChallenger3 := votes.falseVotes[3].(challenge.Challenge)
// 	coins = k.bankKeeper.GetCoins(ctx, winningChallenger3.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "3000000000000", coins.AmountOf(cred).String())

// 	winningVoter1 := votes.falseVotes[4].(TokenVote)
// 	coins = k.bankKeeper.GetCoins(ctx, winningVoter1.Creator())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "2000000000000", coins.AmountOf(cred).String())

// 	losingBacker1 := votes.trueVotes[0].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, losingBacker1.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingBacker2 := votes.trueVotes[1].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, losingBacker2.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingBacker3 := votes.trueVotes[2].(backing.Backing)
// 	coins = k.bankKeeper.GetCoins(ctx, losingBacker3.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingVoter1 := votes.trueVotes[3].(TokenVote)
// 	coins = k.bankKeeper.GetCoins(ctx, losingVoter1.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())

// 	losingVoter2 := votes.trueVotes[4].(TokenVote)
// 	coins = k.bankKeeper.GetCoins(ctx, losingVoter2.Creator())
// 	assert.Equal(t, "1000000000000", coins.AmountOf(params.StakeDenom).String())
// 	assert.Equal(t, "0", coins.AmountOf(cred).String())
// }

// func TestCheckForEmptyPool(t *testing.T) {
// 	pool, _ := sdk.ParseCoin("4trusteak")
// 	voterCount := int64(10)
// 	err := checkForEmptyPoolConfirmed(pool, voterCount)
// 	assert.Nil(t, err)
// }

// func TestCheckForEmptyPool2(t *testing.T) {
// 	pool, _ := sdk.ParseCoin("5trusteak")
// 	voterCount := int64(10)
// 	err := checkForEmptyPoolConfirmed(pool, voterCount)
// 	assert.Nil(t, err)
// }

// func TestCheckForEmptyPool3(t *testing.T) {
// 	pool, _ := sdk.ParseCoin("9trusteak")
// 	voterCount := int64(10)
// 	err := checkForEmptyPoolConfirmed(pool, voterCount)
// 	assert.Nil(t, err)
// }

// func Test_voterRewardAmount(t *testing.T) {
// 	pool, _ := sdk.ParseCoin("1trusteak")
// 	assert.Equal(t, sdk.NewInt(0), voterRewardAmount(pool, 0))
// }
