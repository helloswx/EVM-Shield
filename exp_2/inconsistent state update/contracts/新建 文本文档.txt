// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract InconsistentState5 {
    mapping(address => bool) public authorized;
    address public owner;

    constructor() {
        owner = msg.sender;
        authorized[owner] = true;
    }

    function authorizeUser(address _user) public {
        require(msg.sender == owner, "Not owner");

        // The authorization should be updated before checking the operation
        someSensitiveOperation();

        authorized[_user] = true;
    }

    function someSensitiveOperation() private {
        // xxx
    }
}
