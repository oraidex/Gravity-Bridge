use crate::utils::*;
use crate::MINER_ADDRESS;
use crate::MINER_PRIVATE_KEY;
use crate::OPERATION_TIMEOUT;
use crate::TOTAL_TIMEOUT;
use bytes::BytesMut;
use clarity::{Address as EthAddress, Uint256};
use deep_space::address::Address as CosmosAddress;
use deep_space::Contact;
use ethereum_gravity::send_erc721_to_cosmos::send_erc721_to_cosmos;
use ethereum_gravity::utils::get_gravity_sol_address;
use ethereum_gravity::utils::get_valset_nonce;
use cosmos_gravity::query::get_erc721_attestations;
use gravity_proto::nft::QueryNfTsRequest;
use gravity_proto::nft::query_client::QueryClient as NftQueryClient;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_proto::gravity::MsgSendErc721ToCosmosClaim;
use tonic::transport::Channel;
use gravity_utils::error::GravityError;
use prost::Message;
use std::any::type_name;
use std::time::Duration;
use std::time::Instant;
use tokio::time::sleep as delay_for;
use web30::client::Web3;
use web30::types::SendTxOption;

pub async fn erc721_happy_path_test(
    web30: &Web3,
    grpc_client: GravityQueryClient<Channel>,
    cosmos_node_grpc: String,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
    gravityerc721_address: EthAddress,
    erc721_address: EthAddress,
    validator_out: bool,
) {

    let mut grpc_client = grpc_client;

    let grav_sol_address_in_erc721 =
        get_gravity_sol_address(gravityerc721_address, *MINER_ADDRESS, web30)
            .await
            .unwrap();

    info!("ERC721 address is {}", erc721_address);
    info!("Miner address is {}", *MINER_ADDRESS);
    info!("Miner pk is {}", *MINER_PRIVATE_KEY);
    info!(
        "grav_sol_address_in_erc721 is {}",
        grav_sol_address_in_erc721
    );
    info!("Gravity address is {}", gravity_address);
    info!("GravityERC721 address is {}", gravityerc721_address);

    assert_eq!(grav_sol_address_in_erc721, gravity_address);

    let no_relay_market_config = create_default_test_config();
    start_orchestrators(
        keys.clone(),
        gravity_address,
        gravityerc721_address,
        validator_out,
        no_relay_market_config,
    )
    .await;

    // generate an address for coin sending tests, this ensures test imdepotency
    let user_keys = get_user_key(None);

    info!("testing erc721 deposit");
    // Run test with three ERC721 tokens, 200, 201, 202
    for i in 200_i32..203_i32 {
        test_erc721_deposit_panic(
            web30,
            &mut grpc_client,
            cosmos_node_grpc.clone(),
            contact,
            user_keys.cosmos_address,
            gravity_address,
            gravityerc721_address,
            erc721_address,
            Uint256::from_bytes_be(&i.to_be_bytes()),
            None,
        )
        .await;
    }

    info!("testing ERC721 approval utility");
    let token_id_for_approval = Uint256::from_bytes_be(&203_i32.to_be_bytes());
    test_erc721_transfer_utils(
        web30,
        gravityerc721_address,
        erc721_address,
        token_id_for_approval.clone(),
    )
    .await;
}

// Tests an ERC721 deposit and panics on failure
#[allow(clippy::too_many_arguments)]
pub async fn test_erc721_deposit_panic(
    web30: &Web3,
    grpc_client: &mut GravityQueryClient<Channel>,
    cosmos_node_grpc: String,
    contact: &Contact,
    dest: CosmosAddress,
    gravity_address: EthAddress,
    gravity_erc721_address: EthAddress,
    erc721_address: EthAddress,
    token_id: Uint256,
    timeout: Option<Duration>,
) {
    match test_erc721_deposit_result(
        web30,
        grpc_client,
        cosmos_node_grpc,
        contact,
        dest,
        gravity_address,
        gravity_erc721_address,
        erc721_address,
        token_id.clone(),
        timeout,
    )
    .await
    {
        Ok(_) => {
            info!("Successfully bridged ERC721 for token id {}!", token_id)
        }
        Err(_) => {
            panic!("Failed to bridge ERC721 for token id {}!", token_id)
        }
    }
}

/// this function tests Ethereum -> Cosmos deposits of ERC721 tokens
#[allow(clippy::too_many_arguments)]
pub async fn test_erc721_deposit_result(
    web30: &Web3,
    grpc_client: &mut GravityQueryClient<Channel>,
    cosmos_node_grpc: String,
    contact: &Contact,
    dest: CosmosAddress,
    gravity_address: EthAddress,
    gravityerc721_address: EthAddress,
    erc721_address: EthAddress,
    token_id: Uint256,
    timeout: Option<Duration>,
) -> Result<(), GravityError> {
    get_valset_nonce(gravity_address, *MINER_ADDRESS, web30)
        .await
        .expect("Incorrect Gravity Address or otherwise unable to contact Gravity");

    let val = web30
        .get_erc721_symbol(erc721_address, *MINER_ADDRESS)
        .await
        .expect("Not a valid ERC721 contract address");
    info!("In erc721_happy_path_test symbol is {}", val);

    info!(
        "Sending to Cosmos from miner adddress {} to {} with token id {}",
        *MINER_ADDRESS,
        dest,
        token_id.clone()
    );
    // we send an ERC721 to gravityERC721.sol to register a deposit
    let tx_id = send_erc721_to_cosmos(
        erc721_address,
        gravityerc721_address,
        token_id.clone(),
        dest,
        *MINER_PRIVATE_KEY,
        None,
        web30,
        vec![],
    )
    .await
    .expect("Failed to send tokens to Cosmos");

    let _tx_res = web30
        .wait_for_transaction(tx_id, OPERATION_TIMEOUT, None)
        .await
        .expect("Send to cosmos transaction failed to be included into ethereum side");

    let mut grpc_client = grpc_client.clone();
    check_send_erc721_to_cosmos_attestation(&mut grpc_client, erc721_address, dest, *MINER_ADDRESS, token_id.clone()).await?;

    let start = Instant::now();
    let duration = match timeout {
        Some(w) => w,
        None => TOTAL_TIMEOUT,
    };


    let mut grpc_nft_client = NftQueryClient::connect(cosmos_node_grpc)
        .await
        .unwrap();

    while Instant::now() - start < duration {
        // in this while loop wait for owner to change OR wait for event to fire
        info!("Trying to get owner of token_class {} and token_id {} from Cosmos", format!("{}{}", "gravityerc721", erc721_address), token_id.clone());
        let res = grpc_nft_client
        .nf_ts(QueryNfTsRequest {
            class_id: format!("{}{}", "gravityerc721", erc721_address),
            owner: dest.to_string(),
            pagination: None,
        })
        .await;

        if res.is_err() {
            error!("Failed to get nfts of token_class {} from Cosmos. Retrying...", format!("{}{}", "gravityerc721", erc721_address));
            contact.wait_for_next_block(TOTAL_TIMEOUT).await.unwrap();
            continue;
        }

        let nfts = res.unwrap().into_inner().nfts;
        if nfts.len() == 0 {
            error!("No NFTs found for token_class {} from Cosmos. Retrying...", format!("{}{}", "gravityerc721", erc721_address));
            contact.wait_for_next_block(TOTAL_TIMEOUT).await.unwrap();
            continue;
        }
        let mut nft_found = false;
        for nft in nfts {
            if nft.id == token_id.to_string() {
                nft_found = true;
                break;
            }
        }

        if nft_found {
            info!(
                "Successfully moved token_id {} to {}",
                token_id.clone(),
                dest.to_string()
            );
            return Ok(());
        }
        contact.wait_for_next_block(TOTAL_TIMEOUT).await.unwrap();
    }
    Err(GravityError::InvalidBridgeStateError(
        "Did not complete ERC721 deposit!".to_string(),
    ))
}

async fn check_send_erc721_to_cosmos_attestation(
    grpc_client: &mut GravityQueryClient<Channel>,
    erc721_address: EthAddress,
    receiver: CosmosAddress,
    sender: EthAddress,
    token_id: Uint256,
) -> Result<(), GravityError> {
    let start = Instant::now();
    let mut found = false;
    loop {
        iterate_attestations(grpc_client, &mut |decoded: MsgSendErc721ToCosmosClaim| {
            let right_contract = decoded.token_contract == erc721_address.to_string();
            let right_destination = decoded.cosmos_receiver == receiver.to_string();
            let right_sender = decoded.ethereum_sender == sender.to_string();
            let right_token_id = decoded.token_id == token_id.to_string();
            found = right_contract && right_destination && right_sender && right_token_id;
        })
        .await;
        if found {
            break;
        } else if Instant::now() - start > TOTAL_TIMEOUT {
            return Err(GravityError::InvalidBridgeStateError(
                "Could not find the send_erc721_to_cosmos attestation we were looking for!".to_string(),
            ));
        }
        info!("Looking for send_erc721_to_cosmos attestations");
        delay_for(Duration::from_secs(10)).await;
    }
    info!("Found the expected MsgSendERC721ToCosmosClaim attestation with erc721 contract address {} and token id {}", erc721_address, token_id);
    Ok(())
}

pub async fn iterate_attestations<F: FnMut(T), T: Message + Default>(
    grpc_client: &mut GravityQueryClient<Channel>,
    f: &mut F,
) {
    let attestations = get_erc721_attestations(grpc_client, None)
        .await
        .expect("Something happened while getting attestations after delegating to validator");
    for (i, att) in attestations.into_iter().enumerate() {
        let claim = att.clone().claim;
        trace!("Processing attestation {}", i);
        if claim.is_none() {
            trace!("Attestation returned with no claim: {:?}", att);
            continue;
        }
        let claim = claim.unwrap();
        let mut buf = BytesMut::with_capacity(claim.value.len());
        buf.extend_from_slice(&claim.value);

        let decoded = T::decode(buf);

        if decoded.is_err() {
            debug!(
                "Found an attestation which is not a {}: {:?}",
                type_name::<T>(),
                att,
            );
            continue;
        }
        let decoded = decoded.unwrap();
        f(decoded);
    }
}

pub async fn test_erc721_transfer_utils(
    web30: &Web3,
    gravity_erc721_address: EthAddress,
    erc721_address: EthAddress,
    token_id: Uint256,
) {
    let mut options = vec![];
    let nonce = web30
        .eth_get_transaction_count(*MINER_ADDRESS)
        .await
        .expect("Error retrieving nonce from eth");
    options.push(SendTxOption::Nonce(nonce.clone()));

    match web30
        .approve_erc721_transfers(
            erc721_address,
            *MINER_PRIVATE_KEY,
            gravity_erc721_address,
            token_id.clone(),
            None,
            options,
        )
        .await
    {
        Ok(_) => {
            info!("Successfully called ERC721 transfer util for {}!", token_id)
        }
        Err(_) => {
            panic!("Failed ERC721 transfer util for {}!", token_id)
        }
    }
}
