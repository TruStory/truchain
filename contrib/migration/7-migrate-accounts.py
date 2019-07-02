#!/usr/bin/env python3

import copy

import lib


def process_genesis(genesis, parsed_args):
    genesis['app_state']['account'] = {
        'params': {
            'registrar': None,
            'max_slash_count': '3',
        }
    }

    genesis['app_state']['trustaking']['users_earnings'] = []

    for a in genesis['app_state']['accounts']:
        old_coins = copy.deepcopy(a['coins'])
        a['coins'] = []
        earnings = {}
        earnings['address'] = a['address']
        earnings['coins'] = []
        for c in old_coins:
            if c['denom'] != 'trusteak':
                # move all cred coins to users_earnings
                earnings['coins'].append(c)
            else:
                # keep only trusteak in account coins
                a['coins'].append(c)
        genesis['app_state']['trustaking']['users_earnings'].append(earnings)

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to add account module',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
