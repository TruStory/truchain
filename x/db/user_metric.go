package db

import (
	"time"
)

// UserMetric is the db model to interact with user metrics
type UserMetric struct {
	tableName                 struct{}  `sql:"user_metrics"`
	Address                   string    `json:"address"`
	AsOnDate                  time.Time `json:"as_on_date"`
	CategoryID                int64     `json:"category_id"`
	TotalClaims               uint64    `json:"total_claims"  sql:"type:,notnull"`
	TotalArguments            uint64    `json:"total_arguments"  sql:"type:,notnull"`
	TotalClaimsBacked         uint64    `json:"total_claims_backed"  sql:"type:,notnull"`
	TotalClaimsChallenged     uint64    `json:"total_claims_challenged"  sql:"type:,notnull"`
	TotalAmountBacked         uint64    `json:"total_amount_backed"  sql:"type:,notnull"`
	TotalAmountChallenged     uint64    `json:"total_amount_challenged"  sql:"type:,notnull"`
	TotalEndorsementsGiven    uint64    `json:"total_endorsements_given"  sql:"type:,notnull"`
	TotalEndorsementsReceived uint64    `json:"total_endorsements_received"  sql:"type:,notnull"`
	StakeEarned               uint64    `json:"stake_earned"  sql:"type:,notnull"`
	StakeLost                 uint64    `json:"stake_lost"  sql:"type:,notnull"`
	StakeBalance              uint64    `json:"stake_balance"  sql:"type:,notnull"`
	InterestEarned            uint64    `json:"interest_earned"  sql:"type:,notnull"`
	TotalAmountAtStake        uint64    `json:"total_amount_at_stake"  sql:"type:,notnull"`
	TotalAmountStaked         uint64    `json:"total_amount_staked"  sql:"type:,notnull"`
	CredEarned                uint64    `json:"cred_earned"  sql:"type:,notnull"`
}

// AggregateStatisticsByAddressBetweenDates gets and aggregates the user metrics for a given address on a given date
func (c *Client) AggregateStatisticsByAddressBetweenDates(address string, from string, to string) ([]UserMetric, error) {
	userMetrics := make([]UserMetric, 0)
	err := c.Model(&userMetrics).
		Column("as_on_date", "category_id").
		ColumnExpr(`
			sum(total_claims) as total_claims,
			sum(total_arguments) as total_arguments,
			sum(total_backed) as total_backed,
			sum(total_challenged) as total_challenged,
			sum(total_endorsements_given) as total_endorsements_given,
			sum(total_endorsements_received) as total_endorsements_received,
			sum(stake_earned) as stake_earned,
			sum(stake_lost) as stake_lost,
			sum(stake_balance) as stake_balance,
			sum(interest_earned) as interest_earned,
			sum(total_amount_at_stake) as total_amount_at_stake,
			sum(total_amount_staked) as total_amount_staked,
			sum(cred_earned) as cred_earned
		`).
		Where("address = ?", address).
		Where("as_on_date >= ?", from).
		Where("as_on_date <= ?", to).
		Group("as_on_date").
		Group("category_id").
		Order("as_on_date").
		Order("category_id").
		Select()
	if err != nil {
		return nil, err
	}

	return userMetrics, nil
}
