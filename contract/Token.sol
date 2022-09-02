// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "node_modules/@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract Token is ERC20 {
    constructor() ERC20("Token", "TOK") {
        _mint(msg.sender, 100 * 10**decimals());
    }
}
