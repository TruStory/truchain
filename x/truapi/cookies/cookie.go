package cookies

import (
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/TruStory/truchain/x/db"
	"github.com/gorilla/securecookie"
)

const (
	// UserCookieName contains the name of the cookie that stores the user
	UserCookieName string = "tru-user"

	// AuthenticationExpiry is the period for which,
	// the logged in user must be considered authenticated
	AuthenticationExpiry int64 = 72 // in hours
)

// AuthenticatedUser denotes the data structure of the data inside the encrypted cookie
type AuthenticatedUser struct {
	TwitterProfileID int64
	Address          string
	AuthenticatedAt  int64
}

// GetLoginCookie returns the http cookie that authenticates and identifies the given user
func GetLoginCookie(twitterProfile *db.TwitterProfile) (*http.Cookie, error) {
	value, err := MakeLoginCookieValue(twitterProfile)
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:     UserCookieName,
		HttpOnly: true,
		Value:    value,
		Expires:  time.Now().Add(time.Duration(AuthenticationExpiry) * time.Hour),
		Domain:   os.Getenv("COOKIE_HOST"),
	}

	return &cookie, nil
}

// GetLogoutCookie returns the http cookie that overrides
// the login cookie to practically delete it.
func GetLogoutCookie() *http.Cookie {
	cookie := http.Cookie{
		Name:     UserCookieName,
		HttpOnly: true,
		Value:    "",
		Expires:  time.Now(),
		Domain:   os.Getenv("COOKIE_HOST"),
		MaxAge:   0,
	}

	return &cookie
}

// GetAuthenticatedUser gets the user from the request's http cookie
func GetAuthenticatedUser(r *http.Request) (*AuthenticatedUser, error) {
	cookie, err := r.Cookie(UserCookieName)
	if err != nil {
		return nil, err
	}

	s, err := getSecureCookieInstance()
	if err != nil {
		return nil, err
	}

	user := &AuthenticatedUser{}
	err = s.Decode(UserCookieName, cookie.Value, &user)
	if err != nil {
		return nil, err
	}

	if isStale(user) {
		return nil, errors.New("Stale cookie found")
	}

	return user, nil
}

// MakeLoginCookieValue takes a user and encodes it into a cookie value.
func MakeLoginCookieValue(twitterProfile *db.TwitterProfile) (string, error) {
	s, err := getSecureCookieInstance()
	if err != nil {
		return "", err
	}

	cookieValue := &AuthenticatedUser{
		TwitterProfileID: twitterProfile.ID,
		Address:          twitterProfile.Address,
		AuthenticatedAt:  time.Now().Unix(),
	}
	encodedValue, err := s.Encode(UserCookieName, cookieValue)
	if err != nil {
		return "", err
	}

	return encodedValue, nil
}

// isStale returns whether the cookie older than what is accepted
func isStale(user *AuthenticatedUser) bool {
	return time.
		// if the authentication time...
		Unix(user.AuthenticatedAt, 0).
		// ...exists before in past...
		Before(
			// ...than the valid period.
			time.Now().Add(time.Duration(-1*AuthenticationExpiry) * time.Hour))
}

func getSecureCookieInstance() (*securecookie.SecureCookie, error) {
	// Saves and excrypts the context in the cookie
	hashKey, err := hex.DecodeString(os.Getenv("COOKIE_HASH_KEY"))
	if err != nil {
		return nil, err
	}
	blockKey, err := hex.DecodeString(os.Getenv("COOKIE_ENCRYPT_KEY"))
	if err != nil {
		return nil, err
	}
	return securecookie.New(hashKey, blockKey), nil
}
