package truapi

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/dghubble/go-twitter/twitter"
)

// HandleMockRegistration takes an empty request and returns a `RegistrationResponse`
func (ta *TruAPI) HandleMockRegistration(r *http.Request) chttp.Response {
	// Get the mock Twitter User from the auth token
	twitterUser := getMockTwitterUser()

	return RegisterTwitterUser(ta, twitterUser)
}

func getMockTwitterUser() *twitter.User {
	// getting a random id
	id := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(999999999)
	mocked := &twitter.User{
		ID:              id,
		IDStr:           strconv.FormatInt(int64(id), 10),
		ScreenName:      "trustory_engineering",
		Name:            "Trustory Engineering",
		Email:           "engineering@trustory.io",
		ProfileImageURL: "https://pbs.twimg.com/profile_images/999336936572567552/SY65rL1h_bigger.jpg",
	}

	return mocked
}
