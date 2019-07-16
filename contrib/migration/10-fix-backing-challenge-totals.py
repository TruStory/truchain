#!/usr/bin/env python3

import copy

import lib

STAKE_BACKING = 0
STAKE_CHALLENGE = 1
STAKE_UPVOTE = 2

def process_genesis(genesis, parsed_args):
    # update argument length to 1500
    genesis['app_state']['trustaking']['params']['argument_body_max_length'] = '1500'

    # update reward broker address
    genesis['app_state']['trubank2']['params']['reward_broker_address'] = 'cosmos1tfpcnjzkthft3ynewqvn7mtdk7guf3knjdqg4d'

    # get total_backed, total_challenged per story id
    totals = total_backed_challenged_by_claim_id(genesis)

    # overwrite total backing/total challenge on a claim
    for c in genesis['app_state']['claim']['claims']:
      if totals[c['id']]['total_backed'] != c['total_backed']['amount']:
        if int(totals[c['id']]['total_backed']) < int(c['total_backed']['amount']):
          raise Exception('Total backed decreased', totals[c['id']]['total_backed'], c['total_backed']['amount'])
        c['total_backed'] = totals[c['id']]['total_backed']
      if totals[c['id']]['total_challenged'] != c['total_challenged']['amount']:
        if int(totals[c['id']]['total_challenged']) < int(c['total_challenged']['amount']):
          raise Exception('Total challenged decreased', totals[c['id']]['total_challenged'], c['total_challenged']['amount'])
        c['total_challenged'] = totals[c['id']]['total_challenged']

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def total_backed_challenged_by_claim_id(genesis):
    totals = dict()
    for c in genesis['app_state']['claim']['claims']:
        totals[c['id']] = {}
        totals[c['id']]['total_backed'] = 0
        totals[c['id']]['total_challenged'] = 0

    for s in genesis['app_state']['trustaking']['stakes']:
        argument = get_stake_argument(genesis, s['argument_id'])
        claim_id = argument['claim_id']

        if s['type'] == STAKE_BACKING:
          totals[claim_id]['total_backed'] = totals[claim_id]['total_backed'] + int(s['amount']['amount'])
        elif s['type'] == STAKE_CHALLENGE:
          totals[claim_id]['total_challenged'] = totals[claim_id]['total_challenged'] + int(s['amount']['amount'])
        elif s['type'] == STAKE_UPVOTE:
          if argument['stake_type'] == STAKE_BACKING:
            totals[claim_id]['total_backed'] = totals[claim_id]['total_backed'] + int(s['amount']['amount'])
          elif argument['stake_type'] == STAKE_CHALLENGE:
            totals[claim_id]['total_challenged'] = totals[claim_id]['total_challenged'] + int(s['amount']['amount'])

    for c in genesis['app_state']['claim']['claims']:
        totals[c['id']]['total_backed'] = str(totals[c['id']]['total_backed'])
        totals[c['id']]['total_challenged'] = str(totals[c['id']]['total_challenged'])

    return totals

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
        prog_desc='Migrate genesis.json from stories to claims',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
