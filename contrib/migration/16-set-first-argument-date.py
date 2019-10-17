#!/usr/bin/env python3

import lib


def process_genesis(genesis, parsed_args):
    first_argument_mapping = {}
    for a in genesis['app_state']['trustaking']['arguments']:
      if a['claim_id'] in first_argument_mapping:
        if first_argument_mapping[a['claim_id']] > a['created_time']:
          first_argument_mapping[a['claim_id']] = a['created_time']
      else:
        first_argument_mapping[a['claim_id']] = a['created_time']

    for c in genesis['app_state']['claim']['claims']:
      if c['id'] in first_argument_mapping:
        c['first_argument_time'] = first_argument_mapping[c['id']]
    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to upgrade cosmos modules',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
