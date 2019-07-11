#!/usr/bin/env python3

import copy
import datetime

import lib

TRANSACTION_REGISTRATION = 0
TRANSACTION_UPVOTE_RECEIVED = 8

STAKE_BACKING = 0
STAKE_CHALLENGE = 1
STAKE_UPVOTE = 2

def process_genesis(genesis, parsed_args):
    transaction_id = 0

    genesis['app_state']['trubank2'] = {
        'params': {
            'reward_broker_address': None,
        }
    }

    genesis['app_state']['trubank2']['transactions'] = []

    earned_coins = initialize_earned_coins(genesis)
    stake_refund = initialize_stake_refund(genesis)

    # add agree received transaction for each migrated endorsment at .25 conversion
    for s1 in genesis['app_state']['trustaking']['stakes']:
        if s1['expired'] == False:
            s1['expired'] = True
            s1['end_time'] = datetime.datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%S.%f')[:-3] + 'Z'
            stake_refund[s1['creator']] += int(s1['amount']['amount'])

    # create initial trusteak balance transactions
    for a1 in genesis['app_state']['accounts']:
        for c1 in a1['coins']:
            if c1['denom'] == 'trusteak':
                transaction_id += 1
                transaction = {
                    'id': str(transaction_id),
                    'type': TRANSACTION_REGISTRATION,
                    'app_account_address': a1['address'],
                    'reference_id': '0',
                    'amount': {
                        'amount': str(int(c1['amount']) + stake_refund[a1['address']]),
                        'denom': c1['denom']
                    },
                    'created_time': datetime.datetime.strptime('Apr 1 2019', '%b %d %Y').strftime('%Y-%m-%dT%H:%M:%S.%f')[:-3] + 'Z'
                }
                genesis['app_state']['trubank2']['transactions'].append(transaction)


    # add agree received transaction for each migrated endorsment at .25 conversion
    for s2 in genesis['app_state']['trustaking']['stakes']:
        if s2['type'] == STAKE_UPVOTE:
            argument = get_stake_argument(genesis, s2['argument_id'])
            claim = get_argument_claim(genesis, argument['claim_id'])
            earned_amount = int(float(s2['amount']['amount']) * 0.025)
            earned_coins[argument['creator']][claim['community_id']] += earned_amount
            transaction_id += 1
            transaction = {
                'id': str(transaction_id),
                'type': TRANSACTION_UPVOTE_RECEIVED,
                'app_account_address': argument['creator'],
                'reference_id': s2['id'],
                'amount': {
                    'amount': str(earned_amount),
                    'denom': 'trusteak'
                },
                'community_id': claim['community_id'],
                'created_time': s2['end_time']
            }
            genesis['app_state']['trubank2']['transactions'].append(transaction)

    # convert cred to earned coins at .25 conversion rate
    genesis['app_state']['trustaking']['users_earnings'] = []
    for a2 in genesis['app_state']['accounts']:
        for c2 in a2['coins']:
            if c2['denom'] == 'trusteak':
                c2['amount'] = str(int(c2['amount']) + earned_coin_balance(earned_coins[a2['address']]) + stake_refund[a2['address']])
        earnings = {}
        earnings['address'] = a2['address']
        earnings['coins'] = []
        for community_id, amount in earned_coins[a2['address']].items():
            if amount != 0:
                earnings['coins'].append({
                    'amount': str(amount),
                    'denom': community_id
                })
                check_cred(genesis, a2['address'], amount, community_id)
        genesis['app_state']['trustaking']['users_earnings'].append(earnings)

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

# initialize user's earned coins to 0 for each community
def initialize_earned_coins(genesis):
    earned_coins = {}
    for a in genesis['app_state']['accounts']:
        earned_coins[a['address']] = {}
        for c in genesis['app_state']['community']['communities']:
            earned_coins[a['address']][c['id']] = 0
    return earned_coins

def check_cred(genesis, address, amount, denom):
    for a in genesis['app_state']['accounts']:
        if a['address'] == address:
            for c in a['coins']:
                if c['denom'] == denom:
                    if int(float(c['amount']) * .25) - amount > 250000000:
                        print('Cred conversion failed!', c['amount'], c['denom'], a['address'], int(float(c['amount']) * .25), amount, denom)

def initialize_stake_refund(genesis):
    stake_refund = {}
    for a in genesis['app_state']['accounts']:
        stake_refund[a['address']] = 0
    return stake_refund

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

def earned_coin_balance(earned_coins):
    balance = 0
    for community_id, amount in earned_coins.items():
        balance += amount
    return balance

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to add trubank2',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
