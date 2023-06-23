use std::{str::FromStr, time::Duration};

use actix::System;
use web30::{client::Web3, EthAddress};

use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;

use crate::{
    find_latest_valset::find_latest_valset,
    utils::{get_gravity_id, get_latest_valset_nonce},
};
use std::env;

#[test]
#[ignore]
fn test_get_latest_valset_nonce() {
    let runner = System::new();
    let web3 = Web3::new("https://api.trongrid.io/jsonrpc", Duration::from_secs(30));
    let gravity_contract_addr =
        EthAddress::from_str("0x73Ddc880916021EFC4754Cb42B53db6EAB1f9D64").unwrap();

    runner.block_on(async move {
        // test getting latest valset nonce
        let result = get_latest_valset_nonce(gravity_contract_addr, &web3)
            .await
            .unwrap();
        assert_ne!(result, 0);
        println!("result: {}", result)
    });
}

#[test]
#[ignore]
fn test_get_gravity_id() {
    let runner = System::new();
    let web3 = Web3::new("https://api.trongrid.io/jsonrpc", Duration::from_secs(30));
    let gravity_contract_addr =
        EthAddress::from_str("0x73Ddc880916021EFC4754Cb42B53db6EAB1f9D64").unwrap();

    runner.block_on(async move {
        // test getting latest valset nonce
        let gravity_id = get_gravity_id(gravity_contract_addr, &web3).await.unwrap();
        assert_ne!(gravity_id, "oraibridge-32".to_string());
    });
}

#[test]
#[ignore]
fn test_compare_valset_hash_bsc_network() {
    let runner = System::new();
    let web3 = Web3::new("https://bsc-dataseed.binance.org", Duration::from_secs(30));
    let gravity_contract_address =
        EthAddress::from_str("0xb40C364e70bbD98E8aaab707A41a52A2eAF5733f").unwrap();
    let evm_chain_prefix = "oraib";
    let oraibridge_grpc_url = env::var("ORAIBRIDGE_GRPC").unwrap();

    runner.block_on(async move {
        let mut grpc_client = GravityQueryClient::connect(oraibridge_grpc_url)
            .await
            .unwrap();
        let nonce = get_latest_valset_nonce(gravity_contract_address, &web3)
            .await
            .unwrap();
        let cosmos_valset =
            cosmos_gravity::query::get_valset(&mut grpc_client, evm_chain_prefix, nonce)
                .await
                .unwrap()
                .unwrap();
        let valset = find_latest_valset(
            &mut grpc_client,
            evm_chain_prefix,
            gravity_contract_address,
            &web3,
        )
        .await
        .unwrap();
        assert_eq!(valset, cosmos_valset);
    });
}

#[test]
#[ignore]
fn test_compare_checkpoint_hash_tron_network() {
    let runner = System::new();
    let web3 = Web3::new("https://api.trongrid.io/jsonrpc", Duration::from_secs(30));
    let gravity_contract_address =
        EthAddress::from_str("0x73Ddc880916021EFC4754Cb42B53db6EAB1f9D64").unwrap();
    let evm_chain_prefix = "trontrx-mainnet";
    let oraibridge_grpc_url = env::var("ORAIBRIDGE_GRPC").unwrap();

    runner.block_on(async move {
        let mut grpc_client = GravityQueryClient::connect(oraibridge_grpc_url)
            .await
            .unwrap();
        let nonce = get_latest_valset_nonce(gravity_contract_address, &web3)
            .await
            .unwrap();
        let cosmos_valset =
            cosmos_gravity::query::get_valset(&mut grpc_client, evm_chain_prefix, nonce)
                .await
                .unwrap()
                .unwrap();
        let valset = find_latest_valset(
            &mut grpc_client,
            evm_chain_prefix,
            gravity_contract_address,
            &web3,
        )
        .await
        .unwrap();
        assert_eq!(valset, cosmos_valset);
    });
}

// how to run the test: ORAIBRIDGE_GRPC=http://localhost:9090 cargo test --package relayer
