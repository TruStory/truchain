#!/usr/bin/env python3

import copy

import lib


def process_genesis(genesis, parsed_args):
    # rename story -> claim
    genesis['app_state']['claim'] = copy.deepcopy(genesis['app_state']['story'])
    # migrate story state
    migrate_story_data(genesis['app_state']['claim'])

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def migrate_story_data(story_data):
    del story_data['story_queue']
    story_data['params'] = {
        'min_claim_length': story_data['params']['min_story_length'],
        'max_claim_length': story_data['params']['max_story_length'],
    }
    story_data['claims'] = story_data['stories']
    del story_data['stories']
    for s in story_data['claims']:
        s['id'] = s['id']
        s['community_id'] = s['category_id']
        del s['category_id']
        s['body'] = s['body']
        s['creator'] = s['creator']
        if 'source' in s:
            s['source'] = s['source']
        s['created_time'] = s['timestamp']['created_time']
        del s['timestamp']
        del s['status']
        del s['expire_time']
        del s['type']

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json from stories to claims',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
