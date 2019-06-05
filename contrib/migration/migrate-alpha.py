#!/usr/bin/env python3

import lib


def process_genesis(exported_genesis, new_genesis, parsed_args):
    del exported_genesis['consensus_params']
    del exported_genesis['app_state']['auth']
    del exported_genesis['app_state']['expiration']

    # new genesis state
    exported_genesis['consensus_params'] = new_genesis['consensus_params']
    exported_genesis['app_state']['auth'] = new_genesis['app_state']['auth']
    exported_genesis['app_state']['staking'] = new_genesis['app_state']['staking']
    exported_genesis['app_state']['distr'] = new_genesis['app_state']['distr']
    exported_genesis['app_state']['genutil'] = new_genesis['app_state']['genutil']
    exported_genesis['app_state']['staking'] = new_genesis['app_state']['staking']
    exported_genesis['app_state']['params'] = new_genesis['app_state']['params']

    exported_genesis['app_state']['expiration'] = new_genesis['app_state']['expiration']

    # reconfiguration
    exported_genesis['app_state']['staking']['params']['bond_denom'] = parsed_args.bond_denom.strip()
    exported_genesis['chain_id'] = parsed_args.chain_id.strip()
    exported_genesis['genesis_time'] = parsed_args.start_time
    return exported_genesis


if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Convert genesis.json ',
        default_chain_id='truchain',
        default_start_time='2019-06-05T12:00:00Z',
        bond_denom='trusteak',
    )
    lib.main(parser, process_genesis)
