use clarity::Address;
use ethereum_gravity::message_signatures::encode_valset_confirm;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::{error::GravityError, types::Valset};
use tokio::try_join;
use tonic::transport::Channel;
use web30::client::Web3;

use crate::utils::get_eth_gravity_checkpoint;
use crate::utils::get_gravity_id;
use crate::utils::get_latest_valset_nonce;
use sha3::{Digest, Keccak256};

/// This function finds the latest valset on the Gravity contract by querying the latest valset nonce on Ethereum
/// Then use that nonce to query the valset on cosmos
/// It then hashes the Cosmos's valset and compare it with the latest checkpoint stored on the Gravity contract address
/// If they match then we use the valset on Cosmos
/// If not then we panic for safety reasons.
pub async fn find_latest_valset(
    grpc_client: &mut GravityQueryClient<Channel>,
    evm_chain_prefix: &str,
    gravity_contract_address: Address,
    web3: &Web3,
) -> Result<Valset, GravityError> {
    let nonce_future = get_latest_valset_nonce(gravity_contract_address, &web3);
    let eth_checkpoint_future = get_eth_gravity_checkpoint(gravity_contract_address, &web3);
    let gravity_id_future = get_gravity_id(gravity_contract_address, &web3);
    let join_result = try_join!(nonce_future, eth_checkpoint_future, gravity_id_future);
    if let Ok((nonce, eth_checkpoint, gravity_id)) = join_result {
        let result = cosmos_gravity::query::get_valset(grpc_client, evm_chain_prefix, nonce)
            .await
            .unwrap();
        assert_eq!(result.is_some(), true);
        let valset = result.unwrap();
        let checkpoint = encode_valset_confirm(gravity_id, valset.clone());
        let cosmos_checkpoint = Keccak256::digest(&checkpoint);
        /*
        This function exists to provide a warning if Cosmos and Ethereum have different validator sets
        for a given nonce. In the mundane version of this warning the validator sets disagree on sorting order
        which can happen if some relayer uses an unstable sort, or in a case of a mild griefing attack.
        The Gravity contract validates signatures in order of highest to lowest power. That way it can exit
        the loop early once a vote has enough power, if a relayer were to submit things in the reverse order
        they could grief users of the contract into paying more in gas.
        The other (and far worse) way a disagreement here could occur is if validators are colluding to steal
        funds from the Gravity contract and have submitted a highjacking update. If slashing for off Cosmos chain
        Ethereum signatures is implemented you would put that handler here.
        */
        if eth_checkpoint.ne(&cosmos_checkpoint.to_vec()) {
            panic!("Validator sets for nonce {} of Cosmos and nonce {} of Ethereum differ. Possible bridge highjacking!", nonce, valset.nonce);
        }
        return Ok(valset);
    }
    panic!("Could not find the last validator set for contract {} with error: {:?}, probably not a valid Gravity contract!", gravity_contract_address, join_result.unwrap_err().to_string())
}
