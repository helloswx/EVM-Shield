// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SelfDestruct1 {
    // Constructor sets the owner
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    // Self-destruct function without proper access control
    function destroyContract() public {
        selfdestruct(payable(owner));
    }
}
