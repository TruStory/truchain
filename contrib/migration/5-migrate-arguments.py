#!/usr/bin/env python3

import copy

import lib

# Old argument
# type Argument struct {
# 	ID      int64 `json:"id"`
# 	StoryID int64 `json:"story_id"`
# 	// association with backing or challenge
# 	StakeID int64 `json:"stake_id"`
# 	Body      string         `json:"body"`
# 	Creator   sdk.AccAddress `json:"creator"`
# 	Timestamp app.Timestamp  `json:"timestamp"`
# }

# New argument
# type Argument struct {
# 	ID             id
# 	Creator        creator
# 	ClaimID        story_id
# 	Summary        first 140 chars of body
# 	Body           body 
# 	StakeType      stakeid -> vote
# 	UpvotedCount   count(backing/challenge with argumentId == id && backing/challenge.creator != argument.creator)
# 	UpvotedStake   backing/challenge with argumentId == id && backing/challenge.creator != argument.creator
# 	TotalStake     backing/challenge with argumentId == id
# 	UnhelpfulCount 0
# 	IsUnhelpful    false
# 	CreatedTime    timestamp.created_time
# 	UpdatedTime    timestamp.updated_time
# }

STAKE_BACKING = 0
STAKE_CHALLENGE = 1
STAKE_UPVOTE = 2


def process_genesis(genesis, parsed_args):
    # rename argument.arguments -> truchain/staking.arguments
    genesis['app_state']['trustaking'] = {}
    genesis['app_state']['trustaking']['params'] = {}

    genesis['app_state']['trustaking']['params']['period'] = '604800000000000'
    genesis['app_state']['trustaking']['params']['argument_creation_stake'] = { 'denom': 'trusteak', 'amount': '50000000000' }
    genesis['app_state']['trustaking']['params']['argument_body_max_length'] = "1250"
    genesis['app_state']['trustaking']['params']['argument_body_min_length'] = "25"
    genesis['app_state']['trustaking']['params']['argument_summary_max_length'] = "140"
    genesis['app_state']['trustaking']['params']['argument_summary_min_length'] = "25"
    genesis['app_state']['trustaking']['params']['upvote_stake'] = { 'denom': 'trusteak', 'amount': '10000000000' }
    genesis['app_state']['trustaking']['params']['creator_share'] = "0.500000000000000000"
    genesis['app_state']['trustaking']['params']['interest_rate'] = "1.050000000000000000"
    genesis['app_state']['trustaking']['params']['stake_limit_percent'] = "0.667000000000000000"
    genesis['app_state']['trustaking']['params']['stake_limit_days'] = "604800000000000"
    genesis['app_state']['trustaking']['params']['unjail_upvotes'] = "1"
    genesis['app_state']['trustaking']['params']['max_arguments_per_claim'] = "5"

    genesis['app_state']['trustaking']['arguments'] = copy.deepcopy(genesis['app_state']['argument']['arguments'])

    # get upvote_count, upvote_stake, total_stake, stake_type by argument id
    totals = total_stakes_by_argument_id(genesis)

    # migrate argument state
    migrate_argument_data(genesis, genesis['app_state']['trustaking']['arguments'], totals)

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def migrate_argument_data(genesis, argument_data, totals):
    for a in argument_data:
        a['id'] = a['id']
        a['creator'] = a['creator']
        a['claim_id'] = a['story_id']
        claim = get_argument_claim(genesis, a['claim_id'])
        a['community_id'] = claim['community_id'] 
        del a['story_id']
        a['summary'] = a['body'][0:140]
        a['body'] = a['body']
        a['stake_type'] = totals[a['id']]['stake_type']
        a['upvoted_count'] = str(totals[a['id']]['upvoted_count'])
        a['upvoted_stake'] = { 'amount': str(totals[a['id']]['upvoted_stake']), 'denom': 'trusteak' }
        a['total_stake'] = { 'amount': str(totals[a['id']]['total_stake']), 'denom': 'trusteak' }
        del a['stake_id']
        a['unhelpful_count'] = "0"
        a['is_unhelpful'] = False
        a['created_time'] = a['timestamp']['created_time']
        a['updated_time'] = a['timestamp']['updated_time']
        del a['timestamp']

def get_argument_author_stake(genesis, argument_id, creator):
    for b in genesis['app_state']['backing']['backings']:
      vote = b['vote']
      if vote['argument_id'] == argument_id and vote['creator'] == creator:
        return vote
    for c in genesis['app_state']['challenge']['challenges']:
      vote = c['vote']
      if vote['argument_id'] == argument_id and vote['creator'] == creator:
        return vote
    raise Exception('backing/challenge not found')


def total_stakes_by_argument_id(genesis):
    totals = dict()
    for s in genesis['app_state']['argument']['arguments']:
        totals[s['id']] = {}
        totals[s['id']]['total_count'] = 0
        totals[s['id']]['total_stake'] = 0
    for b in genesis['app_state']['backing']['backings']:
        vote = b['vote']
        totals[vote['argument_id']]['total_count'] = totals[vote['argument_id']]['total_count'] + 1
        totals[vote['argument_id']]['total_stake'] = totals[vote['argument_id']]['total_stake'] + int(vote['amount']['amount'])
    for c in genesis['app_state']['challenge']['challenges']:
        vote = c['vote']
        totals[vote['argument_id']]['total_count'] = totals[vote['argument_id']]['total_count'] + 1
        totals[vote['argument_id']]['total_stake'] = totals[vote['argument_id']]['total_stake'] + int(vote['amount']['amount'])
    # subtract the argument authors stake from total_stake and upvoted_count
    for s in genesis['app_state']['argument']['arguments']:
        author_stake = get_argument_author_stake(genesis, s['id'], s['creator'])
        totals[s['id']]['stake_type'] = STAKE_BACKING if author_stake['vote'] == True else STAKE_CHALLENGE
        totals[s['id']]['upvoted_count'] = totals[s['id']]['total_count'] - 1
        totals[s['id']]['upvoted_stake'] = totals[s['id']]['total_stake'] - int(author_stake['amount']['amount'])
        totals[s['id']]['total_stake'] = totals[s['id']]['total_stake']
        del totals[s['id']]['total_count']
        if totals[s['id']]['upvoted_stake'] > totals[s['id']]['total_stake']:
            raise Exception('Upvoted stake is bigger than total stake for argument id: ' + s['id'])

    return totals

def get_argument_claim(genesis, claim_id):
    for c in genesis['app_state']['claim']['claims']:
        if c['id'] == claim_id:
            return c
    raise Exception('Claim not found with claim id: ' + claim_id)

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json from old arguments to new arguments',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
