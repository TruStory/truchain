#!/usr/bin/env python3

import lib

def process_genesis(genesis, parsed_args):
    # reset community pool
    genesis['app_state']['distribution']['community_tax'] = '0.500000000000000000'
    genesis['app_state']['distribution']['fee_pool']['community_pool'] = []

    genesis['app_state']['trudistribution'] = {
        'params': {
            'user_growth_allocation': '0.500000000000000000',
            'user_reward_allocation': '0.500000000000000000',
        }
    }

    # "mint": {
    #             "minter": {
    #                 "annual_provisions": "0.000000000000000000",
    #                 "inflation": "0.200000000000000000"
    #             },
    #             "params": {
    #                 "blocks_per_year": "6311520",
    #                 "goal_bonded": "0.670000000000000000",
    #                 "inflation_max": "0.250000000000000000",
    #                 "inflation_min": "0.100000000000000000",
    #                 "inflation_rate_change": "0.150000000000000000",
    #                 "mint_denom": "tru"
    #             }
    #         },
    genesis['app_state']['mint']['minter']['inflation'] = '0.700000000000000000'
    genesis['app_state']['mint']['params']['inflation_min'] = '0.700000000000000000'
    genesis['app_state']['mint']['params']['inflation_max'] = '0.700000000000000000'

    # remove validators
    del genesis['validators']

    # change bonded_tokens_pool coins to []
    # because the old staking module had delegated shares which added to the bond pool
    # remove coins from registrar which fucks up supply
    for acc in genesis['app_state']['accounts']:
        if acc['module_name'] == 'bonded_tokens_pool':
            acc['coins'] = []
        if acc['address'] == 'cosmos1tfpcnjzkthft3ynewqvn7mtdk7guf3knjdqg4d':
            acc['coins'] = [{'denom': 'tru', 'amount': '1000000000000'}]
        if acc['address'] == 'cosmos1pmp80ys5kplk0gnvmhtxq086xlerkwvcdhk8gx':
            acc['coins'] = []
        if acc['address'] == 'cosmos1em44grl9ylmmnwawwp5fjn079kesatwp67rxjx':
            acc['coins'] = []

    # staking from init genesis
    genesis['app_state']['staking'] = {
        'params': {
            'unbonding_time': '1814400000000000',
            'max_validators': 100,
            'max_entries': 7,
            'bond_denom': 'tru',
        },
        'last_total_power': '0',
        'last_validator_power': None,
        'validators': None,
        'delegations': None,
        'unbonding_delegations': None,
        'redelegations': None,
        'exported': False,
    }

    # supply is set automatically on chain init, so leave empty
    # i.e:
    # genesis['app_state']['supply'] = {
    #     'supply': [
    #         {'denom': 'tru', 'amount': '1332881859320829'},
    #     ],
    # }
    genesis['app_state']['supply'] = {
        'supply': [],
    }

    genesis['app_state']['gov'] = {
        'starting_proposal_id': '1',
        'deposits': None,
        'votes': None,
        'proposals': None,
        'deposit_params': {
            'min_deposit': [{
                'denom': 'tru',
                'amount': '1000',
            }],
            'max_deposit_period': '172800000000000',
        },
        'voting_params': {
            'voting_period': '172800000000000',
        },
        'tally_params': {
            'quorum': '0.334000000000000000',
            'threshold': '0.500000000000000000',
            'veto': '0.334000000000000000',
        },
    }

    genesis['app_state']['crisis'] = {
        'constant_fee': {
            'denom': 'tru',
            'amount': '1000',
        },
    }

    genesis['app_state']['slashing'] = {
        'params': {
            'max_evidence_age': '120000000000',
            'signed_blocks_window': '100',
            'min_signed_per_window': '0.500000000000000000',
            'downtime_jail_duration': '600000000000',
            'slash_fraction_double_sign': '0.050000000000000000',
            'slash_fraction_downtime': '0.010000000000000000',
        },
        'signing_infos': {},
        'missed_blocks': {},
    }

    genesis['app_state']['genutil'] = {
        'gentx': None,
    }

    # module accounts themselves are automatically created...

    genesis['chain_id'] = 'betanet-1'

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to add inflation',
        default_chain_id='betanet-1',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
