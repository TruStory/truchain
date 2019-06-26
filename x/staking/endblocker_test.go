package staking

import (
	"testing"
	"time"
)

func TestEndBlocker(t *testing.T) {
	ctx, _, _, _ := mockDB()
	initialTime := mustParseTime("2006-01-02", "2019-01-01")
	ctx.WithBlockTime(initialTime)

}

func mustParseTime(layout, date string) time.Time {
	t, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}
	return t
}
