#!/usr/bin/env python3

import copy

import lib


def process_genesis(genesis, parsed_args):
    # rename story -> claim
    genesis['app_state']['claim'] = copy.deepcopy(genesis['app_state']['story'])

    # get total_backed, total_challenged, total_stakers per story id
    totals = total_backed_challenged_stakers_by_story_id(genesis)

    # migrate story state
    migrate_story_data(genesis['app_state']['claim'], totals, genesis['app_state']['category']['categories'])

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def migrate_story_data(story_data, totals, categories):
    del story_data['story_queue']
    story_data['params'] = {
        'min_claim_length': story_data['params']['min_story_length'],
        'max_claim_length': story_data['params']['max_story_length'],
    }
    story_data['claims'] = story_data['stories']
    del story_data['stories']
    for s in story_data['claims']:
        s['id'] = s['id']
        s['community_id'] = get_category_slug(categories, s['category_id'])
        del s['category_id']
        s['body'] = s['body']
        s['creator'] = s['creator']
        if 'source' in s:
            s['source'] = s['source']
        s['created_time'] = s['timestamp']['created_time']
        s['total_backed'] = { 'amount': str(totals[s['id']]['total_backed']), 'denom': 'trusteak' }
        s['total_challenged'] = { 'amount': str(totals[s['id']]['total_challenged']), 'denom': 'trusteak' }
        s['total_stakers'] = str(totals[s['id']]['total_stakers'])
        del s['timestamp']
        del s['status']
        del s['expire_time']
        del s['type']

def get_category_slug(categories, category_id):
    for s in categories:
        if s['id'] == category_id:
            return s['slug']
    raise Exception('Category not found')

def total_backed_challenged_stakers_by_story_id(genesis):
    totals = dict()
    for s in genesis['app_state']['story']['stories']:
        totals[s['id']] = {}
        totals[s['id']]['total_backed'] = 0
        totals[s['id']]['total_challenged'] = 0
        totals[s['id']]['total_stakers'] = 0
    for b in genesis['app_state']['backing']['backings']:
        vote = b['vote']
        totals[vote['story_id']]['total_backed'] = totals[vote['story_id']]['total_backed'] + int(vote['amount']['amount'])
        totals[vote['story_id']]['total_stakers'] = totals[vote['story_id']]['total_stakers'] + 1
    for c in genesis['app_state']['challenge']['challenges']:
        vote = c['vote']
        totals[vote['story_id']]['total_challenged'] = totals[vote['story_id']]['total_challenged'] + int(vote['amount']['amount'])
        totals[vote['story_id']]['total_stakers'] = totals[vote['story_id']]['total_stakers'] + 1
    return totals

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json from stories to claims',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
