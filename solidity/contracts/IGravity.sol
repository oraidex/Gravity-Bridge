//SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.10;

// This is used purely to avoid stack too deep errors
// represents everything about a given validator set
struct ValsetArgs {
	// the validators in this set, represented by an Ethereum address
	address[] validators;
	// the powers of the given validators in the same order as above
	uint256[] powers;
	// the nonce of this validator set
	uint256 valsetNonce;
	// the reward amount denominated in the below reward token, can be
	// set to zero
	uint256 rewardAmount;
	// the reward token, should be set to the zero address if not being used
	address rewardToken;
}

// This represents a validator signature
struct Signature {
	uint8 v;
	bytes32 r;
	bytes32 s;
}

interface IGravity {
	function sendToCosmos(
		address _tokenContract,
		string calldata _destination,
		uint256 _amount
	) external;

	function submitBatch(
		// The validators that approve the batch
		ValsetArgs calldata _currentValset,
		// These are arrays of the parts of the validators signatures
		Signature[] calldata _sigs,
		// The batch of transactions
		uint256[] calldata _amounts,
		address[] calldata _destinations,
		uint256[] calldata _fees,
		uint256 _batchNonce,
		address _tokenContract,
		// a block height beyond which this batch is not valid
		// used to provide a fee-free timeout
		uint256 _batchTimeout
	) external;
}
