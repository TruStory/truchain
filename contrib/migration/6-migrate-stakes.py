#!/usr/bin/env python3

import copy

import lib

# type Stake struct {
# 	ID          generate auto-incrementing ID because backing/challenge ids are not globally unique
# 	ArgumentID  backing/challenge vote->argument_id
# 	Type        backing/challenge, endorse -> upvote
# 	Amount      backing/challenge vote->amount
# 	Creator     backing/challenge vote->creator
# 	CreatedTime backing/challenge timestamp->created_time
# 	EndTime     backing/challenge timestamp->created_time
# 	Expired     true
# }

# In progress stakes will have the EndTime of the story end time

STAKE_BACKING = 0
STAKE_CHALLENGE = 1
STAKE_UPVOTE = 2

PENDING = 0
EXPIRED = 1

stake_id = 0

def process_genesis(genesis, parsed_args):
    genesis['app_state']['trustaking']['stakes'] = []

    migrate_votes_to_stakes(genesis, genesis['app_state']['backing']['backings'])
    migrate_votes_to_stakes(genesis, genesis['app_state']['challenge']['challenges'])

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def migrate_votes_to_stakes(genesis, votes):
    global stake_id
    for v in votes:
      stake_id += 1
      stake = {}
      vote = v['vote']
      stake['id'] = str(stake_id)
      stake['argument_id'] = vote['argument_id']
      argument = get_stake_argument(genesis, vote['argument_id'])
      claim = get_argument_claim(genesis, argument['claim_id'])
      stake['community_id'] = claim['community_id']
      stake['type'] = get_vote_type(genesis, vote)
      stake['amount'] = vote['amount']
      stake['creator'] = vote['creator']
      stake['created_time'] = vote['timestamp']['created_time']
      stake['end_time'] = get_vote_expiration(genesis, vote)
      stake['expired'] = is_vote_expired(genesis, vote)
      genesis['app_state']['trustaking']['stakes'].append(stake)

def get_vote_type(genesis, vote):
  for a in genesis['app_state']['argument']['arguments']:
    if a['creator'] == vote['creator'] and a['id'] == vote['argument_id']:
      return (STAKE_BACKING if vote['vote'] == True else STAKE_CHALLENGE)
  return STAKE_UPVOTE

def get_vote_expiration(genesis, vote):
  for s in genesis['app_state']['story']['stories']:
    if s['id'] == vote['story_id']:
      return s['expire_time']
  raise Exception('Story expire time not found for story id: ' + vote['story_id'])

def is_vote_expired(genesis, vote):
  for s in genesis['app_state']['story']['stories']:
    if s['id'] == vote['story_id']:
      return s['status'] == EXPIRED
  raise Exception('Story status not found for story id: ' + vote['story_id'])

def get_stake_argument(genesis, argument_id):
    for a in genesis['app_state']['trustaking']['arguments']:
        if a['id'] == argument_id:
            return a
    raise Exception('Argument not found with argument id: ' + argument_id)

def get_argument_claim(genesis, claim_id):
    for c in genesis['app_state']['claim']['claims']:
        if c['id'] == claim_id:
            return c
    raise Exception('Claim not found with claim id: ' + claim_id)

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json backings and challenges to stakes',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
