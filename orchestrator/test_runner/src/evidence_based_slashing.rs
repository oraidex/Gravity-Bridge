//! This is a test for Evidence based slashing, we simply create a bad signature and submit it as evidence
//! we don't launch the orchestrators here as they are not required.

use crate::utils::ValidatorKeys;
use crate::EVM_CHAIN_PREFIX;
use crate::STAKING_TOKEN;
use crate::STARTING_STAKE_PER_VALIDATOR;
use crate::TOTAL_TIMEOUT;
use crate::{get_fee, utils::get_operator_address};
use clarity::utils::bytes_to_hex_str;
use clarity::Address as EthAddress;
use cosmos_gravity::query::get_latest_valsets;
use cosmos_gravity::{send::submit_bad_signature_evidence, utils::BadSignatureEvidence};
use deep_space::PrivateKey;
use deep_space::{Coin, Contact};
use ethereum_gravity::message_signatures::encode_valset_confirm_hashed;
use ethereum_gravity::{message_signatures::encode_valset_confirm, utils::get_gravity_id};
use gravity_proto::cosmos_sdk_proto::cosmos::staking::v1beta1::QueryValidatorsRequest;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::types::{Valset, ValsetMember};
use tonic::transport::Channel;
use web30::client::Web3;

// takes a vector of valsets and returns the latest nonce or none if there
// are no valsets currently generated
fn latest_valset_nonce(valsets: &[Valset]) -> Option<u64> {
    let mut ret = None;
    for v in valsets {
        if let Some(n) = ret {
            if v.nonce > n {
                ret = Some(v.nonce)
            }
        } else {
            ret = Some(v.nonce)
        }
    }
    ret
}

pub async fn evidence_based_slashing(
    web30: &Web3,
    grpc_client: GravityQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
) {
    let mut grpc_client = grpc_client;

    // first step in this test is to ensure that our slashing vitim does not
    // have 33% and therefore will not halt the chain when jailed. Since our
    // default validator set is 3 validators with 33% each we need to change it some.
    delegate_to_validator(&keys, keys[1].validator_key, contact).await;
    delegate_to_validator(&keys, keys[2].validator_key, contact).await;

    let starting_valset_nonce = latest_valset_nonce(
        &get_latest_valsets(&mut grpc_client, EVM_CHAIN_PREFIX.as_str())
            .await
            .unwrap(),
    );
    info!("Starting valset nonce is {:?}", starting_valset_nonce);

    // our slashing victim is just the first validator
    let cosmos_private_key = keys[0].validator_key;
    let eth_private_key = keys[0].eth_key;
    let eth_addr = eth_private_key.to_address();
    // reporter is another validator using their delegate key
    let submitter_private_key = keys[1].orch_key;
    // this is a false valset, one that happens to contain only the
    // validator signing it, as if they were trying to take over the
    // bridge. This valset isn't valid for submitting but that's not a
    // condition of the slashing
    let false_valset = Valset {
        nonce: 500,
        members: vec![ValsetMember {
            power: 1337,
            eth_address: eth_addr,
        }],
        reward_amount: 0u8.into(),
        reward_token: None,
    };
    let gravity_id = get_gravity_id(gravity_address, eth_addr, web30)
        .await
        .unwrap();
    let message = encode_valset_confirm(gravity_id.clone(), false_valset.clone());
    let checkpoint = encode_valset_confirm_hashed(gravity_id.clone(), false_valset.clone());
    let eth_signature = eth_private_key.sign_ethereum_msg(&message);
    info!(
        "Created signature {} over checkpoint {} with Gravity ID {} using address {}",
        eth_signature,
        bytes_to_hex_str(&checkpoint),
        gravity_id,
        eth_addr
    );

    // now we are prepared to submit our evidence, we check first that validator 0 is in the set
    print_validator_status(contact).await;
    let (is_in_set, jailed) =
        check_validator(contact, cosmos_private_key, "BOND_STATUS_BONDED").await;
    assert!(is_in_set);
    assert!(!jailed);
    info!("Target validator is in the set and not jailed");

    info!("Submitting Evidence");
    // submit the evidence
    let res = submit_bad_signature_evidence(
        EVM_CHAIN_PREFIX.as_str(),
        submitter_private_key,
        get_fee(None),
        contact,
        BadSignatureEvidence::Valset(false_valset),
        eth_signature,
    )
    .await
    .unwrap();
    trace!("{:?}", res);

    // confirm that the validator for which the evidence has been submitted is removed
    let (is_in_set, jailed) =
        check_validator(contact, cosmos_private_key, "BOND_STATUS_UNBONDING").await;
    assert!(is_in_set);
    assert!(jailed);

    // a validator set should have been generated by this validator getting slashed
    let ending_valset_nonce = latest_valset_nonce(
        &get_latest_valsets(&mut grpc_client, EVM_CHAIN_PREFIX.as_str())
            .await
            .unwrap(),
    );
    match (starting_valset_nonce, ending_valset_nonce) {
        (Some(starting), Some(ending)) => {
            if ending <= starting {
                panic!("No new validator set created due to slashing!")
            } else {
                info!("Successfully generated valset due to slashing");
            }
        }
        (None, None) => panic!("No new validator set created due to delegating or slashing!"),
        (Some(_), None) => panic!("Nonsense valset state!"),
        (None, Some(_)) => info!("Successfully generated valset due to slashing"),
    }

    info!("Evidence based slashing test succeeded! Validator now jailed!");
}

async fn check_validator(contact: &Contact, key: impl PrivateKey, filter: &str) -> (bool, bool) {
    let validators = contact
        .get_validators_list(QueryValidatorsRequest {
            pagination: None,
            status: filter.to_string(),
        })
        .await
        .unwrap();
    let addr = get_operator_address(key);
    for val in validators {
        if val.operator_address == addr.to_string() {
            return (true, val.jailed);
        }
    }
    (false, false)
}

async fn print_validator_status(contact: &Contact) {
    let validators = contact.get_active_validators().await.unwrap();
    for val in validators.iter() {
        info!(
            "Validator: {} Power: {} Jailed: {}",
            val.operator_address, val.tokens, val.jailed
        )
    }
}

/// Delegates to a specific validator
async fn delegate_to_validator(keys: &[ValidatorKeys], to: impl PrivateKey, contact: &Contact) {
    let delegate_address = get_operator_address(to);
    let amount = Coin {
        denom: STAKING_TOKEN.to_string(),
        amount: (STARTING_STAKE_PER_VALIDATOR / 4).into(),
    };
    let res = contact
        .delegate_to_validator(
            delegate_address,
            amount,
            get_fee(None),
            keys[1].validator_key,
            Some(TOTAL_TIMEOUT),
        )
        .await
        .unwrap();
    trace!("{:?}", res);
}
