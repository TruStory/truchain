package truapi

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	trubank "github.com/TruStory/truchain/x/trubank"
	"github.com/TruStory/truchain/x/voting"
	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/graphql"
	"github.com/TruStory/truchain/x/params"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/users"
	"github.com/TruStory/truchain/x/vote"
	sdk "github.com/cosmos/cosmos-sdk/types"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
)

// TruAPI implements an HTTP server for TruStory functionality using `chttp.API`
type TruAPI struct {
	*chttp.API
	GraphQLClient *graphql.Client
	DBClient      db.Datastore
}

// NewTruAPI returns a `TruAPI` instance populated with the existing app and a new GraphQL client
func NewTruAPI(aa *chttp.App) *TruAPI {
	ta := TruAPI{
		API:           chttp.NewAPI(aa, supported),
		GraphQLClient: graphql.NewGraphQLClient(),
		DBClient:      db.NewDBClient(),
	}

	return &ta
}

// RegisterModels registers types for off-chain DB models
func (ta *TruAPI) RegisterModels() {
	err := ta.DBClient.RegisterModel(&db.TwitterProfile{})
	if err != nil {
		panic(err)
	}
	err = ta.DBClient.RegisterModel(&db.KeyPair{})
	if err != nil {
		panic(err)
	}
}

// WrapHandler wraps a chttp.Handler and returns a standar http.Handler
func WrapHandler(h chttp.Handler) http.Handler {
	return h.HandlerFunc()
}

// RegisterRoutes applies the TruStory API routes to the `chttp.API` router
func (ta *TruAPI) RegisterRoutes() {
	api := ta.Subrouter("/api/v1")
	api.Use(chttp.JSONResponseMiddleware)
	api.Handle("/ping", WrapHandler(ta.HandlePing))
	api.Handle("/graphql", WrapHandler(ta.HandleGraphQL))
	api.Handle("/presigned", WrapHandler(ta.HandlePresigned))
	api.Handle("/unsigned", WrapHandler(ta.HandleUnsigned))
	api.Handle("/register", WrapHandler(ta.HandleRegistration))
	api.Handle("/user", WrapHandler(ta.HandleUserDetails))

	if os.Getenv("MOCK_REGISTRATION") == "true" {
		api.Handle("/mock_register", WrapHandler(ta.HandleMockRegistration))
	}

	ta.RegisterOAuthRoutes()

	// Register routes for Trustory React web app

	appDir := os.Getenv("CHAIN_WEB_DIR")
	if appDir == "" {
		appDir = "build"
	}
	fs := http.FileServer(http.Dir(appDir))

	ta.PathPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if it is not requesting a file with a valid extension serve the index
		if filepath.Ext(path.Base(r.URL.Path)) == "" {
			w.Header().Add("Content-Type", "text/html")
			http.ServeFile(w, r, filepath.Join(appDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))
}

// RegisterOAuthRoutes adds the proper routes needed for the oauth
func (ta *TruAPI) RegisterOAuthRoutes() {
	oauth1Config := &oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_API_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_API_SECRET"),
		CallbackURL:    os.Getenv("CHAIN_OAUTH_CALLBACK"),
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}

	ta.Handle("/auth-twitter", twitter.LoginHandler(oauth1Config, nil))
	ta.Handle("/auth-twitter-callback", twitter.CallbackHandler(oauth1Config, IssueSession(ta), nil))
	ta.Handle("/auth-logout", Logout())
}

// RegisterResolvers builds the app's GraphQL schema from resolvers (declared in `resolver.go`)
func (ta *TruAPI) RegisterResolvers() {
	formatTime := func(t time.Time) string {
		return t.UTC().Format(time.UnixDate)
	}

	getUser := func(ctx context.Context, addr sdk.AccAddress) users.User {
		res := ta.usersResolver(ctx, users.QueryUsersByAddressesParams{Addresses: []string{addr.String()}})
		if len(res) > 0 {
			return res[0]
		}
		return users.User{}
	}

	getBackings := func(ctx context.Context, storyID int64) []backing.Backing {
		return ta.backingsResolver(ctx, app.QueryByIDParams{ID: storyID})
	}

	getChallenges := func(ctx context.Context, storyID int64) []challenge.Challenge {
		return ta.challengesResolver(ctx, app.QueryByIDParams{ID: storyID})
	}

	getVotes := func(ctx context.Context, storyID int64) []vote.TokenVote {
		return ta.votesResolver(ctx, app.QueryByIDParams{ID: storyID})
	}
	getVoteResults := func(ctx context.Context, storyID int64) voting.VoteResult {
		return ta.voteResultsResolver(ctx, app.QueryByIDParams{ID: storyID})
	}

	getTransactions := func(ctx context.Context, creator string) []trubank.Transaction {
		return ta.transactionsResolver(ctx, app.QueryByCreatorParams{Creator: creator})
	}

	getStory := func(ctx context.Context, storyID int64) story.Story {
		return ta.storyResolver(ctx, story.QueryStoryByIDParams{ID: storyID})
	}

	ta.GraphQLClient.RegisterQueryResolver("backing", ta.backingResolver)
	ta.GraphQLClient.RegisterObjectResolver("Backing", backing.Backing{}, map[string]interface{}{
		"amount":    func(ctx context.Context, q backing.Backing) sdk.Coin { return q.Amount() },
		"argument":  func(ctx context.Context, q backing.Backing) string { return q.Argument },
		"weight":    func(ctx context.Context, q backing.Backing) string { return q.Weight().String() },
		"vote":      func(ctx context.Context, q backing.Backing) bool { return q.VoteChoice() },
		"creator":   func(ctx context.Context, q backing.Backing) users.User { return getUser(ctx, q.Creator()) },
		"timestamp": func(ctx context.Context, q backing.Backing) app.Timestamp { return q.Timestamp() },

		// Deprecated: interest is no longer saved in backing
		"interest": func(ctx context.Context, q backing.Backing) sdk.Coin { return sdk.Coin{} },
	})

	ta.GraphQLClient.RegisterQueryResolver("categories", ta.allCategoriesResolver)
	ta.GraphQLClient.RegisterQueryResolver("category", ta.categoryResolver)
	ta.GraphQLClient.RegisterObjectResolver("Category", category.Category{}, map[string]interface{}{
		"id":      func(_ context.Context, q category.Category) int64 { return q.ID },
		"stories": ta.categoryStoriesResolver,
	})

	ta.GraphQLClient.RegisterQueryResolver("challenge", ta.challengeResolver)
	ta.GraphQLClient.RegisterObjectResolver("Challenge", challenge.Challenge{}, map[string]interface{}{
		"amount":    func(ctx context.Context, q challenge.Challenge) sdk.Coin { return q.Amount() },
		"argument":  func(ctx context.Context, q challenge.Challenge) string { return q.Argument },
		"weight":    func(ctx context.Context, q challenge.Challenge) string { return q.Weight().String() },
		"vote":      func(ctx context.Context, q challenge.Challenge) bool { return q.VoteChoice() },
		"creator":   func(ctx context.Context, q challenge.Challenge) users.User { return getUser(ctx, q.Creator()) },
		"timestamp": func(ctx context.Context, q challenge.Challenge) app.Timestamp { return q.Timestamp() },
	})

	ta.GraphQLClient.RegisterObjectResolver("Coin", sdk.Coin{}, map[string]interface{}{
		"amount": func(_ context.Context, q sdk.Coin) string { return q.Amount.String() },
		"denom":  func(_ context.Context, q sdk.Coin) string { return q.Denom },
		"unit":   func(_ context.Context, q sdk.Coin) string { return "preethi" },
	})

	ta.GraphQLClient.RegisterQueryResolver("params", ta.paramsResolver)
	ta.GraphQLClient.RegisterObjectResolver("Params", params.Params{}, map[string]interface{}{
		"amountWeight":      func(_ context.Context, p params.Params) string { return p.StakeParams.AmountWeight.String() },
		"periodWeight":      func(_ context.Context, p params.Params) string { return p.StakeParams.PeriodWeight.String() },
		"minInterestRate":   func(_ context.Context, p params.Params) string { return p.StakeParams.MinInterestRate.String() },
		"maxInterestRate":   func(_ context.Context, p params.Params) string { return p.StakeParams.MaxInterestRate.String() },
		"minArgumentLength": func(_ context.Context, p params.Params) int { return p.StakeParams.MinArgumentLength },
		"maxArgumentLength": func(_ context.Context, p params.Params) int { return p.StakeParams.MaxArgumentLength },

		"storyExpireDuration": func(_ context.Context, p params.Params) string { return p.StoryParams.ExpireDuration.String() },
		"storyMinLength":      func(_ context.Context, p params.Params) int { return p.StoryParams.MinStoryLength },
		"storyMaxLength":      func(_ context.Context, p params.Params) int { return p.StoryParams.MaxStoryLength },
		"storyVotingDuration": func(_ context.Context, p params.Params) string { return p.StoryParams.VotingDuration.String() },

		"challengeMinStake": func(_ context.Context, p params.Params) string { return p.ChallengeParams.MinChallengeStake.String() },
		"challengeMinThreshold": func(_ context.Context, p params.Params) string {
			return p.ChallengeParams.MinChallengeThreshold.String()
		},
		"challengeThresholdPercent": func(_ context.Context, p params.Params) string {
			return p.ChallengeParams.ChallengeToBackingRatio.String()
		},

		"voteStake": func(_ context.Context, p params.Params) string { return p.VoteParams.StakeAmount.String() },

		"stakerRewardRatio": func(_ context.Context, p params.Params) string {
			return p.VotingParams.StakerRewardPoolShare.String()
		},

		"stakeDenom": func(_ context.Context, _ params.Params) string { return app.StakeDenom },

		// Deprecated: replaced by "stakerRewardRatio"
		"challengeRewardRatio": func(_ context.Context, p params.Params) string {
			return p.VotingParams.StakerRewardPoolShare.String()
		},
	})

	ta.GraphQLClient.RegisterQueryResolver("stories", ta.allStoriesResolver)
	ta.GraphQLClient.RegisterQueryResolver("story", ta.storyResolver)
	ta.GraphQLClient.RegisterObjectResolver("Story", story.Story{}, map[string]interface{}{
		"id":                 func(_ context.Context, q story.Story) int64 { return q.ID },
		"backings":           func(ctx context.Context, q story.Story) []backing.Backing { return getBackings(ctx, q.ID) },
		"challenges":         func(ctx context.Context, q story.Story) []challenge.Challenge { return getChallenges(ctx, q.ID) },
		"backingPool":        ta.backingPoolResolver,
		"challengePool":      ta.challengePoolResolver,
		"votingPool":         ta.votingPoolResolver,
		"challengeThreshold": ta.challengeThresholdResolver,
		"category":           ta.storyCategoryResolver,
		"creator":            func(ctx context.Context, q story.Story) users.User { return getUser(ctx, q.Creator) },
		"source":             func(ctx context.Context, q story.Story) string { return q.Source.String() },
		"votes":              func(ctx context.Context, q story.Story) []vote.TokenVote { return getVotes(ctx, q.ID) },
		"voteResults":        func(ctx context.Context, q story.Story) voting.VoteResult { return getVoteResults(ctx, q.ID) },
		"state":              func(ctx context.Context, q story.Story) story.Status { return q.Status },
		"expireTime":         func(_ context.Context, q story.Story) string { return formatTime(q.ExpireTime) },
		"votingStartTime":    func(_ context.Context, q story.Story) string { return formatTime(q.VotingStartTime) },
		"votingEndTime":      func(_ context.Context, q story.Story) string { return formatTime(q.VotingEndTime) },
	})

	ta.GraphQLClient.RegisterObjectResolver("Timestamp", app.Timestamp{}, map[string]interface{}{
		"createdTime": func(_ context.Context, t app.Timestamp) string { return formatTime(t.CreatedTime) },
		"updatedTime": func(_ context.Context, t app.Timestamp) string { return formatTime(t.UpdatedTime) },
	})

	ta.GraphQLClient.RegisterObjectResolver("TwitterProfile", db.TwitterProfile{}, map[string]interface{}{
		"id": func(_ context.Context, q db.TwitterProfile) string { return string(q.ID) },
	})

	ta.GraphQLClient.RegisterQueryResolver("users", ta.usersResolver)
	ta.GraphQLClient.RegisterObjectResolver("User", users.User{}, map[string]interface{}{
		"id":             func(_ context.Context, q users.User) string { return q.Address },
		"coins":          func(_ context.Context, q users.User) sdk.Coins { return q.Coins },
		"pubkey":         func(_ context.Context, q users.User) string { return q.Pubkey.String() },
		"twitterProfile": ta.twitterProfileResolver,
		"transactions": func(ctx context.Context, q users.User) []trubank.Transaction {
			return getTransactions(ctx, q.Address)
		},
	})

	ta.GraphQLClient.RegisterObjectResolver("Transactions", trubank.Transaction{}, map[string]interface{}{
		"id":              func(_ context.Context, q trubank.Transaction) int64 { return q.ID },
		"transactionType": func(_ context.Context, q trubank.Transaction) trubank.TransactionType { return q.TransactionType },
		"amount":          func(_ context.Context, q trubank.Transaction) sdk.Coin { return q.Amount },
		"createdTime":     func(_ context.Context, q trubank.Transaction) time.Time { return q.Timestamp.CreatedTime },
		"story": func(ctx context.Context, q trubank.Transaction) story.Story {
			return getStory(ctx, q.GroupID)
		},
	})

	ta.GraphQLClient.RegisterObjectResolver("URL", url.URL{}, map[string]interface{}{
		"url": func(_ context.Context, q url.URL) string { return q.String() },
	})

	ta.GraphQLClient.RegisterQueryResolver("vote", ta.voteResolver)
	ta.GraphQLClient.RegisterObjectResolver("Vote", vote.TokenVote{}, map[string]interface{}{
		"amount":    func(ctx context.Context, q vote.TokenVote) sdk.Coin { return q.Amount() },
		"argument":  func(ctx context.Context, q vote.TokenVote) string { return q.Argument },
		"vote":      func(ctx context.Context, q vote.TokenVote) bool { return q.VoteChoice() },
		"weight":    func(ctx context.Context, q vote.TokenVote) string { return q.Weight().String() },
		"creator":   func(ctx context.Context, q vote.TokenVote) users.User { return getUser(ctx, q.Creator()) },
		"timestamp": func(ctx context.Context, q vote.TokenVote) app.Timestamp { return q.Timestamp() },
	})

	ta.GraphQLClient.RegisterObjectResolver("voteResults", voting.VoteResult{}, map[string]interface{}{
		"backedCredTotal":     func(_ context.Context, q voting.VoteResult) string { return q.BackedCredTotal.String() },
		"challengedCredTotal": func(_ context.Context, q voting.VoteResult) string { return q.ChallengedCredTotal.String() },
	})

	ta.GraphQLClient.BuildSchema()
}
