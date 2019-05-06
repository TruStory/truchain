package truapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	app "github.com/TruStory/truchain/types"
)

func (ta *TruAPI) sendCommentNotification(n CommentNotificationRequest) {
	if !ta.notificationsInitialized || ta.commentsNotificationsCh == nil {
		return
	}
	ta.commentsNotificationsCh <- n
}

func (ta *TruAPI) runCommentNotificationSender(notifications <-chan CommentNotificationRequest, endpoint string) {
	url := fmt.Sprintf("%s/%s", strings.TrimRight(strings.TrimSpace(endpoint), "/"), "sendCommentNotification")

	for n := range notifications {
		argument := ta.argumentResolver(context.Background(), app.QueryArgumentByID{ID: n.ArgumentID, Raw: true})
		n.StoryID = argument.StoryID
		n.ArgumentCreator = argument.Creator.String()
		b, err := json.Marshal(&n)
		if err != nil {
			fmt.Println("error encoding comment notification request", err)
			continue
		}
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
		if err != nil {
			fmt.Println("error sending comment notification request", err)
			continue
		}
		// only read the status
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusAccepted {
			fmt.Printf("error sending comment notification request status [%s] \n", resp.Status)
			continue
		}
		fmt.Printf("comment notification sent id[%d]\n", n.ID)
	}
}
