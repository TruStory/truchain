#!/usr/bin/env python3

import lib

def process_genesis(genesis, parsed_args):
    genesis['app_state']['bank']['send_enabled'] = False

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to disable send txs',
        default_chain_id='betanet-1',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
