// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TxOriginWallet {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    // Vulnerability: Using tx.origin for withdrawal authorization
    function withdraw(uint amount) public {
        require(tx.origin == owner, "Not the owner");
        payable(tx.origin).transfer(amount);
    }

    receive() external payable {}
}
