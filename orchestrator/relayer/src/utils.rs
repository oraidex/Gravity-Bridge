use clarity::abi::encode_call;
use clarity::{Address, Uint256};
use gravity_utils::num_conversion::downcast_uint256;
use std::env;
use web30::client::Web3;
use web30::jsonrpc::error::Web3Error;
use web30::types::TransactionRequest;

/// This is roughly the maximum number of blocks a reasonable Ethereum node
/// can search in a single request before it starts timing out or behaving badly
pub const BLOCKS_TO_SEARCH: u128 = 5_000u128;

pub fn convert_block_to_search() -> u128 {
    env::var("BLOCK_TO_SEARCH")
        .unwrap_or_else(|_| BLOCKS_TO_SEARCH.to_string())
        .parse::<u128>()
        .unwrap_or_else(|_| BLOCKS_TO_SEARCH)
}

pub async fn get_latest_valset_nonce(
    gravity_contract_address: Address,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let payload = encode_call("state_lastValsetNonce()", &[]).unwrap();
    let result = web3
        .eth_call(TransactionRequest::Legacy {
            from: gravity_contract_address,
            to: gravity_contract_address,
            gas: None,
            gas_price: None,
            value: None,
            data: Some(web30::types::Data(payload)),
            nonce: None,
        })
        .await?;
    let real_num = Uint256::from_be_bytes(&result.0);
    Ok(downcast_uint256(real_num).expect("Valset nonce overflow! Bridge Halt!"))
}

pub async fn get_gravity_id(
    gravity_contract_address: Address,
    web3: &Web3,
) -> Result<String, Web3Error> {
    let payload = encode_call("state_gravityId()", &[]).unwrap();
    let gravity_id_data = web3
        .eth_call(TransactionRequest::Legacy {
            from: gravity_contract_address,
            to: gravity_contract_address,
            gas: None,
            gas_price: None,
            value: None,
            data: Some(web30::types::Data(payload)),
            nonce: None,
        })
        .await
        .unwrap();
    Ok(String::from_utf8(gravity_id_data.0).map_err(|err| {
        Web3Error::ContractCallError(format!(
            "cannot decode gravity id from vec<u8> with error: {}",
            err.to_string()
        ))
    })?)
}

pub async fn get_eth_gravity_checkpoint(
    gravity_contract_address: Address,
    web3: &Web3,
) -> Result<Vec<u8>, Web3Error> {
    let payload = encode_call("state_lastValsetCheckpoint()", &[]).unwrap();
    web3.eth_call(TransactionRequest::Legacy {
        from: gravity_contract_address,
        to: gravity_contract_address,
        gas: None,
        gas_price: None,
        value: None,
        data: Some(web30::types::Data(payload)),
        nonce: None,
    })
    .await
    .map(|data| data.0)
}
