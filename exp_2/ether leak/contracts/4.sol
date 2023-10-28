// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EtherLeak4 {
    address public richest;
    uint public mostSent;

    constructor() {
        richest = msg.sender;
        mostSent = 0;
    }

    function becomeRichest() public payable {
        require(msg.value > mostSent, "Not enough Ether sent");

        // Ether leak as it sends the Ether back to the previous richest without proper checks
        payable(richest).transfer(address(this).balance);
        richest = msg.sender;
        mostSent = msg.value;
    }

    receive() external payable {}
}
