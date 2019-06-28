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
        'min_id_length': '3',
        'max_id_length': '15',
        'min_name_length': '5',
        'max_name_length': '25',
        'max_description_length': '140',
    }
    category_data['communities'] = category_data['categories']
    del category_data['categories']
    for s in category_data['communities']:
        s['name'] = s['title']
        del s['title']
        if 'description' in s:
          s['description'] = s['description']
        if 'timestamp' in s:
          s['created_time'] = s['timestamp']['created_time']
          del s['timestamp']
        del s['total_cred']
        if s['slug'] == 'crypto':
          s['description'] = 'Satoshi inspired a new generation of technologists to rethink how the world’s financial and economic system works. But the toxicity on crypto Twitter leaves little room for constructive debate. Until now.'
        if s['slug'] == 'product':
          s['description'] = 'TruStory is a startup. We make critical product decisions every week and want you to be a part of the process. Let’s debate how to make TruStory better together.'
        if s['slug'] == 'tech':
          s['description'] = 'Technology has revolutionized our lives in every facet. What are the pros & cons of the latest innovations and applications that shape our future?'
        if s['slug'] == 'entertainment':
          s['description'] = 'Content is being created and consumed at a breakneck pace, whether it’s via videos, podcasts, books or music. Debate the value of specific content or the medium itself.'
        if s['slug'] == 'sports':
          s['description'] = 'Sports inspire, motivate, and capture the attention of young kids and senior citizens alike. The estimated size of the global sports industry is over a trillion dollars. Ready, set, go!'
        s['id'] = s['slug']
        del s['slug']

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json from categories to communities',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
