#!/usr/bin/env python3

import lib

def process_genesis(genesis, parsed_args):
    genesis['app_state']['distribution']['community_tax'] = '0.800000000000000000'

    genesis['app_state']['account']['params']['user_growth_allocation'] = '0.200000000000000000'
    genesis['app_state']['account']['params']['stakeholder_allocation'] = '0.200000000000000000'

    genesis['app_state']['trustaking']['params']['user_reward_allocation'] = '0.200000000000000000'

    genesis['app_state']['supply'] = {
        'supply': [
            {'denom': 'tru', 'amount': '253644628704484'},
        ],
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

    # Add modules accounts that hold inflation fund pools

    # This is the account that holds the total inflation from each block.
    # The naming is bit of a misnomer due to historical reasons with Cosmos.
    # The address of the account is derived from the name "fee_collector", which
    # the mint modules uses. Thus it cannot be changed.
    feeCollectorAccnt = {
        'address': 'cosmos17xpfvakm2amg962yls6f84z3kell8c5lserqta',
        'coins': [{'denom': 'tru', 'amount': '0'}],
        'sequence_number': '0',
        'account_number':'0',
        'original_vesting': [],
        'delegated_free': [],
        'delegated_vesting': [],
        'start_time': '0',
        'end_time': '0',
        'module_name': 'fee_collector',
        'module_permissions': [],
    }
    genesis['app_state']['accounts'].append(feeCollectorAccnt)

    userRewardPoolAccnt = {
        'address': 'cosmos1ed82m7snyk8mux8xxpwygvtyq633a4k43rfp8l',
        'coins': [{'denom': 'tru', 'amount': '2500000'}],
        'sequence_number': '0',
        'account_number':'0',
        'original_vesting': [],
        'delegated_free': [],
        'delegated_vesting': [],
        'start_time': '0',
        'end_time': '0',
        'module_name': 'user_reward_tokens_pool',
        'module_permissions': [],
    }
    genesis['app_state']['accounts'].append(userRewardPoolAccnt)

    userGrowthPoolAccnt = {
        'address': 'cosmos1f7x5wx3adh6klcurmd8n36etx4elgu9d4wkys3',
        'coins': [{'denom': 'tru', 'amount': '2500000'}],
        'sequence_number': '0',
        'account_number':'0',
        'original_vesting': [],
        'delegated_free': [],
        'delegated_vesting': [],
        'start_time': '0',
        'end_time': '0',
        'module_name': 'user_growth_tokens_pool',
        'module_permissions': [],
    }
    genesis['app_state']['accounts'].append(userGrowthPoolAccnt)

    stakeholderPoolAccnt = {
        'address': 'cosmos1m9rhdryf059x684um3pa9n30tsdxuww84pxemz',
        'coins': [{'denom': 'tru', 'amount': '0'}],
        'sequence_number': '0',
        'account_number':'0',
        'original_vesting': [],
        'delegated_free': [],
        'delegated_vesting': [],
        'start_time': '0',
        'end_time': '0',
        'module_name': 'stakeholder_tokens_pool',
        'module_permissions': [],
    }
    genesis['app_state']['accounts'].append(stakeholderPoolAccnt)

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to adjust community tax',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
