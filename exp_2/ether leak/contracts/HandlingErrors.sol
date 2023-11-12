// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EtherLeak1 {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    receive() external payable {}

    function attack(uint _amount) public {
        require(msg.sender == owner, "Not owner");
        require(_amount <= address(this).balance, "Insufficient balance");

        // Ether leak due to sending to an incorrect address
        payable(address(0)).transfer(_amount);
    }
}
