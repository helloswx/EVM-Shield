// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TxOriginToken {
    mapping(address => uint) public balances;
    address public admin;

    constructor() {
        admin = msg.sender;
    }

    // Vulnerability: Using tx.origin for admin checks
    function mint(address to, uint amount) public {
        require(tx.origin == admin, "Not the admin");
        balances[to] += amount;
    }
}
