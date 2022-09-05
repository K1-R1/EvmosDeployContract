// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "node_modules/@openzeppelin/contracts/token/ERC20/ERC20.sol";

/** @title Token
 *  @author K1-R1
 *  @notice This contract is an implemenation of ERC20 specification
 */

contract Token is ERC20 {
    /**
     * @notice 100 TOK are minted, and assigned to the deployer of the contract
     *         upon deployment.
     * @dev decimals is set to 18
     */
    constructor() ERC20("Token", "TOK") {
        _mint(msg.sender, 100 * 10**decimals());
    }
}
