#!/usr/bin/env python3

import copy

import lib


def process_genesis(genesis, parsed_args):
    # rename category -> community 
    genesis['app_state']['community'] = copy.deepcopy(genesis['app_state']['category'])
    # migrate category state
    migrate_category_data(genesis['app_state']['community'])

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def migrate_category_data(category_data):
    category_data['params'] = {
        'min_name_length': '5',
        'max_name_length': '25',
        'min_slug_length': '3',
        'max_slug_length': '15',
        'max_description_length': '140',
    }
    category_data['communities'] = category_data['categories']
    del category_data['categories']
    for s in category_data['communities']:
        s['id'] = s['id']
        s['name'] = s['title']
        del s['title']
        s['slug'] = s['slug']
        if 'description' in s:
          s['description'] = s['description']
        if 'timestamp' in s:
          s['created_time'] = s['timestamp']['created_time']
          del s['timestamp']
        s['total_earned_stake'] = s['total_cred']
        del s['total_cred']
        s['total_earned_stake']['denom'] = 'trusteak'

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json from categories to communities',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
