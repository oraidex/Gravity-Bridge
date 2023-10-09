//SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.10;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// One of three testing coins
contract TestERC20Custom is ERC20 {
	constructor(address[] memory _addresses) ERC20("Oraichain Token", "ORAI") {
		for (uint i = 0; i < _addresses.length; ++i) {
			_mint(_addresses[i], 100000000000000000000000000);
		}
	}
}
