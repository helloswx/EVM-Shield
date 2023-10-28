// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TxOriginAuth {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    // Vulnerability: Using tx.origin for ownership check
    function changeOwner(address newOwner) public {
        require(tx.origin == owner, "Not the owner");

        owner = newOwner;
    }
}
