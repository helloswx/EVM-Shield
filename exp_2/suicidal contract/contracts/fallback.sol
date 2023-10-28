// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SelfDestruct2 {
    // Self-destruct triggered in fallback function
    fallback() external payable {
        selfdestruct(payable(msg.sender));
    }
}
