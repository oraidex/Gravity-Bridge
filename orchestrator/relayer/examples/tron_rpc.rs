use heliosphere::{core::Address, signer::keypair::Keypair, signer::signer::Signer, RpcClient};
use std::time::Duration;
use web30::client::Web3;

#[actix_rt::main]
pub async fn main() {
    let api = "https://api.trongrid.io";
    let mut web3 = Web3::new(&format!("{}/jsonrpc", api), Duration::from_secs(120));

    web3.set_header("TRON-PRO-API-KEY", option_env!("API_KEY").unwrap());
    web3.set_check_sync(false);

    let ret = web3
        .eth_get_balance(
            "0xf2846a1e4dafaea38c1660a618277d67605bd2b5"
                .parse()
                .unwrap(),
        )
        .await;

    println!("{:?}", ret);

    let mut client = RpcClient::new(api, Duration::from_secs(120)).unwrap();
    client.set_header("TRON-PRO-API-KEY", option_env!("API_KEY").unwrap());
    let address: Address =
        web30::utils::get_base58_address("0xf2846a1e4dafaea38c1660a618277d67605bd2b5")
            .parse()
            .unwrap();
    let ret = client.get_account_balance(&address).await;
    println!("{:?}", ret);

    let keypair = Keypair::from_hex_key(option_env!("PRIV_KEY").unwrap()).unwrap();
    let mut trade_obj = client
        .trx_transfer(
            &"TPwTVfDDvmWSawsP7Ki1t3ecSBmaFeMMXc"
                .parse::<Address>()
                .unwrap(),
            &"TY5X9ocQACH9YGAyiK3WUxLcLw3t2ethnc"
                .parse::<Address>()
                .unwrap(),
            100u64,
        )
        .await
        .unwrap();

    keypair.sign_transaction(&mut trade_obj).unwrap();

    println!("{:?}", serde_json::to_string(&trade_obj).unwrap());

    // broad cast
    // client.broadcast_transaction(&trade_obj).await.unwrap();
}
