#!/usr/bin/env python3

import copy
import datetime

import lib

TRANSACTION_REGISTRATION = 0

def process_genesis(genesis, parsed_args):
    transaction_id = 0

    genesis['app_state']['trubank2'] = {
        'params': {
            'reward_broker_address': None,
        }
    }

    genesis['app_state']['trubank2']['transactions'] = []

    # create initial trusteak balance transactions
    for a in genesis['app_state']['accounts']:
        for c in a['coins']:
            if c['denom'] == 'trusteak':
                transaction_id += 1
                transaction = {
                    'id': str(transaction_id),
                    'type': TRANSACTION_REGISTRATION,
                    'app_account_address': a['address'],
                    'reference_id': '0',
                    'amount': c,
                    'created_time': datetime.datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%S.%f')[:-3] + 'Z'
                }
                genesis['app_state']['trubank2']['transactions'].append(transaction)

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to add trubank2',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
