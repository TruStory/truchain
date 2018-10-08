package story

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Story errors reserve 700 ~ 799.
const (
	DefaultCodespace sdk.CodespaceType = 7

	CodeInvalidStoryBody   sdk.CodeType = 701
	CodeInvalidStoryID     sdk.CodeType = 702
	CodeStoryNotFound      sdk.CodeType = 703
	CodeInvalidEvidenceURL sdk.CodeType = 704
	CodeInvalidStoryType   sdk.CodeType = 706
)

// ErrInvalidStoryBody throws an error on invalid title
func ErrInvalidStoryBody(body string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidStoryBody, "Invalid story body: "+body)
}

// ErrInvalidStoryID throws an error on invalid storyID
func ErrInvalidStoryID(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidStoryID, "Invalid story id: "+fmt.Sprintf("%d", id))
}

// ErrInvalidStoryKind throws an error on invalid story type
func ErrInvalidStoryKind(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidStoryType, "Invalid story kind: "+msg)

}

// ErrInvalidEvidenceURL throws an error when a URL in invalid
func ErrInvalidEvidenceURL(url string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidEvidenceURL, "Invalid evidence URL: "+url)
}

// ErrStoryNotFound throws an error when the searched story is not found
func ErrStoryNotFound(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeStoryNotFound, "Story with id "+fmt.Sprintf("%d", id)+" not found")
}
