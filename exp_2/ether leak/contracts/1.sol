// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EtherLeak2 {
    function attack() public {
        // Ether leak as it sends all the balance without checks
        payable(msg.sender).transfer(address(this).balance);
    }

    receive() external payable {}
}
