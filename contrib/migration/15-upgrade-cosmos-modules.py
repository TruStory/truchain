#!/usr/bin/env python3

import lib

def process_genesis(genesis, parsed_args):
    del genesis['app_state']['auth']['collected_fees']

    v = genesis['app_state']['distr']
    del genesis['app_state']['distr']
    genesis['app_state']['distribution'] = v

    del genesis['app_state']['staking']['pool']

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to upgrade cosmos modules',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
