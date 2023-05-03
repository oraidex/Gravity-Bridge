
use crate::args::Erc721ToCosmosOpts;
use crate::utils::TIMEOUT;
use ethereum_gravity::send_erc721_to_cosmos::send_erc721_to_cosmos;
use ethereum_gravity::utils::get_valset_nonce;
use gravity_utils::{
    connection_prep::{check_for_eth, create_rpc_connections},
};

pub async fn erc721_to_cosmos(args: Erc721ToCosmosOpts, prefix: String) {
    let gravity_address = args.gravity_contract_address;
    let gravityerc721_address = args.gravityerc721_contract_address;
    let erc721_address = args.token_contract_address;
    let cosmos_dest = args.destination;
    let ethereum_key = args.ethereum_key;
    let ethereum_public_key = ethereum_key.to_address();
    let ethereum_rpc = args.ethereum_rpc;
    let token_id = args.token_id.clone();

    let connections = create_rpc_connections(prefix, None, Some(ethereum_rpc), TIMEOUT).await;

    let web3 = connections.web3.unwrap();

    get_valset_nonce(gravity_address, ethereum_public_key, &web3)
        .await
        .expect("Incorrect Gravity Address or otherwise unable to contact Gravity");

    check_for_eth(ethereum_public_key, &web3).await;
    
    info!(
        "Sending {} / {} to Cosmos from {} to {}",
        token_id.clone(), erc721_address, ethereum_public_key, cosmos_dest
    );
    // we send erc721 token to the gravityerc721 contract to register a deposit
    let res = send_erc721_to_cosmos(
        erc721_address,
        gravityerc721_address,
        token_id,
        cosmos_dest,
        ethereum_key,
        Some(TIMEOUT),
        &web3,
        vec![],
    )
    .await;
    match res {
        Ok(tx_id) => info!("Send to Cosmos txid: {:#066x}", tx_id),
        Err(e) => info!("Failed to send token! {:?}", e),
    }

    info!(
        "Your token should show up in the account {} on the destination chain within 5 minutes",
        cosmos_dest
    )
}
