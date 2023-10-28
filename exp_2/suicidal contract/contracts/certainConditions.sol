// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SelfDestruct5 {
    address public owner;
    uint public creationTime;

    constructor() {
        owner = msg.sender;
        creationTime = block.timestamp;
    }

    // Self-destruct if certain time has passed
    function conditionalDestroy() public {
        require(msg.sender == owner, "Only owner can destroy");
        if (block.timestamp >= creationTime + 1 weeks) {
            selfdestruct(payable(owner));
        }
    }
}
