package cookies

import (
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/TruStory/truchain/x/db"
	"github.com/gorilla/securecookie"
)

// GetUserFromCookie gets the user context from the request's cookie
func GetUserFromCookie(r *http.Request) (map[string]string, error) {
	truUser, err := r.Cookie("tru-user")
	if err != nil {
		return nil, err
	}

	hashKey, err := hex.DecodeString(os.Getenv("COOKIE_HASH_KEY"))
	if err != nil {
		return nil, err
	}
	blockKey, err := hex.DecodeString(os.Getenv("COOKIE_ENCRYPT_KEY"))
	if err != nil {
		return nil, err
	}
	var s = securecookie.New(hashKey, blockKey)

	decodedTruUser := make(map[string]string)
	err = s.Decode("tru-user", truUser.Value, &decodedTruUser)
	if err != nil {
		return nil, err
	}

	// if the cookie is stale
	cookieTime, err := strconv.ParseInt(decodedTruUser["created_at"], 10, 64)
	if err != nil {
		return nil, err
	}
	if time.Unix(cookieTime, 0).Before(time.Now().Add(-2 * time.Hour)) {
		return nil, errors.New("Stale cookie found")
	}

	return decodedTruUser, nil
}

// SetUserToCookie takes a user and encodes it into a cookie value.
func SetUserToCookie(twitterProfile *db.TwitterProfile) (string, error) {
	// Saves and excrypts the context in the cookie
	hashKey, err := hex.DecodeString(os.Getenv("COOKIE_HASH_KEY"))
	if err != nil {
		return "", err
	}
	blockKey, err := hex.DecodeString(os.Getenv("COOKIE_ENCRYPT_KEY"))
	if err != nil {
		return "", err
	}
	s := securecookie.New(hashKey, blockKey)
	cookieValue := map[string]string{
		"twitter-profile-id": strconv.FormatInt(twitterProfile.ID, 10),
		"address":            twitterProfile.Address,
		"created_at":         strconv.FormatInt(time.Now().Unix(), 10),
	}
	encodedValue, err := s.Encode("tru-user", cookieValue)
	if err != nil {
		return "", err
	}

	return encodedValue, nil
}
