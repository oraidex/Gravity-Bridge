# General Concept

Transfer flow of Ethereum-based NFTs to Stargaze bases heavily on how ERC20 tokens are swapped through Gravity. Since Gravity.sol contract is immutable, new GravityERC721.sol contract has been deployed, which partialy relies on the security of Gravity.sol, because NFTs can be withdrawn from it only if called by `Gravity.sol` directly. Additionally, `x/nft` and `x/nft-transfer` modules have been added to Gravity Chain to enable NFT and ICS721 support.

ERC721 tokens can be deposited in GravityERC721.sol contract using `sendErc721ToCosmos(tokenContractAddress, cosmosDestinationAddress, tokenId)` call. Successful lockup emits the event that is picked up by validators' orchestrators and `Claim` that NFT has been locked up is reported by the validator. When enough validators report the lockup event, `Attestation` is being processed, which mints the NFT on Gravity Chain. If the destination address is not gravity-based, NFT is being added to interchain transfer queue to be sent to destination chain through IBC. If the destination is on Gravity, then NFT is simply transfered to the indicated receiver. If something goes wrong and freshly minted NFT cannot be delivered, it is being sent to the community pool.

From the user perspective, to transfer NFT to Cosmos, two eth transactions has to be sent - one to approve `GravityERC721.sol` to use user's NFT and the other to deposit the NFT. Then, after the trust period, deposit event is picked up by validators and transfer process continues without any further actions needed from the user.

![eth-to-cosmos](./media/eth-to-cosmos.png)

To transfer Cosmos-originating NFTs to Ethereum, they will have to be locked in the module account. Next, `GravityERC721.sol` will work as a factory and deploy ERC721 contracts that will correspond to Cosmos' NFT classes as well as mint actual NFTs. These `GravityERC721.sol` calls will verify the presence of Gravity validator signatures similarly to how corresponding calls for ERC20s currently work.

From the UX perspective, user will have to transfer NFT from a Cosmos chain to Gravity through IBC and then send a deposit message to gather validator signatures required for minting on the Ethereum side. With ERC20, user would pay a transaction relaying fee to relayers in the same ERC20 token that is being transfered, but with NFT it's not possible, therefore, user will have to *relay* Ethereum transactions manually with the help of UI.

![cosmos-to-eth](./media/cosmos-to-eth.png)

Similar processes will take place to transfer the NFT back to it's origin, but instead of locking it up on the source chain and minting on the destination chain, it's going to be burned at the remote location and unlocked at origin.

// TODO @gjermundgaraba we need ibc channel maintained for every destination chain, no?

# Work done so far
- Added NFT metadata support to existing GravityERC721.sol
- Connected `x/nft` and `x/nft-transfer` to Gravity Chain
- Added `SendERC721ToCosmosClaim` to `x/gravity` and implemented attestiation handling
- Added separate orchestartor oracle flow for `GravityERC721.sol` contract
- Set up a testnet (Gravitygaze) connecting Sepolia and Stargaze Testnet
- frontend? TODO @gjermundgaraba

# How to use

// TODO @gjermundgaraba our UI? Or just etherscan's UI?

# Future plans

At this point, implementation supports only transfers of Ethereum-based NFTs to Cosmos and that's what can be tested on the running testnet. `GraivtyERC721.sol` has a `withdrawERC721()` function, so it would be possible to unlock NFTs held in the contract with the appropriate Gravity Bridge upgrade. However, we plan to implement flows that will allow NFT transfers back and forth (see General Concept above). UI with good UX is also needed for transfers from Cosmos to Ethereum, because of the need to relay Ethereum transactions manually. Additionaly, for a better security, we plan to move NFT bridging logic from `x/gravity` to a dedicated module.

# Other considerations

During our work, we also considered two different approaches:
1. Implementing `x/gravity` directly on Stargaze instead of using IBC and Gravity Bridge.
2. Utilizing Gravity's arbitrary logic calls for NFT bridging.

Ad. 1. Using `x/gravity` directly implies that every Stargaze validator would have to maintain a trusted full Ethereum node as well as run the orchestrator, which would require validators to significantly increase the resources they allocate for validation, therefore this solution is in our opinion inferior to utilizing Gravity Bridge and IBC protocol.

Ad. 2. While we could potentially benefit from less code duplication, this solution would require additional work on arbitrary logic call mechanism itself, as in it's current state, it wouldn't be able to handle our use-case. Moreover, arbitrary logica calls, are very generic in their design, so using them instead of developing our own dedicated module, would decrease the overall security of the design.

--------
**Authors**
