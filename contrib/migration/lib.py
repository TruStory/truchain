#!/usr/bin/env python3

import argparse
import json
import sys


def init_default_argument_parser(prog_desc, default_chain_id, default_start_time, bond_denom):
    parser = argparse.ArgumentParser(description=prog_desc)
    parser.add_argument(
        '--exported-genesis',
        help='exported genesis.json file',
        type=argparse.FileType('r'),
        required=True,
    )
    parser.add_argument(
        '--new-genesis',
        help='new genesis.json file format',
        type=argparse.FileType('r'),
        required=True,
    )
    parser.add_argument('--chain-id', type=str, default=default_chain_id)
    parser.add_argument('--start-time', type=str, default=default_start_time)
    parser.add_argument('--bond-denom', type=str, default=bond_denom)
    return parser


def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    if args.chain_id.strip() == '':
        sys.exit('chain-id required')
    exported_genesis = json.loads(args.exported_genesis.read())
    new_genesis = json.loads(args.new_genesis.read())

    print(json.dumps(process_genesis_func(
        exported_genesis=exported_genesis, new_genesis=new_genesis, parsed_args=args,), indent=True, sort_keys=True))
