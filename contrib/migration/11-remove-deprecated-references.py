#!/usr/bin/env python3

import copy

import lib


def process_genesis(genesis, parsed_args):
    del genesis['app_state']['backing']
    del genesis['app_state']['challenge']
    del genesis['app_state']['expiration']
    del genesis['app_state']['story']
    del genesis['app_state']['trubank']
    
    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    genesis['app_state']['truslashing']['params']['max_detailed_reason_length'] = '140'
    genesis['app_state']['claim']['params']['claim_admins'] = ['cosmos1tfpcnjzkthft3ynewqvn7mtdk7guf3knjdqg4d','cosmos1xqc5gwzpg3fyv5en2fzyx36z2se5ks33tt57e7', 'cosmos1xqc5gwz923znjvzyg3pnxdfsgcu4jv34mep8hp','cosmos1xqc5gwzpgdr4wjz8xscnys2jx3f9x4zy223g9w','cosmos1xqc5gwzpgdr4gk3nfdxn24jegc6rv5zewn82ch','cosmos1xqc5gwzy2ge9ysec2vursk2etqm5yjzceu04ez','cosmos1xqc5gwzpgdr4jkjkx3z9xs2x2gurgs22fzksza','cosmos1xqc5gwz9g4pnvvznfpzyxkp5t92yx5pkx9lnsh','cosmos1xqc5gwzgx9z9xvjjxq6rzkfs23f9j5jx748tp4', 'cosmos1xqc5gwzpfp8ygkzdfdpnq4j3xd8y6djy5z8gfn', 'cosmos1xqc5gwfhx4f5k3jcf4r55j6ntge5y3jtxesy8r']
    genesis['app_state']['truslashing']['params']['slash_admins'] = ['cosmos1tfpcnjzkthft3ynewqvn7mtdk7guf3knjdqg4d','cosmos1xqc5gwzpg3fyv5en2fzyx36z2se5ks33tt57e7', 'cosmos1xqc5gwz923znjvzyg3pnxdfsgcu4jv34mep8hp','cosmos1xqc5gwzpgdr4wjz8xscnys2jx3f9x4zy223g9w','cosmos1xqc5gwzpgdr4gk3nfdxn24jegc6rv5zewn82ch','cosmos1xqc5gwzy2ge9ysec2vursk2etqm5yjzceu04ez','cosmos1xqc5gwzpgdr4jkjkx3z9xs2x2gurgs22fzksza','cosmos1xqc5gwz9g4pnvvznfpzyxkp5t92yx5pkx9lnsh','cosmos1xqc5gwzgx9z9xvjjxq6rzkfs23f9j5jx748tp4', 'cosmos1xqc5gwzpfp8ygkzdfdpnq4j3xd8y6djy5z8gfn', 'cosmos1xqc5gwfhx4f5k3jcf4r55j6ntge5y3jtxesy8r']
    genesis['app_state']['community']['params']['community_admins'] = ['cosmos1tfpcnjzkthft3ynewqvn7mtdk7guf3knjdqg4d','cosmos1xqc5gwzpg3fyv5en2fzyx36z2se5ks33tt57e7', 'cosmos1xqc5gwz923znjvzyg3pnxdfsgcu4jv34mep8hp','cosmos1xqc5gwzpgdr4wjz8xscnys2jx3f9x4zy223g9w','cosmos1xqc5gwzpgdr4gk3nfdxn24jegc6rv5zewn82ch','cosmos1xqc5gwzy2ge9ysec2vursk2etqm5yjzceu04ez','cosmos1xqc5gwzpgdr4jkjkx3z9xs2x2gurgs22fzksza','cosmos1xqc5gwz9g4pnvvznfpzyxkp5t92yx5pkx9lnsh','cosmos1xqc5gwzgx9z9xvjjxq6rzkfs23f9j5jx748tp4', 'cosmos1xqc5gwzpfp8ygkzdfdpnq4j3xd8y6djy5z8gfn', 'cosmos1xqc5gwfhx4f5k3jcf4r55j6ntge5y3jtxesy8r']
    genesis['app_state']['trustaking']['params']['staking_admins'] = ['cosmos1tfpcnjzkthft3ynewqvn7mtdk7guf3knjdqg4d','cosmos1xqc5gwzpg3fyv5en2fzyx36z2se5ks33tt57e7', 'cosmos1xqc5gwz923znjvzyg3pnxdfsgcu4jv34mep8hp','cosmos1xqc5gwzpgdr4wjz8xscnys2jx3f9x4zy223g9w','cosmos1xqc5gwzpgdr4gk3nfdxn24jegc6rv5zewn82ch','cosmos1xqc5gwzy2ge9ysec2vursk2etqm5yjzceu04ez','cosmos1xqc5gwzpgdr4jkjkx3z9xs2x2gurgs22fzksza','cosmos1xqc5gwz9g4pnvvznfpzyxkp5t92yx5pkx9lnsh','cosmos1xqc5gwzgx9z9xvjjxq6rzkfs23f9j5jx748tp4', 'cosmos1xqc5gwzpfp8ygkzdfdpnq4j3xd8y6djy5z8gfn', 'cosmos1xqc5gwfhx4f5k3jcf4r55j6ntge5y3jtxesy8r']

    genesis['app_state']['truslashing']['params']['slash_min_stake']['amount'] = '25000000000'

    genesis['app_state']['account']['app_accounts'] = []

    for a in genesis['app_state']['accounts']:
        appAccount = {
            'addresses': [a['address']]
        }
        genesis['app_state']['account']['app_accounts'].append(appAccount)
    
    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json from deprecated dont',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
