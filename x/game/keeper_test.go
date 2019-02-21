package game

// var creator = sdk.AccAddress([]byte{1, 2})

// func TestCreateGame(t *testing.T) {
// 	ctx, k, categoryKeeper := mockDB()

// 	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)

// 	gameID, err := k.Create(ctx, storyID, creator)
// 	assert.Nil(t, err)
// 	assert.Equal(t, int64(1), gameID)
// }

// func TestRegisterChallengeNoBackersMeetMinChallenge(t *testing.T) {
// 	ctx, k, categoryKeeper := mockDB()

// 	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
// 	gameID, _ := k.Create(ctx, storyID, creator)

// 	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
// 	err := k.AddToChallengePool(ctx, gameID, amount)
// 	assert.Nil(t, err)

// 	game, _ := k.Game(ctx, gameID)
// 	assert.Equal(t, sdk.NewInt(5000000), game.ChallengePool.Amount)
// }

// func TestRegisterChallengeNoBackersNotMeetMinChallenge(t *testing.T) {
// 	ctx, k, categoryKeeper := mockDB()

// 	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
// 	gameID, _ := k.Create(ctx, storyID, creator)

// 	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
// 	err := k.AddToChallengePool(ctx, gameID, amount)
// 	assert.Nil(t, err)

// 	game, _ := k.Game(ctx, gameID)
// 	assert.Equal(t, sdk.NewInt(5000000), game.ChallengePool.Amount)
// }

// func TestRegisterChallengeHaveBackersMeetThreshold(t *testing.T) {
// 	ctx, k, categoryKeeper := mockDB()

// 	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
// 	gameID, _ := k.Create(ctx, storyID, creator)
// 	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(100))
// 	argument := "cool story brew"

// 	// back story with 100trusteak
// 	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
// 	duration := 30 * 24 * time.Hour
// 	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration)

// 	// challenge with 33trusteak (33% of total backings)
// 	amount, _ = sdk.ParseCoin("33trusteak")
// 	err := k.AddToChallengePool(ctx, gameID, amount)
// 	assert.Nil(t, err)
// }

// func TestRegisterChallengeHaveBackersNotMeetThreshold(t *testing.T) {
// 	ctx, k, categoryKeeper := mockDB()

// 	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
// 	gameID, _ := k.Create(ctx, storyID, creator)
// 	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(100))
// 	argument := "cool story brew"

// 	// back story with 100trusteak
// 	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
// 	duration := 30 * 24 * time.Hour
// 	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration)

// 	// challenge with 32truteak (32% of total backings)
// 	amount = sdk.NewCoin(app.StakeDenom, sdk.NewInt(32))
// 	err := k.AddToChallengePool(ctx, gameID, amount)
// 	assert.Nil(t, err)
// }

// func TestSetGame(t *testing.T) {
// 	ctx, k, _ := mockDB()

// 	game := Game{ID: int64(5)}
// 	k.set(ctx, game)

// 	savedGame, err := k.Game(ctx, int64(5))
// 	assert.Nil(t, err)
// 	assert.Equal(t, game.ID, savedGame.ID)
// }

// func Test_challengeThresholdNoBacking(t *testing.T) {
// 	_, k, _ := mockDB()
// 	amt := k.ChallengeThreshold(sdk.NewCoin(app.StakeDenom, sdk.ZeroInt()))

// 	assert.Equal(t, "10000000000trusteak", amt.String())
// }

// // challenge threshold should not go below min challenge stake
// // if challenge threshold is 1/3 of backing and min challenge stake is 10
// // when a story has backing amount of 21, challenge threshold should be 10, not 7
// // then instead
// func Test_challengeThresholdWithSmallBacking(t *testing.T) {
// 	_, k, _ := mockDB()
// 	amt := k.ChallengeThreshold(sdk.NewCoin(app.StakeDenom, sdk.NewInt(21000000000)))

// 	assert.Equal(t, "21000000000trusteak", amt.String())
// }

// func Test_challengeThresholdWithBacking(t *testing.T) {
// 	_, k, _ := mockDB()
// 	amt := k.ChallengeThreshold(sdk.NewCoin(app.StakeDenom, sdk.NewInt(100000000000)))

// 	assert.Equal(t, "100000000000trusteak", amt.String())
// }

// func Test_start(t *testing.T) {
// 	ctx, k, categoryKeeper := mockDB()

// 	storyID := createFakeStory(ctx, k.storyKeeper, categoryKeeper)
// 	gameID, _ := k.Create(ctx, storyID, creator)
// 	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(100000000000))
// 	argument := "cool story brew"

// 	// back story with 100trusteak
// 	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
// 	duration := 30 * 24 * time.Hour
// 	k.backingKeeper.Create(ctx, storyID, amount, argument, creator, duration)

// 	// challenge with 33trusteak (33% of total backings)
// 	amount = sdk.NewCoin(app.StakeDenom, sdk.NewInt(100000000000))
// 	err := k.AddToChallengePool(ctx, gameID, amount)
// 	assert.Nil(t, err)

// 	// test queue sizes
// 	// assert.Equal(t, uint64(1), k.pendingList(ctx).Len())
// 	// assert.Equal(t, uint64(1), k.queue(ctx).List.Len())
// }
