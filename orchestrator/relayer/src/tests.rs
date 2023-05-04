use std::{str::FromStr, time::Duration};

use actix::System;
use ethereum_gravity::message_signatures::encode_valset_confirm;
use futures::try_join;
use sha3::{Digest, Keccak256};
use web30::{client::Web3, EthAddress};

use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;

use crate::utils::{get_eth_gravity_checkpoint, get_gravity_id, get_latest_valset_nonce};
use std::env;

#[test]
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
fn test_compare_checkpoint_hash_bsc_network() {
    let runner = System::new();
    let web3 = Web3::new("https://bsc-dataseed.binance.org", Duration::from_secs(30));
    let gravity_contract_addr =
        EthAddress::from_str("0xb40C364e70bbD98E8aaab707A41a52A2eAF5733f").unwrap();
    let evm_chain_prefix = "oraib";
    let oraibridge_grpc_url = env::var("ORAIBRIDGE_GRPC").unwrap();

    runner.block_on(async move {
        // getting latest valset nonce
        let nonce_future = get_latest_valset_nonce(gravity_contract_addr, &web3);
        let eth_checkpoint_hash_future = get_eth_gravity_checkpoint(gravity_contract_addr, &web3);
        let gravity_id_future = get_gravity_id(gravity_contract_addr, &web3);
        let mut grpc_client = GravityQueryClient::connect(oraibridge_grpc_url)
            .await
            .unwrap();

        if let Ok((nonce, eth_checkpoint_hash, gravity_id)) =
            try_join!(nonce_future, eth_checkpoint_hash_future, gravity_id_future)
        {
            let result =
                cosmos_gravity::query::get_valset(&mut grpc_client, evm_chain_prefix, nonce)
                    .await
                    .unwrap();
            assert_eq!(result.is_some(), true);
            let valset = result.unwrap();
            let checkpoint = encode_valset_confirm(gravity_id, valset);
            let cosmos_checkpoint_hash = Keccak256::digest(&checkpoint);
            assert_eq!(eth_checkpoint_hash, cosmos_checkpoint_hash.to_vec())
        }
    });
}

#[test]
fn test_compare_checkpoint_hash_tron_network() {
    let runner = System::new();
    let web3 = Web3::new("https://api.trongrid.io/jsonrpc", Duration::from_secs(30));
    let gravity_contract_addr =
        EthAddress::from_str("0x73Ddc880916021EFC4754Cb42B53db6EAB1f9D64").unwrap();
    let evm_chain_prefix = "trontrx-mainnet";
    let oraibridge_grpc_url = env::var("ORAIBRIDGE_GRPC").unwrap();

    runner.block_on(async move {
        // getting latest valset nonce
        let nonce_future = get_latest_valset_nonce(gravity_contract_addr, &web3);
        let eth_checkpoint_hash_future = get_eth_gravity_checkpoint(gravity_contract_addr, &web3);
        let gravity_id_future = get_gravity_id(gravity_contract_addr, &web3);
        let mut grpc_client = GravityQueryClient::connect(oraibridge_grpc_url)
            .await
            .unwrap();

        if let Ok((nonce, eth_checkpoint_hash, gravity_id)) =
            try_join!(nonce_future, eth_checkpoint_hash_future, gravity_id_future)
        {
            let result =
                cosmos_gravity::query::get_valset(&mut grpc_client, evm_chain_prefix, nonce)
                    .await
                    .unwrap();
            assert_eq!(result.is_some(), true);
            let valset = result.unwrap();
            let checkpoint = encode_valset_confirm(gravity_id, valset);
            let cosmos_checkpoint_hash = Keccak256::digest(&checkpoint);
            assert_eq!(eth_checkpoint_hash, cosmos_checkpoint_hash.to_vec())
        }
    });
}

// how to run the test: ORAIBRIDGE_GRPC=http://localhost:9090 cargo test --package relayer
