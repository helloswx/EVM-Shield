// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EtherLeak5 {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    function attack() public {
        // Ether leak due to lack of owner check
        payable(msg.sender).transfer(address(this).balance);
    }

    receive() external payable {}
}
