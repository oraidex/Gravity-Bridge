pub const TRANSACTION_BATCH_EXECUTED_EVENT_SIG: &str =
    "TransactionBatchExecutedEvent(uint256,address,uint256)";

pub const SENT_TO_COSMOS_EVENT_SIG: &str =
    "SendToCosmosEvent(address,address,string,uint256,uint256)";

pub const SENT_ERC721_TO_COSMOS_EVENT_SIG: &str =
    "SendERC721ToCosmosEvent(address,address,string,uint256,uint256,string)";

pub const GRAVITYERC721_DEPLOYED_EVENT_SIG: &str =
    "GravityERC721DeployedEvent()";

pub const ERC20_DEPLOYED_EVENT_SIG: &str =
    "ERC20DeployedEvent(string,address,string,string,uint8,uint256)";

pub const LOGIC_CALL_EVENT_SIG: &str = "LogicCallEvent(bytes32,uint256,bytes,uint256)";

pub const VALSET_UPDATED_EVENT_SIG: &str =
    "ValsetUpdatedEvent(uint256,uint256,uint256,address,address[],uint256[])";
