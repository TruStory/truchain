#!/usr/bin/env python3

import copy
import argparse
import sys
import re


def process_genesis(genesis, parsed_args):
    # create a capture group around the address part of the cosmos address
    # replace cosmos prefix with tru, and add back captured group
    genesis = re.sub(r'cosmos(.{39}.)', r'tru\1', genesis)

    return genesis


def init_default_argument_parser(prog_desc):
    parser = argparse.ArgumentParser(description=prog_desc)
    parser.add_argument(
        '--exported-genesis',
        help='exported genesis.json file',
        type=argparse.FileType('r'),
        required=True,
    )
    return parser


def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    genesis = args.exported_genesis.read()

    print(process_genesis_func(genesis=genesis, parsed_args=args))

if __name__ == '__main__':
    parser = init_default_argument_parser(
        prog_desc='Migrate genesis.json trusteak to tru',
    )
    main(parser, process_genesis)
