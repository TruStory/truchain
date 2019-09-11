#!/usr/bin/env python3

import lib

def process_genesis(genesis, parsed_args):
    genesis['app_state']['distribution']['community_tax'] = '0.800000000000000000'

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to adjust community tax',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
