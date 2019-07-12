#!/usr/bin/env python3

import copy

import lib

def process_genesis(genesis, parsed_args):
    genesis['app_state']['truslashing'] = {
        'params': {
            'min_slash_count': '5',
            'slash_magnitude': '3',
            'slash_min_stake': '10000000000trusteak',
            'slash_admins': ['cosmos1xqc5gwzpg3fyv5en2fzyx36z2se5ks33tt57e7', 'cosmos1xqc5gwz923znjvzyg3pnxdfsgcu4jv34mep8hp'],
            'curator_share': '0.250000000000000000',
        }
    }

    genesis['app_state']['account']['params']['jail_duration'] = '604800000000000'

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to add slashing module',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
